import { Modal } from '@heroui/react';
import { useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react';
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

type PlaceTarget = { place: number; handValue: number; tiedThrough: number };
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
  return rankOf(
    players.map((p) => ({ id: p.id, base: p.score, slope: 0 })),
    selectedPlayerId,
    0
  );
}

// One threshold per attainable place. A sole place takes precedence over a tie
// for the same place, even when the tie can be reached with fewer fan.
// If a place is skipped entirely (e.g. jumping from 4th to 2nd), it is omitted.
function selectTargets(
  scores: LinearScore[],
  selectedPlayerId: number,
  isRelevant: (rank: Rank) => boolean
): PlaceTarget[] {
  const candidates = new Map<number, PlaceTarget>();
  for (const handValue of candidateHands(scores)) {
    const rank = rankOf(scores, selectedPlayerId, handValue);
    if (!isRelevant(rank)) continue;
    const target: PlaceTarget = {
      place: rank.first,
      tiedThrough: rank.last,
      handValue,
    };
    const existing = candidates.get(target.place);
    if (!existing ||
      (existing.tiedThrough !== existing.place && target.tiedThrough === target.place) ||
      ((existing.tiedThrough === existing.place) === (target.tiedThrough === target.place) &&
        target.handValue < existing.handValue)
    ) {
      candidates.set(target.place, target);
    }
  }
  return [...candidates.values()].sort((a, b) => a.place - b.place);
}

function getGainTargets(
  players: DeltaPlayer[],
  selectedPlayerId: number,
  mode: 'self' | 'direct',
  discarderId?: number
): PlaceTarget[] {
  const initial = currentRank(players, selectedPlayerId);
  if (initial.first === 1 && initial.last === 1) return [];
  const scores = buildScenario(players, {
    mode, selectedPlayerId, opponentId: discarderId,
  });
  return selectTargets(scores, selectedPlayerId,
    (rank) => rank.first < initial.first ||
      (rank.first === initial.first && rank.last < initial.last));
}

function getDiscardTargets(
  players: DeltaPlayer[],
  selectedPlayerId: number,
  winnerId: number
): LossTarget[] {
  const initial = currentRank(players, selectedPlayerId);
  if (initial.first === players.length && initial.last === players.length) return [];
  const scores = buildScenario(players, {
    mode: 'deal-in', selectedPlayerId, opponentId: winnerId,
  });
  return selectTargets(scores, selectedPlayerId,
    (rank) => rank.first > initial.first ||
      (rank.first === initial.first && rank.last > initial.last))
    .map(({ place, tiedThrough, handValue }) => ({
      placesLost: Math.max(0, place - initial.first),
      finalPlace: place,
      tiedThrough,
      handValue,
    }));
}

// Keep every player at their original table seat, including the selected one.
function getTablePositions(players: DeltaPlayer[], selectedPlayerId: number) {
  return {
    selectedPlayer: players.find(p => p.id === selectedPlayerId) ?? null,
    topPlayer: players.find(p => p.position === 'top') ?? null,
    leftPlayer: players.find(p => p.position === 'left') ?? null,
    rightPlayer: players.find(p => p.position === 'right') ?? null,
    bottomPlayer: players.find(p => p.position === 'bottom') ?? null,
  };
}

type TargetListProps = {
  targets: PlaceTarget[];
};

function TargetList({ targets }: TargetListProps) {
  return (
    <div className="delta-thresholds">
      {targets.map(({ place, handValue, tiedThrough }) => (
        <span
          key={place}
          className="delta-threshold"
          title={`Reach ${place}${ordinal(place)} place${tiedThrough > place ? " (tie)" : ""}`}
        >
          <strong className="tabular-nums">{handValue}</strong>
          <span className="text-muted">
            → {place}{ordinal(place)}{tiedThrough > place ? ' tie' : ''}
          </span>
        </span>
      ))}
    </div>
  );
}

type OpponentCardProps = {
  player: DeltaPlayer;
  selectedScore: number;
  directTargets: PlaceTarget[];
  discardTargets: LossTarget[];
  area: OpponentArea;
};

function OpponentCard({
  player,
  selectedScore,
  directTargets,
  discardTargets,
  area,
}: OpponentCardProps) {
  return (
    <section
      data-delta-help="opponent" className={`delta-opponent delta-opponent--${area} rounded-lg border border-border bg-surface p-2 shadow-sm sm:p-3`}
    >
      <div data-delta-help="identity" className="delta-identity border-b border-border pb-1.5">
        <div
          className="delta-wind text-muted"
          aria-label={windLabel(player.wind) ? `Wind ${windLabel(player.wind)}` : undefined}
        >
          {windLabel(player.wind)}
        </div>
        <div className="delta-person">
          <div className="delta-name font-semibold text-foreground" title={player.name}>
            {player.name}
          </div>
          <div className="delta-score-line">
            <span data-delta-help="score" className="delta-score tabular-nums text-muted">
              {player.score.toLocaleString('en-US')}
            </span>
            <span
              data-delta-help="gap" className={`delta-score-gap tabular-nums ${
                selectedScore > player.score
                  ? 'delta-score-gap--ahead'
                  : selectedScore < player.score
                    ? 'delta-score-gap--behind'
                    : 'delta-score-gap--equal'
              }`}
              title="Selected player's score minus this player's score"
              aria-label={`Score advantage of selected player: ${selectedScore - player.score}`}
            >
              {selectedScore - player.score > 0 ? '+' : ''}
              {(selectedScore - player.score).toLocaleString('en-US')}
            </span>
          </div>
        </div>
      </div>
      {directTargets.length > 0 && (
        <div
          data-delta-help="direct" className="delta-action delta-action--direct"
          title={`Win on a discard from ${player.name}`}
        >
          <TargetList targets={directTargets} />
        </div>
      )}
      {discardTargets.length > 0 && (
        <div
          data-delta-help="loss" className="delta-action delta-action--loss"
          title={`Deal into a win by ${player.name}`}
        >
          <div className="delta-thresholds">
            {discardTargets.map(({ finalPlace, tiedThrough, handValue }) => (
              <span
                key={`${finalPlace}-${tiedThrough}-${handValue}`}
                className="delta-threshold"
              >
                <strong className="tabular-nums">{handValue}</strong>
                <span className="text-muted">
                  → {finalPlace}{ordinal(finalPlace)}
                  {tiedThrough > finalPlace ? ' tie' : ''}
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
    opponentTargets: {} as Record<
      number,
      { direct: PlaceTarget[]; discard: LossTarget[] }
    >,
  };
}


type HelpAnnotation = {
  id: string;
  label: string;
  rect: { left: number; top: number; width: number; height: number };
  labelRect: { left: number; top: number; width: number; height: number };
};

const HELP_ITEMS = [
  { id: 'selected', label: 'Selected player' },
  { id: 'score', label: 'Current score' },
  { id: 'gap', label: 'Δ: selected score − opponent score' },
  { id: 'self', label: 'Self-draw: minimum fan to reach this place' },
  { id: 'direct', label: 'Win on this player’s discard' },
  { id: 'loss', label: 'Deal into this player’s winning hand' },
] as const;

function HelpOverlay({
  onClose,
  onFallback,
}: {
  onClose: () => void;
  onFallback: () => void;
}) {
  const [annotations, setAnnotations] = useState<HelpAnnotation[]>([]);
  const [fallback, setFallback] = useState(false);
  const [visible, setVisible] = useState(false);

  useLayoutEffect(() => {
    const calculate = () => {
      const width = window.innerWidth;
      const height = window.innerHeight;
      const elements = HELP_ITEMS.flatMap((item) => {
        // Prefer the selected player's score; all other labels refer to an opponent.
        const nodes = Array.from(document.querySelectorAll<HTMLElement>(
          `[data-delta-help="${item.id}"]`
        ));
        const node = item.id === 'score'
          ? nodes.find((el) => el.closest('.delta-selected')) ?? nodes[0]
          : item.id === 'gap' || item.id === 'direct' || item.id === 'loss'
            ? nodes.find((el) => el.getBoundingClientRect().width > 0)
            : nodes[0];
        if (!node) return [];
        const r = node.getBoundingClientRect();
        if (!r.width || !r.height) return [];
        return [{
          id: item.id,
          label: item.label,
          rect: { left: r.left, top: r.top, width: r.width, height: r.height },
        }];
      });
      const occupied: Array<{ left: number; top: number; width: number; height: number }> = [];
      const result: HelpAnnotation[] = [];
      const margin = 10;
      // Search the whole viewport for the nearest free label position.
      // Previously only four adjacent positions were considered; these are
      // usually blocked by the other three player cards, even on desktop.
      for (const item of elements) {
        const labelWidth = Math.min(width - 2 * margin, width < 600 ? 124 : 184);
        const labelHeight = width < 600 ? 42 : 38;
        const r = item.rect;
        const cx = r.left + r.width / 2;
        const cy = r.top + r.height / 2;
        const candidates: Array<{ left: number; top: number; distance: number }> = [];
        const step = width < 600 ? 14 : 24;
        for (let top = 48; top <= height - labelHeight - margin; top += step) {
          for (let left = margin; left <= width - labelWidth - margin; left += step) {
            const lx = left + labelWidth / 2;
            const ly = top + labelHeight / 2;
            candidates.push({ left, top, distance: Math.hypot(lx - cx, ly - cy) });
          }
        }
        candidates.sort((a, b) => a.distance - b.distance);
        const free = candidates.find((pos) =>
          occupied.every((o) =>
            pos.left + labelWidth + 6 <= o.left ||
            pos.left >= o.left + o.width + 6 ||
            pos.top + labelHeight + 6 <= o.top ||
            pos.top >= o.top + o.height + 6
          ) &&
          elements.every((other) =>
            pos.left + labelWidth + 4 <= other.rect.left ||
            pos.left >= other.rect.left + other.rect.width + 4 ||
            pos.top + labelHeight + 4 <= other.rect.top ||
            pos.top >= other.rect.top + other.rect.height + 4
          )
        );
        if (!free) {
          setFallback(true);
          setVisible(true);
          return;
        }
        const labelRect = { left: free.left, top: free.top, width: labelWidth, height: labelHeight };
        occupied.push(labelRect);
        result.push({ ...item, labelRect });
      }
      if (result.length < 3) {
        setFallback(true);
      } else {
        setFallback(false);
        setAnnotations(result);
      }
      setVisible(true);
    };
    calculate();
    window.addEventListener('resize', calculate);
    return () => window.removeEventListener('resize', calculate);
  }, []);

  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        event.stopPropagation();
        onClose();
      }
    };
    document.addEventListener('keydown', onKey, true);
    return () => document.removeEventListener('keydown', onKey, true);
  }, [onClose]);

  if (!visible) return null;
  if (fallback) {
    return (
          <div className="delta-help-backdrop">
            <button
              type="button"
              className="delta-help-backdrop-button"
              onClick={onClose}
              aria-label="Close help by clicking outside"
            />
            <section
              className="delta-help-panel"
              role="dialog"
              aria-modal="true"
              aria-label="How to read Delta"
            >
          <div className="delta-help-heading">
            <h2>How to read Delta</h2>
            <button className="delta-help-close" type="button" onClick={onClose} aria-label="Close help">×</button>
          </div>
          <p>
            The cards follow the actual table seating: top, left, right and bottom.
            The gold-bordered card is the <strong>selected player</strong>, whose position
            changes are being calculated.
          </p>
          <div className="delta-help-example">
            <div className="delta-help-example__row">
              <strong>Wind and name</strong>
              <span>Identify the player and their seat at the table.</span>
            </div>
            <div className="delta-help-example__row">
              <strong>Score and place</strong>
              <span>Current points and current standing (1st–4th).</span>
            </div>
            <div className="delta-help-example__row">
              <strong>Δ +8</strong>
              <span>Selected player&apos;s score minus this opponent&apos;s score.
                Green means ahead; red means behind.</span>
            </div>
          </div>
          <div className="delta-help-example">
            <div className="delta-help-example__row">
              <span className="delta-legend__item delta-action--self">Self</span>
              <span>On the selected player&apos;s card: minimum hand fan needed to reach a higher place by self-draw (tsumo).</span>
            </div>
            <div className="delta-help-example__row">
              <span className="delta-legend__item delta-action--direct">Win from this player</span>
              <span>On an opponent&apos;s card: minimum hand fan needed to reach the indicated place by winning on that opponent&apos;s discard.</span>
            </div>
            <div className="delta-help-example__row">
              <span className="delta-legend__item delta-action--loss">Deal into this player</span>
              <span>On an opponent&apos;s card: hand value at which discarding into that opponent&apos;s winning hand causes the indicated fall in standing.</span>
            </div>
          </div>
          <div className="delta-help-demo">
            <span className="delta-action delta-action--direct">8 → 2nd</span>
            <span>An 8-fan winning hand is enough to reach second place in this example.</span>
          </div>
          <div className="delta-help-demo">
            <span className="delta-action delta-action--loss">9 → 1st tie</span>
            <span>“Tie” means sharing a place with another player.</span>
          </div>
          <p>
            The numbers before the arrows are <strong>hand fan values</strong>, not changes
            to the player&apos;s total score. Each attainable place appears once.
            When both a sole place and a tied place are possible, the sole place
            takes priority, even if it requires a more expensive hand.
          </p>
          <button type="button" className="delta-back" onClick={onFallback}>Back to Delta</button>
        </section>
      </div>
    );
  }

  return (
    <div className="delta-annotation-layer" role="dialog" aria-modal="true" aria-label="How to read Delta">
      <button
        type="button"
        className="delta-annotation-layer__dismiss"
        onClick={onClose}
        aria-label="Close explanations"
      />
      <div className="delta-annotation-layer__top">
        <span>How to read Delta</span>
        <button type="button" onClick={onClose} aria-label="Close explanations">×</button>
      </div>
      <svg className="delta-annotation-lines" width="100%" height="100%" aria-hidden="true">
        {annotations.map((item) => {
          const tx = item.rect.left + item.rect.width / 2;
          const ty = item.rect.top + item.rect.height / 2;
          const lx = item.labelRect.left + item.labelRect.width / 2;
          const ly = item.labelRect.top + item.labelRect.height / 2;
          return (
            <line key={item.id} x1={lx} y1={ly} x2={tx} y2={ty}
              stroke="#fbbf24" strokeWidth="1.5" strokeDasharray="4 3" />
          );
        })}
      </svg>
      {annotations.map((item) => (
        <div key={item.id}>
          <div className={`delta-annotation-target delta-annotation-target--${item.id}`}
            style={{
              left: item.rect.left - 3, top: item.rect.top - 3,
              width: item.rect.width + 6, height: item.rect.height + 6,
            }}
          />
          <div className="delta-annotation-label"
            style={{
              left: item.labelRect.left, top: item.labelRect.top,
              width: item.labelRect.width, minHeight: item.labelRect.height,
            }}
          >{item.label}</div>
        </div>
      ))}
    </div>
  );
}

export default function Delta({
  isOpen,
  onClose,
  selectedPlayerId,
  players,
}: DeltaProps) {
  const [helpOpen, setHelpOpen] = useState(false);
  useEffect(() => {
    if (!isOpen) setHelpOpen(false);
  }, [isOpen]);

  // CurrentGame passes a newly created array on each render.
  const playersKey = JSON.stringify(
    players.map((p) => [p.id, p.name, p.score, p.position, p.wind])
  );

  const lastCalculations = useRef<ReturnType<typeof computeEmpty> | null>(null);
  const calculations = useMemo(() => {
    const empty = computeEmpty();
    if (!isOpen || selectedPlayerId === null) {
      return lastCalculations.current ?? empty;
    }
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
    rank,
    selectedPlayer,
    topPlayer,
    leftPlayer,
    rightPlayer,
    bottomPlayer,
    selfTargets,
    opponentTargets,
  } = calculations;

  return (
    <Modal.Backdrop
      isOpen={isOpen}
      onOpenChange={(open) => {
        if (!open) onClose();
      }}
    >
      <Modal.Container
        placement="center"
        size="full"
        scroll="inside"
        className="delta-modal"
      >
        <Modal.Dialog aria-label="Position thresholds" className="bg-surface text-foreground">
          <Modal.Header className="delta-modal__header relative border-b border-border pr-14">
            <div className="flex flex-col gap-1">
              <span>Position thresholds</span>
              {selectedPlayer && (
                <span className="text-sm font-normal text-muted">
                  {selectedPlayer.name} · {rank.first}
                  {ordinal(rank.first)}
                  {rank.last > rank.first
                    ? `–${rank.last}${ordinal(rank.last)} (tie)`
                    : ''}{' '}
                  place
                </span>
              )}
            </div>
            <Modal.CloseTrigger
              className="absolute right-3 top-1/2 z-20 -translate-y-1/2"
              aria-label="Close position thresholds"
            >
              <svg
                aria-hidden="true"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
              >
                <path d="M18 6 6 18M6 6l12 12" />
              </svg>
            </Modal.CloseTrigger>
          </Modal.Header>
          <Modal.Body className="delta-modal__body">
            {!selectedPlayer ? (
              <div className="flex min-h-48 items-center justify-center text-muted">
                Player not found
              </div>
            ) : (
              <div className="delta-layout">
                {([
                  ['top', topPlayer],
                  ['left', leftPlayer],
                  ['right', rightPlayer],
                  ['bottom', bottomPlayer],
                ] as const).map(([area, player]) => {
                  if (!player) return null;
                  if (player.id === selectedPlayer.id) {
                    return (
                      <section key={area}
                        data-delta-help="selected" className={`delta-selected delta-seat--${area} rounded-lg border-2 p-2 shadow-sm sm:p-3`}
                        aria-label={`Selected player: ${player.name}`}>
                        <div data-delta-help="identity" className="delta-identity border-b border-border pb-1.5">
                          <div className="delta-wind" aria-label={`Wind ${windLabel(player.wind)}`}>
                            {windLabel(player.wind)}
                          </div>
                          <div className="delta-person">
                            <div className="delta-name font-semibold" title={player.name}>{player.name}</div>
                            <div className="delta-score-line">
                              <span data-delta-help="score" className="delta-score tabular-nums text-muted">{player.score.toLocaleString('en-US')}</span>
                              <span className="delta-current-rank text-muted">
                                {rank.first}{ordinal(rank.first)}
                                {rank.last > rank.first ? `–${rank.last}${ordinal(rank.last)} (tie)` : ''}
                              </span>
                            </div>
                          </div>
                        </div>
                        {selfTargets.length > 0 ? (
                          <div data-delta-help="self" className="delta-action delta-action--self" title="Self-draw thresholds for higher places">
                            <TargetList targets={selfTargets} />
                          </div>
                        ) : (
                          <div className="delta-no-targets text-muted">No higher place by self-draw</div>
                        )}
                      </section>
                    );
                  }
                  const targets = opponentTargets[player.id];
                  return targets ? (
                    <OpponentCard key={area} player={player} selectedScore={selectedPlayer.score}
                      directTargets={targets.direct} discardTargets={targets.discard} area={area} />
                  ) : null;
                })}
                <div className="delta-center">
                  <button type="button" className="delta-back" onClick={onClose}>Back to table</button>
                  <button
                    type="button"
                    className="delta-help-trigger"
                    onClick={() => setHelpOpen(true)}
                    aria-label="How to read Delta thresholds"
                  >
                    <span className="delta-help-trigger__title">How to read</span>
                    <span className="delta-help-trigger__swatches" aria-hidden="true">
                      <span className="delta-legend__item delta-action--self">Self</span>
                      <span className="delta-legend__item delta-action--direct">Win</span>
                      <span className="delta-legend__item delta-action--loss">Deal in</span>
                    </span>
                  </button>
                </div>
                {helpOpen && (
                  <HelpOverlay
                    onClose={() => setHelpOpen(false)}
                    onFallback={() => setHelpOpen(false)}
                  />
                )}
              </div>
            )}
          </Modal.Body>
        </Modal.Dialog>
      </Modal.Container>
    </Modal.Backdrop>
  );
}
