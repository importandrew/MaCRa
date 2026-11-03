import { Button, Modal, ModalBody, ModalContent, ModalHeader } from '@heroui/react';
import { useMemo, useRef } from 'react';
import './Delta.css';

function windLabel(wind: string | string[] | undefined): string {
  // The API model types wind as typeof WIND_SYMBOLS (string[]), although a
  // player wind is normally one symbol. Do not display the entire symbol list.
  return typeof wind === 'string' ? wind : '';
}

function ordinal(value: number): string {
  if (value % 100 >= 11 && value % 100 <= 13) return 'th';
  return ['th', 'st', 'nd', 'rd'][value % 10] ?? 'th';
}

type TablePosition = 'top' | 'right' | 'bottom' | 'left';

export type DeltaPlayer = {
  id: number;
  name: string;
  wind?: string | string[];
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
type OpponentArea = 'top' | 'left' | 'right' | 'bottom';

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
  const seenPlaces = new Set<number>();
  for (const handValue of candidateHands(scores)) {
    const rank = rankOf(scores, selectedPlayerId, handValue);
    // A shared place is not a new threshold: show only sole lower places.
    if (rank.first !== rank.last || rank.first <= initial.first) continue;
    if (seenPlaces.has(rank.first)) continue;
    seenPlaces.add(rank.first);
    targets.push({
      placesLost: rank.first - initial.first,
      finalPlace: rank.first,
      tiedThrough: rank.first,
      handValue,
    });
  }
  return targets;
}

// The table always stays oriented toward the device owner (bottom seat).
// The selected player's self-draw card moves to the centre; their original
// seat is intentionally left vacant, without rotating any other seats.
function getTablePositions(players: DeltaPlayer[], selectedPlayerId: number) {
  const selectedPlayer = players.find((p) => p.id === selectedPlayerId) ?? null;
  const opponentAt = (position: TablePosition) =>
    players.find((p) => p.position === position && p.id !== selectedPlayerId) ?? null;

  return {
    selectedPlayer,
    topPlayer: opponentAt('top'),
    leftPlayer: opponentAt('left'),
    rightPlayer: opponentAt('right'),
    bottomPlayer: opponentAt('bottom'),
  };
}

type TargetListProps = {
  targets: PlaceTarget[];
};

function TargetList({ targets }: TargetListProps) {
  return (
    <div className="delta-thresholds">
      {targets.map(({ place, handValue }) => (
        <span key={place} className="delta-threshold" title={`Reach ${ordinal(place)} place`}>
          <strong className="tabular-nums">{handValue}</strong>
          <span className="text-default-500">→ {place}{ordinal(place)}</span>
        </span>
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

function OpponentCard({
  player,
  directTargets,
  discardTargets,
  area,
}: OpponentCardProps) {
  return (
    <section className={`delta-opponent delta-opponent--${area} rounded-large border border-divider bg-content1 p-2 sm:p-3 shadow-small`}>
      <div className="delta-identity border-b border-divider pb-1.5">
        <div className="delta-wind text-default-500" aria-label={windLabel(player.wind) ? `Wind ${windLabel(player.wind)}` : undefined}>
          {windLabel(player.wind)}
        </div>
        <div className="delta-person">
          <div className="delta-name font-semibold text-foreground" title={player.name}>{player.name}</div>
          <div className={`delta-score tabular-nums ${player.score < 0 ? 'text-danger' : 'text-default-500'}`}>
            {player.score.toLocaleString('en-US')}
          </div>
        </div>
      </div>
      {directTargets.length > 0 && (
        <div className="delta-action delta-action--gain bg-success-100 text-success-700" title={`Win on a discard from ${player.name}`}>
          <span className="delta-action__label">Win</span>
          <TargetList targets={directTargets} />
        </div>
      )}
      {discardTargets.length > 0 && (
        <div className="delta-action delta-action--loss bg-danger-100 text-danger-700" title={`Deal into a win by ${player.name}`}>
          <span className="delta-action__label">Deal in</span>
          <div className="delta-thresholds">
            {discardTargets.map(({ finalPlace, tiedThrough, handValue }) => (
              <span key={`${finalPlace}-${tiedThrough}-${handValue}`} className="delta-threshold">
                <strong className="tabular-nums">{handValue}</strong>
                <span className="text-default-500">
                  → {finalPlace}{ordinal(finalPlace)}
                  {tiedThrough > finalPlace ? `–${tiedThrough}${ordinal(tiedThrough)}` : ''}
                </span>
              </span>
            ))}
          </div>
        </div>
      )}
    </section>
  );
}

function computeEmpty() {
  return {
    rank: { first: 0, last: 0 },
    selectedPlayer: null as DeltaPlayer | null,
    topPlayer: null as DeltaPlayer | null,
    leftPlayer: null as DeltaPlayer | null,
    rightPlayer: null as DeltaPlayer | null,
    bottomPlayer: null as DeltaPlayer | null,
    selfTargets: [] as PlaceTarget[],
    opponentTargets: {} as Record<number, {
      direct: PlaceTarget[];
      discard: LossTarget[];
    }>,
  };
}

export default function Delta({ isOpen, onClose, selectedPlayerId, players }: DeltaProps) {
  // CurrentGame passes a newly created array on each render.
  const playersKey = JSON.stringify(players.map((p) =>
    [p.id, p.name, p.score, p.position, p.wind]
  ));

  const lastCalculations = useRef<ReturnType<typeof computeEmpty> | null>(null);
  const calculations = useMemo(() => {
    const empty = computeEmpty();
    if (!isOpen || selectedPlayerId === null) return lastCalculations.current ?? empty;
    const positions = getTablePositions(players, selectedPlayerId);
    if (!positions.selectedPlayer) return empty;
    const opponentTargets: typeof empty.opponentTargets = {};
    for (const opponent of players) {
      if (opponent.id === selectedPlayerId) continue;
      opponentTargets[opponent.id] = {
        direct: getGainTargets(players, selectedPlayerId, 'direct', opponent.id),
        discard: getDiscardTargets(players, selectedPlayerId, opponent.id),
      };
    }
    const result = {
      ...positions,
      rank: currentRank(players, selectedPlayerId),
      selfTargets: getGainTargets(players, selectedPlayerId, 'self'),
      opponentTargets,
    };
    lastCalculations.current = result;
    return result;
  // The primitive key includes every player field used by these calculations.
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen, selectedPlayerId, playersKey]);

  const {
    rank, selectedPlayer, topPlayer, leftPlayer, rightPlayer, bottomPlayer,
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
                  <section className="delta-selected rounded-large border-2 border-primary bg-content1 p-2 sm:p-3 shadow-small">
                    <div className="delta-identity border-b border-divider pb-1.5">
                      <div className="delta-wind text-primary">{windLabel(selectedPlayer.wind)}</div>
                      <div className="delta-person">
                        <div className="delta-name font-semibold" title={selectedPlayer.name}>
                          {selectedPlayer.name}
                        </div>
                        <div className="delta-score tabular-nums">
                          {selectedPlayer.score.toLocaleString('en-US')}
                        </div>
                        <div className="text-xs text-default-500">
                          {rank.first}{ordinal(rank.first)}
                          {rank.last > rank.first ? `–${rank.last}${ordinal(rank.last)} (tie)` : ''}
                        </div>
                      </div>
                    </div>
                    {selfTargets.length > 0 && (
                      <div className="delta-action delta-action--gain bg-success-100 text-success-700">
                        <span className="delta-action__label">Self</span>
                        <TargetList targets={selfTargets} />
                      </div>
                    )}
                  </section>
                  {rightPlayer && (
                    <OpponentCard player={rightPlayer}
                      directTargets={opponentTargets[rightPlayer.id].direct}
                      discardTargets={opponentTargets[rightPlayer.id].discard}
                      area="right" />
                  )}
                  {bottomPlayer && (
                    <OpponentCard player={bottomPlayer}
                      directTargets={opponentTargets[bottomPlayer.id].direct}
                      discardTargets={opponentTargets[bottomPlayer.id].discard}
                      area="bottom" />
                  )}
                  <div className="delta-legend text-xs text-default-500">
                    Green: gain a place · Red: risk losing a place · Values: minimum fan
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
