import { Button, Modal, ModalBody, ModalContent, ModalHeader } from '@heroui/react';
import { useMemo } from 'react';
import './Delta.css';

function ordinal(value: number): string {
  if (value % 100 >= 11 && value % 100 <= 13) return 'th';
  return ['th', 'st', 'nd', 'rd'][value % 10] ?? 'th';
}

type TablePosition = 'top' | 'right' | 'bottom' | 'left';

export type DeltaPlayer = {
  id: number;
  name: string;
  score: number;
  position: TablePosition;
};

type DeltaProps = {
  isOpen: boolean;
  onClose: () => void;
  selectedPlayerId: number | null;
  players: DeltaPlayer[];
};

type PlaceTarget = { place: number; handValue: number };
type LossTarget = {
  placesLost: number;
  finalPlace: number;
  tiedThrough: number;
  handValue: number;
};
type OpponentArea = 'top' | 'left' | 'right';

const MIN_HAND = 8;

type HandMode = 'self' | 'direct' | 'deal-in';
type Scenario = {
  mode: HandMode;
  selectedPlayerId: number;
  opponentId?: number;
};
type LinearScore = { id: number; base: number; slope: number };
type Rank = { first: number; last: number };

function rankOf(scores: LinearScore[], selectedId: number, handValue: number): Rank {
  const own = scores.find((p) => p.id === selectedId);
  if (!own) return { first: 0, last: 0 };
  const value = own.base + own.slope * handValue;
  let above = 0;
  let equal = 0;
  for (const player of scores) {
    const otherValue = player.base + player.slope * handValue;
    if (otherValue > value) above++;
    if (otherValue === value) equal++;
  }
  return { first: above + 1, last: above + equal };
}

// MCR: self-draw = +3(H+8) / -(H+8);
// discard win = +(H+24) / -(H+8) from discarder / -8 from others.
function buildScenario(players: DeltaPlayer[], scenario: Scenario): LinearScore[] {
  const { mode, selectedPlayerId, opponentId } = scenario;
  return players.map((player) => {
    let base = player.score;
    let slope = 0;
    if (mode === 'self') {
      if (player.id === selectedPlayerId) {
        base += 24;
        slope = 3;
      } else {
        base -= 8;
        slope = -1;
      }
    } else {
      const winnerId = mode === 'direct' ? selectedPlayerId : opponentId;
      const discarderId = mode === 'direct' ? opponentId : selectedPlayerId;
      if (player.id === winnerId) {
        base += 24;
        slope = 1;
      } else if (player.id === discarderId) {
        base -= 8;
        slope = -1;
      } else {
        base -= 8;
      }
    }
    return { id: player.id, base, slope };
  });
}

// Rankings can only change where two linear score functions cross.
// Evaluate those boundaries instead of scanning every possible hand value.
function candidateHands(scores: LinearScore[]): number[] {
  const values = new Set<number>([MIN_HAND]);
  for (let i = 0; i < scores.length; i++) {
    for (let j = i + 1; j < scores.length; j++) {
      const a = scores[i];
      const b = scores[j];
      if (a.slope === b.slope) continue;
      const crossing = (b.base - a.base) / (a.slope - b.slope);
      if (crossing < MIN_HAND - 1) continue;
      const floor = Math.floor(crossing);
      for (const hand of [floor, floor + 1, floor + 2]) {
        if (hand >= MIN_HAND && Number.isSafeInteger(hand)) values.add(hand);
      }
    }
  }
  return [...values].sort((a, b) => a - b);
}

function currentRank(players: DeltaPlayer[], selectedPlayerId: number): Rank {
  return rankOf(players.map((p) => ({ id: p.id, base: p.score, slope: 0 })), selectedPlayerId, 0);
}

function getGainTargets(
  players: DeltaPlayer[],
  selectedPlayerId: number,
  mode: 'self' | 'direct',
  discarderId?: number
): PlaceTarget[] {
  const initial = currentRank(players, selectedPlayerId);
  if (initial.first === 0 || (initial.first === 1 && initial.last === 1)) return [];
  const scores = buildScenario(players, {
    mode,
    selectedPlayerId,
    opponentId: discarderId,
  });
  const bestByPlace = new Map<number, number>();
  for (const handValue of candidateHands(scores)) {
    const rank = rankOf(scores, selectedPlayerId, handValue);
    // Only sole places: never display a tie as a promotion.
    if (rank.first !== rank.last || rank.first >= initial.last) continue;
    if (!bestByPlace.has(rank.first)) bestByPlace.set(rank.first, handValue);
  }
  // The same minimum hand may skip multiple places. Keep the highest.
  const bestAtHand = new Map<number, number>();
  for (const [place, handValue] of bestByPlace) {
    bestAtHand.set(handValue, Math.min(place, bestAtHand.get(handValue) ?? 5));
  }
  return [...bestAtHand].map(([handValue, place]) => ({ place, handValue }))
    .sort((a, b) => a.handValue - b.handValue);
}

function getDiscardTargets(
  players: DeltaPlayer[],
  selectedPlayerId: number,
  winnerId: number
): LossTarget[] {
  const initial = currentRank(players, selectedPlayerId);
  if (initial.first === 0 || initial.last === players.length) return [];
  const scores = buildScenario(players, {
    mode: 'deal-in',
    selectedPlayerId,
    opponentId: winnerId,
  });
  const targets: LossTarget[] = [];
  let worstRankSum = initial.first + initial.last;
  for (const handValue of candidateHands(scores)) {
    const rank = rankOf(scores, selectedPlayerId, handValue);
    const rankSum = rank.first + rank.last;
    if (rankSum <= worstRankSum) continue;
    worstRankSum = rankSum;
    targets.push({
      placesLost: Math.max(1, rank.last - initial.first),
      finalPlace: rank.first,
      tiedThrough: rank.last,
      handValue,
    });
  }
  return targets;
}

function getRelativePositions(players: DeltaPlayer[], selectedPlayerId: number) {
  const selectedPlayer = players.find((p) => p.id === selectedPlayerId) ?? null;
  if (!selectedPlayer) {
    return { selectedPlayer: null, topPlayer: null, leftPlayer: null, rightPlayer: null };
  }

  const order: TablePosition[] = ['bottom', 'right', 'top', 'left'];
  const index = order.indexOf(selectedPlayer.position);
  const at = (offset: number) =>
    players.find((p) => p.position === order[(index + offset) % order.length]) ?? null;

  return {
    selectedPlayer,
    rightPlayer: at(1),
    topPlayer: at(2),
    leftPlayer: at(3),
  };
}

type TargetListProps = {
  targets: PlaceTarget[];
  emptyText: string;
  className: string;
};

function TargetList({ targets, emptyText, className }: TargetListProps) {
  if (targets.length === 0) {
    return <div className={`${className}__empty text-center text-sm text-default-500`}>{emptyText}</div>;
  }
  return (
    <div className={`${className}__targets`}>
      {targets.map(({ place, handValue }) => (
        <div key={place} className={`${className}__target rounded-medium bg-default-100 px-2 py-1.5 text-sm`}>
          <span className={`${className}__value font-semibold tabular-nums`}>{handValue}</span>
          <span className={`${className}__destination text-xs text-default-500`}>→ {place}{ordinal(place)}</span>
        </div>
      ))}
    </div>
  );
}

type OpponentCardProps = {
  player: DeltaPlayer;
  directTargets: PlaceTarget[];
  discardTargets: LossTarget[];
  area: OpponentArea;
};

function OpponentCard({ player, directTargets, discardTargets, area }: OpponentCardProps) {
  return (
    <section className={`delta-opponent delta-opponent--${area} rounded-large border border-divider bg-content1 p-4 shadow-small`}>
      <div className="delta-opponent__identity border-b border-divider pb-3 text-center">
        <div className="delta-opponent__name font-semibold text-foreground" title={player.name}>{player.name}</div>
        <div className={[
          'delta-opponent__score font-semibold text-foreground tabular-nums',
          player.score < 0 ? 'text-danger' : '',
        ].filter(Boolean).join(' ')}>
          {player.score.toLocaleString('en-US')}
        </div>
      </div>

      {directTargets.length > 0 && (
        <div className="delta-direct">
          <div className="delta-direct__title rounded-medium bg-primary-100 px-2 py-1 text-center text-xs font-semibold text-primary" title={`Win on a discard from ${player.name}`}>
            Win from this player
          </div>
          <TargetList targets={directTargets} emptyText="" className="delta-direct" />
        </div>
      )}

      {discardTargets.length > 0 && (
        <div className="delta-discard">
          <div className="delta-discard__title rounded-medium bg-danger-100 px-2 py-1 text-center text-xs font-semibold text-danger" title={`Discard into a win by ${player.name}`}>
            Deal into this player
          </div>
          <div className="delta-discard__targets">
            {discardTargets.map(({ finalPlace, tiedThrough, handValue, placesLost }) => (
              <div key={`${finalPlace}-${tiedThrough}-${handValue}`} className="delta-discard__target rounded-medium bg-default-100 px-2 py-1.5 text-sm">
                <span className="delta-discard__value font-semibold tabular-nums">{handValue}</span>
                <span className="delta-discard__destination text-xs text-default-500">→ {finalPlace}{ordinal(finalPlace)}{tiedThrough > finalPlace ? `–${tiedThrough}${ordinal(tiedThrough)} (tie)` : ''} </span>
                <span className="delta-discard__arrows text-xs text-danger" aria-label={`${placesLost} places lost`}>
                  {'↓'.repeat(placesLost)}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}
    </section>
  );
}

export default function Delta({ isOpen, onClose, selectedPlayerId, players }: DeltaProps) {
  // CurrentGame passes a newly created array on each render.
  const playersKey = JSON.stringify(players.map((p) =>
    [p.id, p.name, p.score, p.position]
  ));

  const calculations = useMemo(() => {
    const empty = {
      rank: { first: 0, last: 0 },
      selectedPlayer: null as DeltaPlayer | null,
      topPlayer: null as DeltaPlayer | null,
      leftPlayer: null as DeltaPlayer | null,
      rightPlayer: null as DeltaPlayer | null,
      selfTargets: [] as PlaceTarget[],
      opponentTargets: {} as Record<number, {
        direct: PlaceTarget[];
        discard: LossTarget[];
      }>,
    };
    if (!isOpen || selectedPlayerId === null) return empty;
    const positions = getRelativePositions(players, selectedPlayerId);
    if (!positions.selectedPlayer) return empty;
    const opponentTargets: typeof empty.opponentTargets = {};
    for (const opponent of players) {
      if (opponent.id === selectedPlayerId) continue;
      opponentTargets[opponent.id] = {
        direct: getGainTargets(players, selectedPlayerId, 'direct', opponent.id),
        discard: getDiscardTargets(players, selectedPlayerId, opponent.id),
      };
    }
    return {
      ...positions,
      rank: currentRank(players, selectedPlayerId),
      selfTargets: getGainTargets(players, selectedPlayerId, 'self'),
      opponentTargets,
    };
  // The primitive key includes every player field used by these calculations.
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen, selectedPlayerId, playersKey]);

  const {
    rank, selectedPlayer, topPlayer, leftPlayer, rightPlayer,
    selfTargets, opponentTargets,
  } = calculations;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      size="full"
      placement="center"
      scrollBehavior="inside"
      classNames={{
        base: 'delta-modal bg-background text-foreground rounded-none',
        header: 'delta-modal__header border-b border-divider',
        body: 'delta-modal__body',
        closeButton: 'hidden',
      }}
    >
      <ModalContent>
        {() => (
          <>
            <ModalHeader className="relative pr-14">
              <div className="flex flex-col gap-1">
                <span>Position thresholds</span>
                {selectedPlayer && (
                  <span className="text-sm font-normal text-default-500">
                    {selectedPlayer.name} · {rank.first}{ordinal(rank.first)}{rank.last > rank.first ? `–${rank.last}${ordinal(rank.last)} (tie)` : ''} place
                  </span>
                )}
              </div>
              <Button
                isIconOnly
                variant="light"
                radius="full"
                size="sm"
                aria-label="Close position thresholds"
                onPress={onClose}
                className="absolute right-3 top-1/2 z-20 -translate-y-1/2 text-default-600"
              >
                <svg aria-hidden="true" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                  <path d="M18 6 6 18M6 6l12 12" />
                </svg>
              </Button>
            </ModalHeader>
            <ModalBody>
              {!selectedPlayer ? (
                <div className="flex min-h-48 items-center justify-center text-default-500">Player not found</div>
              ) : (
                <div className="delta-layout">
                  {topPlayer && (
                    <OpponentCard player={topPlayer}
                      directTargets={opponentTargets[topPlayer.id].direct}
                      discardTargets={opponentTargets[topPlayer.id].discard}
                      area="top" />
                  )}
                  {leftPlayer && (
                    <OpponentCard player={leftPlayer}
                      directTargets={opponentTargets[leftPlayer.id].direct}
                      discardTargets={opponentTargets[leftPlayer.id].discard}
                      area="left" />
                  )}
                  <section className="delta-selected rounded-large border-2 border-primary bg-content1 p-4 shadow-small">
                    <div className="flex flex-col items-center gap-2">
                      
                      <div className="delta-selected__identity">
                        <div className="delta-selected__name font-semibold text-foreground" title={selectedPlayer.name}>
                          {selectedPlayer.name}
                        </div>
                        <div className={[
                          'delta-selected__score font-semibold text-foreground tabular-nums',
                          selectedPlayer.score < 0 ? 'text-danger' : '',
                        ].filter(Boolean).join(' ')}>
                          {selectedPlayer.score.toLocaleString('en-US')}
                        </div>
                        <div className="delta-selected__place text-xs text-default-500">
                          Currently {rank.first}{ordinal(rank.first)}{rank.last > rank.first ? `–${rank.last}${ordinal(rank.last)} (tie)` : ''} place
                        </div>
                      </div>
                    </div>
                    {selfTargets.length > 0 && (
                      <div className="delta-self">
                        <div className="delta-self__title rounded-medium bg-success-100 px-2 py-1 text-center text-xs font-semibold text-success">Self-draw</div>
                        <TargetList targets={selfTargets} emptyText="" className="delta-self" />
                      </div>
                    )}
                  </section>
                  {rightPlayer && (
                    <OpponentCard player={rightPlayer}
                      directTargets={opponentTargets[rightPlayer.id].direct}
                      discardTargets={opponentTargets[rightPlayer.id].discard}
                      area="right" />
                  )}
                  <div className="delta-legend text-xs text-default-500">
                    <span>Self-draw: minimum hand value</span>
                    <span>Discard win: minimum hand value</span>
                    <span>→ resulting place</span>
                  </div>
                </div>
              )}
            </ModalBody>
          </>
        )}
      </ModalContent>
    </Modal>
  );
}
