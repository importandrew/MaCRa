import { Trans, useTranslation } from 'react-i18next';
import { Button, Card, Modal } from '@heroui/react';
import { useMemo, useState } from 'react';
import { STANDARD_RULES } from '../../models/GameRules';
import './Delta.css';

type TablePosition = 'top' | 'right' | 'bottom' | 'left';

export type DeltaPlayer = {
  id: number;
  name: string;
  wind?: string;
  score: number;
  position: TablePosition;
};

type DeltaProps = {
  isOpen: boolean;
  onClose: () => void;
  selectedPlayerId: number | null;
  onSelectPlayer: (playerId: number) => void;
  players: DeltaPlayer[];
  // minHand is the game's minimum hand value: no target is below it.
  minHand: number;
};

type PlaceTarget = { place: number; handValue: number; tiedThrough: number };
type LossTarget = {
  placesLost: number;
  finalPlace: number;
  tiedThrough: number;
  handValue: number;
};
type OpponentArea = 'top' | 'left' | 'right' | 'bottom';

type HandMode = 'self' | 'direct' | 'deal-in';
type Scenario = {
  mode: HandMode;
  selectedPlayerId: number;
  opponentId?: number;
};
type LinearScore = { id: number; base: number; slope: number };
type Rank = { first: number; last: number };

function rankOf(
  scores: LinearScore[],
  selectedId: number,
  handValue: number
): Rank {
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
function buildScenario(
  players: DeltaPlayer[],
  scenario: Scenario
): LinearScore[] {
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
function candidateHands(scores: LinearScore[], minHand: number): number[] {
  const values = new Set<number>([minHand]);
  for (let i = 0; i < scores.length; i++) {
    for (let j = i + 1; j < scores.length; j++) {
      const a = scores[i];
      const b = scores[j];
      if (a.slope === b.slope) continue;
      const crossing = (b.base - a.base) / (a.slope - b.slope);
      if (crossing < minHand - 1) continue;
      const floor = Math.floor(crossing);
      for (const hand of [floor, floor + 1, floor + 2]) {
        if (hand >= minHand && Number.isSafeInteger(hand)) values.add(hand);
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
  minHand: number,
  isRelevant: (rank: Rank) => boolean
): PlaceTarget[] {
  const candidates = new Map<number, PlaceTarget>();
  for (const handValue of candidateHands(scores, minHand)) {
    const rank = rankOf(scores, selectedPlayerId, handValue);
    if (!isRelevant(rank)) continue;
    const target: PlaceTarget = {
      place: rank.first,
      tiedThrough: rank.last,
      handValue,
    };
    const existing = candidates.get(target.place);
    if (
      !existing ||
      (existing.tiedThrough !== existing.place &&
        target.tiedThrough === target.place) ||
      ((existing.tiedThrough === existing.place) ===
        (target.tiedThrough === target.place) &&
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
  minHand: number,
  mode: 'self' | 'direct',
  discarderId?: number
): PlaceTarget[] {
  const initial = currentRank(players, selectedPlayerId);
  if (initial.first === 1 && initial.last === 1) return [];
  const scores = buildScenario(players, {
    mode,
    selectedPlayerId,
    opponentId: discarderId,
  });
  return selectTargets(
    scores,
    selectedPlayerId,
    minHand,
    (rank) =>
      rank.first < initial.first ||
      (rank.first === initial.first && rank.last < initial.last)
  );
}

function getDiscardTargets(
  players: DeltaPlayer[],
  selectedPlayerId: number,
  minHand: number,
  winnerId: number
): LossTarget[] {
  const initial = currentRank(players, selectedPlayerId);
  if (initial.first === players.length && initial.last === players.length)
    return [];
  const scores = buildScenario(players, {
    mode: 'deal-in',
    selectedPlayerId,
    opponentId: winnerId,
  });
  return selectTargets(
    scores,
    selectedPlayerId,
    minHand,
    (rank) =>
      rank.first > initial.first ||
      (rank.first === initial.first && rank.last > initial.last)
  ).map(({ place, tiedThrough, handValue }) => ({
    placesLost: Math.max(0, place - initial.first),
    finalPlace: place,
    tiedThrough,
    handValue,
  }));
}

// Keep every player at their original table seat, including the selected one.
function getTablePositions(players: DeltaPlayer[], selectedPlayerId: number) {
  return {
    selectedPlayer: players.find((p) => p.id === selectedPlayerId) ?? null,
    topPlayer: players.find((p) => p.position === 'top') ?? null,
    leftPlayer: players.find((p) => p.position === 'left') ?? null,
    rightPlayer: players.find((p) => p.position === 'right') ?? null,
    bottomPlayer: players.find((p) => p.position === 'bottom') ?? null,
  };
}

// A place is a number in a colored circle; an equals sign denotes a tie.
function PlaceBadge({
  place,
  tied = false,
}: {
  place: number;
  tied?: boolean;
}) {
  const { t } = useTranslation();
  return (
    <span
      className="delta-place-badge"
      aria-label={
        tied ? t('delta.tiedPlace', { place }) : t('delta.place', { place })
      }
    >
      {tied ? '=' : ''}
      {place}
    </span>
  );
}

type TargetListProps = { targets: PlaceTarget[] };
function TargetList({ targets }: TargetListProps) {
  return (
    <div className="delta-thresholds">
      {targets.map(({ place, handValue, tiedThrough }) => (
        <div className="delta-threshold" key={`${place}-${tiedThrough}`}>
          <strong className="delta-hand-value tabular-nums">{handValue}</strong>
          <span aria-hidden="true">→</span>
          <PlaceBadge place={place} tied={tiedThrough > place} />
        </div>
      ))}
    </div>
  );
}

function PlayerIdentity({
  player,
  selectedScore,
  rank,
  selected = false,
}: {
  player: DeltaPlayer;
  selectedScore?: number;
  rank?: Rank;
  selected?: boolean;
}) {
  const { t } = useTranslation();
  const gap = selectedScore === undefined ? null : selectedScore - player.score;
  return (
    <div className="delta-identity">
      <span
        className="delta-wind"
        aria-label={
          player.wind ? t('delta.wind', { wind: player.wind }) : undefined
        }
      >
        {player.wind}
      </span>
      <div className="delta-person">
        <div className="delta-name" title={player.name}>
          {player.name}
        </div>
        <div className="delta-score-line">
          <span className="delta-score tabular-nums">
            {player.score.toLocaleString('en-US')}
          </span>
          {selected && rank && (
            <span className="delta-current-rank">
              <PlaceBadge place={rank.first} tied={rank.last > rank.first} />
            </span>
          )}
          {gap !== null && (
            <span
              className={`delta-score-gap ${gap > 0 ? 'delta-score-gap--ahead' : gap < 0 ? 'delta-score-gap--behind' : ''}`}
              aria-label={t('delta.scoreDifferenceValue', { gap })}
            >
              Δ {gap > 0 ? '+' : ''}
              {gap.toLocaleString('en-US')}
            </span>
          )}
        </div>
      </div>
    </div>
  );
}

type OpponentCardProps = {
  player: DeltaPlayer;
  selectedScore: number;
  directTargets: PlaceTarget[];
  discardTargets: LossTarget[];
  area: OpponentArea;
  onSelect?: (id: number) => void;
};
function OpponentCard({
  player,
  selectedScore,
  directTargets,
  discardTargets,
  area,
  onSelect,
}: OpponentCardProps) {
  const { t } = useTranslation();
  const losses = discardTargets.map(
    ({ finalPlace, tiedThrough, handValue }) => ({
      place: finalPlace,
      tiedThrough,
      handValue,
    })
  );
  return (
    <Card
      className={`delta-player-card delta-opponent delta-seat--${area} rounded-lg p-2 gap-1.5 ${onSelect ? 'delta-card--clickable' : ''}`}
    >
      {onSelect && (
        <button
          type="button"
          className="delta-card-hitarea"
          aria-label={t('delta.showFor', { name: player.name })}
          onClick={() => onSelect(player.id)}
        />
      )}
      <Card.Header className="delta-card-header">
        <PlayerIdentity player={player} selectedScore={selectedScore} />
      </Card.Header>
      <Card.Content className="delta-card-content">
        {directTargets.length > 0 && (
          <div
            className="delta-action delta-action--direct"
            aria-label={t('delta.winOnDiscard')}
          >
            <TargetList targets={directTargets} />
          </div>
        )}
        {losses.length > 0 && (
          <div
            className="delta-action delta-action--loss"
            aria-label={t('delta.dealIntoWin')}
          >
            <TargetList targets={losses} />
          </div>
        )}
      </Card.Content>
    </Card>
  );
}

function SelectedCard({
  player,
  rank,
  selfTargets,
  area,
}: {
  player: DeltaPlayer;
  rank: Rank;
  selfTargets: PlaceTarget[];
  area: OpponentArea;
}) {
  const { t } = useTranslation();
  return (
    <Card
      className={`delta-player-card delta-selected delta-seat--${area} rounded-lg p-2 gap-1.5`}
    >
      <Card.Header className="delta-card-header">
        <PlayerIdentity player={player} rank={rank} selected />
      </Card.Header>
      <Card.Content className="delta-card-content">
        {selfTargets.length > 0 ? (
          <div
            className="delta-action delta-action--self"
            aria-label={t('delta.selfDrawThresholds')}
          >
            <TargetList targets={selfTargets} />
          </div>
        ) : (
          <div className="delta-no-targets">{t('delta.noHigherPlace')}</div>
        )}
      </Card.Content>
    </Card>
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

// The example is static relative to live play, but its thresholds use the real
// MCR calculation. Scores sum to zero, as they must in an MCR game. The names
// are placeholders: HelpExample shows them in the app's language.
const EXAMPLE_PLAYERS: DeltaPlayer[] = [
  { id: 101, name: 'North', wind: '北', score: -40, position: 'top' },
  { id: 102, name: 'West', wind: '西', score: 25, position: 'left' },
  { id: 103, name: 'South', wind: '南', score: 40, position: 'right' },
  { id: 104, name: 'You', wind: '東', score: -25, position: 'bottom' },
];
const EXAMPLE_SELECTED_ID = 104;
const EXAMPLE_SELECTED = EXAMPLE_PLAYERS[3];
// The static help example uses the standard minimum hand.
const EXAMPLE_MIN_HAND = STANDARD_RULES.min_hand;
const EXAMPLE_RANK = currentRank(EXAMPLE_PLAYERS, EXAMPLE_SELECTED_ID);
const EXAMPLE_SELF = getGainTargets(
  EXAMPLE_PLAYERS,
  EXAMPLE_SELECTED_ID,
  EXAMPLE_MIN_HAND,
  'self'
);
const EXAMPLE_WEST_WIN = getGainTargets(
  EXAMPLE_PLAYERS,
  EXAMPLE_SELECTED_ID,
  EXAMPLE_MIN_HAND,
  'direct',
  102
);
const EXAMPLE_SOUTH_LOSS = getDiscardTargets(
  EXAMPLE_PLAYERS,
  EXAMPLE_SELECTED_ID,
  EXAMPLE_MIN_HAND,
  103
);

function HelpExample({
  isOpen,
  onClose,
}: {
  isOpen: boolean;
  onClose: () => void;
}) {
  const { t } = useTranslation();
  const names = [
    t('winds.north'),
    t('winds.west'),
    t('winds.south'),
    t('delta.you'),
  ];
  const [north, west, south, you] = EXAMPLE_PLAYERS.map((player, index) => ({
    ...player,
    name: names[index],
  }));
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
        className="delta-help-modal"
      >
        <Modal.Dialog
          aria-label={t('delta.helpTitle')}
          className="bg-surface text-foreground"
        >
          <Modal.Header className="delta-modal__header border-b border-border">
            <span>{t('delta.helpTitle')}</span>
            <Modal.CloseTrigger
              aria-label={t('delta.closeHelp')}
              className="delta-help-close-trigger"
            >
              ×
            </Modal.CloseTrigger>
          </Modal.Header>
          <Modal.Body className="delta-help-body">
            <p className="delta-help-intro">{t('delta.helpIntro')}</p>
            <div className="delta-help-illustration">
              <div className="delta-help-canvas">
                <div className="delta-help-mini-table">
                  <OpponentCard
                    player={north}
                    selectedScore={EXAMPLE_SELECTED.score}
                    directTargets={[]}
                    discardTargets={[]}
                    area="top"
                  />
                  <OpponentCard
                    player={west}
                    selectedScore={EXAMPLE_SELECTED.score}
                    directTargets={EXAMPLE_WEST_WIN}
                    discardTargets={[]}
                    area="left"
                  />
                  <OpponentCard
                    player={south}
                    selectedScore={EXAMPLE_SELECTED.score}
                    directTargets={[]}
                    discardTargets={EXAMPLE_SOUTH_LOSS}
                    area="right"
                  />
                  <SelectedCard
                    player={you}
                    rank={EXAMPLE_RANK}
                    selfTargets={EXAMPLE_SELF}
                    area="bottom"
                  />
                </div>
                <svg
                  className="delta-help-leaders"
                  viewBox="0 0 320 410"
                  aria-hidden="true"
                >
                  <path d="M78 29 L118 55 L137 72" />
                  <path d="M246 29 L206 57 L160 76" />
                  <path d="M49 193 L73 193 L89 201" />
                  <path d="M270 193 L248 193 L229 201" />
                  <path d="M73 384 L112 346 L139 332" />
                  <path d="M251 384 L216 346 L180 332" />
                </svg>
                <span className="delta-help-label delta-help-label--one">
                  {t('delta.labelIdentity')}
                </span>
                <span className="delta-help-label delta-help-label--two">
                  {t('delta.scoreDifference')}
                </span>
                <span className="delta-help-label delta-help-label--three">
                  {t('delta.win')}
                </span>
                <span className="delta-help-label delta-help-label--four">
                  {t('delta.dealIn')}
                </span>
                <span className="delta-help-label delta-help-label--five">
                  {t('delta.selectedPlayer')}
                </span>
                <span className="delta-help-label delta-help-label--six">
                  {t('delta.selfDraw')}
                </span>
              </div>
            </div>
            <div className="delta-help-explanation">
              {(
                [
                  'helpSelf',
                  'helpWin',
                  'helpDealIn',
                  'helpDelta',
                  'helpPlaces',
                ] as const
              ).map((key) => (
                <p key={key}>
                  <Trans
                    i18nKey={`delta.${key}`}
                    components={{ b: <strong /> }}
                  />
                </p>
              ))}
            </div>
            <Button variant="secondary" onPress={onClose}>
              {t('delta.backToDelta')}
            </Button>
          </Modal.Body>
        </Modal.Dialog>
      </Modal.Container>
    </Modal.Backdrop>
  );
}

export default function Delta({
  isOpen,
  onClose,
  selectedPlayerId,
  onSelectPlayer,
  players,
  minHand,
}: DeltaProps) {
  const { t } = useTranslation();
  const [helpOpen, setHelpOpen] = useState(false);
  // close also closes the help, so it doesn't come back open next time.
  const close = () => {
    setHelpOpen(false);
    onClose();
  };
  const playersKey = JSON.stringify(
    players.map((p) => [p.id, p.name, p.score, p.position, p.wind])
  );
  // Computed whenever a player is selected, open or not: while the modal
  // closes (its exit animation) it keeps showing the same figures.
  const calculations = useMemo(() => {
    const empty = computeEmpty();
    if (selectedPlayerId === null) return empty;
    const positions = getTablePositions(players, selectedPlayerId);
    if (!positions.selectedPlayer) return empty;
    const opponentTargets: typeof empty.opponentTargets = {};
    for (const opponent of players) {
      if (opponent.id === selectedPlayerId) continue;
      opponentTargets[opponent.id] = {
        direct: getGainTargets(
          players,
          selectedPlayerId,
          minHand,
          'direct',
          opponent.id
        ),
        discard: getDiscardTargets(
          players,
          selectedPlayerId,
          minHand,
          opponent.id
        ),
      };
    }
    const result = {
      ...positions,
      rank: currentRank(players, selectedPlayerId),
      selfTargets: getGainTargets(players, selectedPlayerId, minHand, 'self'),
      opponentTargets,
    };
    return result;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedPlayerId, playersKey, minHand]);
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
        if (!open) close();
      }}
    >
      <Modal.Container
        placement="center"
        size="full"
        scroll="inside"
        className="delta-modal"
      >
        <Modal.Dialog
          aria-label={t('delta.title')}
          className="bg-surface text-foreground"
        >
          <Modal.Header className="delta-modal__header relative border-b border-border pr-14">
            <div className="flex flex-col gap-1">
              <span>{t('delta.title')}</span>
              {selectedPlayer && (
                <span className="delta-header-subtitle">
                  {selectedPlayer.name} ·{' '}
                  <PlaceBadge
                    place={rank.first}
                    tied={rank.last > rank.first}
                  />
                </span>
              )}
            </div>
            <Modal.CloseTrigger
              className="absolute right-3 top-1/2 z-20 -translate-y-1/2"
              aria-label={t('delta.close')}
            >
              ×
            </Modal.CloseTrigger>
          </Modal.Header>
          <Modal.Body className="delta-modal__body">
            {!selectedPlayer ? (
              <div className="flex min-h-48 items-center justify-center text-muted">
                {t('delta.playerNotFound')}
              </div>
            ) : (
              <div className="delta-layout">
                {(
                  [
                    ['top', topPlayer],
                    ['left', leftPlayer],
                    ['right', rightPlayer],
                    ['bottom', bottomPlayer],
                  ] as const
                ).map(([area, player]) => {
                  if (!player) return null;
                  if (player.id === selectedPlayer.id)
                    return (
                      <SelectedCard
                        key={area}
                        player={player}
                        rank={rank}
                        selfTargets={selfTargets}
                        area={area}
                      />
                    );
                  const targets = opponentTargets[player.id];
                  return targets ? (
                    <OpponentCard
                      key={area}
                      player={player}
                      selectedScore={selectedPlayer.score}
                      directTargets={targets.direct}
                      discardTargets={targets.discard}
                      area={area}
                      onSelect={onSelectPlayer}
                    />
                  ) : null;
                })}
                <div className="delta-center">
                  <Button variant="secondary" size="sm" onPress={close}>
                    {t('delta.backToTable')}
                  </Button>
                  <Button
                    variant="tertiary"
                    size="sm"
                    onPress={() => setHelpOpen(true)}
                  >
                    {t('delta.howToRead')}
                  </Button>
                  <div className="delta-help-swatches" aria-hidden="true">
                    <span className="delta-legend__item delta-action--self">
                      {t('delta.self')}
                    </span>
                    <span className="delta-legend__item delta-action--direct">
                      {t('delta.win')}
                    </span>
                    <span className="delta-legend__item delta-action--loss">
                      {t('delta.dealIn')}
                    </span>
                  </div>
                </div>
              </div>
            )}
          </Modal.Body>
          {helpOpen && (
            <HelpExample isOpen={helpOpen} onClose={() => setHelpOpen(false)} />
          )}
        </Modal.Dialog>
      </Modal.Container>
    </Modal.Backdrop>
  );
}
