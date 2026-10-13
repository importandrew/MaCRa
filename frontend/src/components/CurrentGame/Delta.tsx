import { Modal, ModalBody, ModalContent, ModalHeader } from '@heroui/react';
import './Delta.css';

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
type LossTarget = { placesLost: number; finalPlace: number; handValue: number };
type OpponentArea = 'top' | 'left' | 'right';

const MIN_HAND = 8;

function sortPlayers(players: DeltaPlayer[]): DeltaPlayer[] {
  // Keep the input order as the tie-breaker until the project defines
  // a dedicated rule for players with equal scores.
  return [...players].sort((a, b) => b.score - a.score);
}

function validHand(required: number): number {
  return Math.max(MIN_HAND, required);
}

function getRequiredSelfHandAgainstPlayer(selectedScore: number, targetScore: number): number {
  return validHand(Math.ceil((targetScore - selectedScore + 1) / 4 - 8));
}

function getRequiredDirectHandAgainstPlayer(
  selectedScore: number,
  targetScore: number,
  isDiscarder: boolean
): number {
  const difference = targetScore - selectedScore + 1;
  return validHand(
    isDiscarder ? Math.ceil((difference - 32) / 2) : Math.ceil(difference - 32)
  );
}

function getPlayersToOvertake(
  sortedPlayers: DeltaPlayer[],
  selectedPlayerId: number,
  targetPlace: number
): DeltaPlayer[] | null {
  const selectedIndex = sortedPlayers.findIndex((p) => p.id === selectedPlayerId);
  if (selectedIndex < 0 || targetPlace < 1 || targetPlace > selectedIndex) return null;
  return sortedPlayers.slice(targetPlace - 1, selectedIndex);
}

function getGainTargets(
  players: DeltaPlayer[],
  selectedPlayerId: number,
  mode: 'self' | 'direct',
  discarderId?: number
): PlaceTarget[] {
  const sorted = sortPlayers(players);
  const selectedIndex = sorted.findIndex((p) => p.id === selectedPlayerId);
  if (selectedIndex <= 0) return [];
  const selected = sorted[selectedIndex];
  const targets: PlaceTarget[] = [];

  for (let place = selectedIndex; place >= 1; place--) {
    const toOvertake = getPlayersToOvertake(sorted, selectedPlayerId, place);
    if (!toOvertake) continue;
    const handValue = Math.max(
      MIN_HAND,
      ...toOvertake.map((p) =>
        mode === 'self'
          ? getRequiredSelfHandAgainstPlayer(selected.score, p.score)
          : getRequiredDirectHandAgainstPlayer(
              selected.score,
              p.score,
              p.id === discarderId
            )
      )
    );
    targets.push({ place, handValue });
  }
  return targets;
}

function getDiscardTargets(
  players: DeltaPlayer[],
  selectedPlayerId: number,
  winnerId: number
): LossTarget[] {
  const sorted = sortPlayers(players);
  const selectedIndex = sorted.findIndex((p) => p.id === selectedPlayerId);
  if (selectedIndex < 0 || selectedIndex === sorted.length - 1) return [];

  const selected = sorted[selectedIndex];
  const currentPlace = selectedIndex + 1;
  const thresholds = sorted.slice(selectedIndex + 1).map((target) => {
    const difference = selected.score - target.score + 1;
    const required = target.id === winnerId
      ? Math.ceil((difference - 32) / 2)
      : Math.ceil(difference);
    return validHand(required);
  });

  const uniqueValues = [...new Set(thresholds)].sort((a, b) => a - b);
  return uniqueValues.map((handValue) => {
    const placesLost = thresholds.filter((threshold) => threshold <= handValue).length;
    return {
      placesLost,
      finalPlace: currentPlace + placesLost,
      handValue,
    };
  });
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
    return <div className={`${className}__empty`}>{emptyText}</div>;
  }
  return (
    <div className={`${className}__targets`}>
      {targets.map(({ place, handValue }) => (
        <div key={place} className={`${className}__target`}>
          <span className={`${className}__value`}>{handValue}</span>
          <span className={`${className}__destination`}>→ {place}-е</span>
        </div>
      ))}
    </div>
  );
}

type OpponentCardProps = {
  player: DeltaPlayer;
  selectedPlayer: DeltaPlayer;
  players: DeltaPlayer[];
  area: OpponentArea;
};

function OpponentCard({ player, selectedPlayer, players, area }: OpponentCardProps) {
  const directTargets = getGainTargets(players, selectedPlayer.id, 'direct', player.id);
  const discardTargets = getDiscardTargets(players, selectedPlayer.id, player.id);

  return (
    <section className={`delta-opponent delta-opponent--${area}`}>
      <div className="delta-opponent__identity">
        <div className="delta-opponent__name" title={player.name}>{player.name}</div>
        <div className={[
          'delta-opponent__score',
          player.score < 0 ? 'delta-opponent__score--negative' : '',
        ].filter(Boolean).join(' ')}>
          {player.score.toLocaleString('ru-RU')}
        </div>
      </div>

      {directTargets.length > 0 && (
        <div className="delta-direct">
          <div className="delta-direct__title" title={`Прямой выигрыш с ${player.name}`}>
            🎯 Взять с этого игрока
          </div>
          <TargetList targets={directTargets} emptyText="" className="delta-direct" />
        </div>
      )}

      {discardTargets.length > 0 && (
        <div className="delta-discard">
          <div className="delta-discard__title" title={`Набросить игроку ${player.name}`}>
            💥 Набросить этому игроку
          </div>
          <div className="delta-discard__targets">
            {discardTargets.map(({ finalPlace, handValue, placesLost }) => (
              <div key={finalPlace} className="delta-discard__target">
                <span className="delta-discard__value">{handValue}</span>
                <span className="delta-discard__destination">→ {finalPlace}-е</span>
                <span className="delta-discard__arrows" aria-label={`Потеряно мест: ${placesLost}`}>
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
  const sorted = sortPlayers(players);
  const selectedIndex = sorted.findIndex((p) => p.id === selectedPlayerId);
  const { selectedPlayer, topPlayer, leftPlayer, rightPlayer } = selectedPlayerId === null
    ? { selectedPlayer: null, topPlayer: null, leftPlayer: null, rightPlayer: null }
    : getRelativePositions(players, selectedPlayerId);
  const selfTargets = selectedPlayer
    ? getGainTargets(players, selectedPlayer.id, 'self')
    : [];

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      size="full"
      placement="center"
      scrollBehavior="inside"
      classNames={{
        base: 'delta-modal',
        header: 'delta-modal__header',
        body: 'delta-modal__body',
      }}
    >
      <ModalContent>
        {() => (
          <>
            <ModalHeader>
              <div className="delta-title">
                <span>Дельты по местам</span>
                {selectedPlayer && (
                  <span className="delta-title__player">
                    {selectedPlayer.name} · {selectedIndex + 1}-е место
                  </span>
                )}
              </div>
            </ModalHeader>
            <ModalBody onClick={(event) => {
              const target = event.target as HTMLElement;
              if (!target.closest('.delta-opponent, .delta-selected, .delta-legend')) {
                onClose();
              }
            }}>
              {!selectedPlayer ? (
                <div className="delta-empty">Игрок не найден</div>
              ) : (
                <div className="delta-layout">
                  {topPlayer && (
                    <OpponentCard player={topPlayer} selectedPlayer={selectedPlayer}
                      players={players} area="top" />
                  )}
                  {leftPlayer && (
                    <OpponentCard player={leftPlayer} selectedPlayer={selectedPlayer}
                      players={players} area="left" />
                  )}
                  <section className="delta-selected">
                    <div className="delta-selected__header">
                      <span className="delta-selected__icon">👤</span>
                      <div className="delta-selected__identity">
                        <div className="delta-selected__name" title={selectedPlayer.name}>
                          {selectedPlayer.name}
                        </div>
                        <div className={[
                          'delta-selected__score',
                          selectedPlayer.score < 0 ? 'delta-selected__score--negative' : '',
                        ].filter(Boolean).join(' ')}>
                          {selectedPlayer.score.toLocaleString('ru-RU')}
                        </div>
                        <div className="delta-selected__place">
                          Сейчас {selectedIndex + 1}-е место
                        </div>
                      </div>
                    </div>
                    {selfTargets.length > 0 && (
                      <div className="delta-self">
                        <div className="delta-self__title">🀄 Селф (со стены)</div>
                        <TargetList targets={selfTargets} emptyText="" className="delta-self" />
                      </div>
                    )}
                  </section>
                  {rightPlayer && (
                    <OpponentCard player={rightPlayer} selectedPlayer={selectedPlayer}
                      players={players} area="right" />
                  )}
                  <div className="delta-legend">
                    <span>🀄 минимальная рука при селфе</span>
                    <span>🎯 минимальная рука с указанного игрока</span>
                    <span>→ итоговое место</span>
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
