import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';
import { useMemo, useState } from 'react';
import { useFetchGameData } from '../../providers/game/gameDataApi';
import ContextMenu from '../ContextMenu/ContextMenu';
import RefreshButton from '../RefreshButton';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import Delta from './Delta';
import TableMenu from './Menu/TableMenu';
import { RoundEntryPanel } from './RoundEntry';
import CurrentGameSeat from './CurrentGameSeat';
import { useRoundEntry } from './useRoundEntry';
import { useShowGameGone } from './useShowGameGone';
import { useEvent } from '../../providers/event/eventConstants';
import TableLayout from '../Table/TableLayout';
import { useScoreAnimation } from '../Table/useScoreAnimation';

export default function CurrentGame() {
  const { t } = useTranslation();
  const {
    data: roundData,
    isLoading,
    isError,
    refetch,
    setData: setRoundData,
  } = useFetchGameData();
  const [isContextMenuOpen, setIsContextMenuOpen] = useState(false);
  const [isEntryOpen, setIsEntryOpen] = useState(false);
  const [isDeltaOpen, setIsDeltaOpen] = useState(false);
  const [deltaPlayerId, setDeltaPlayerId] = useState<number | null>(null);
  const openDelta = (playerId: number) => {
    setDeltaPlayerId(playerId);
    setIsDeltaOpen(true);
  };
  // A 404 for a table on screen means its game has ended or been aborted:
  // open the result (or home) instead of a loading error.
  const showGameGone = useShowGameGone();
  const { selectedEventId: eventId } = useEvent();
  // refreshTable reloads the table and returns it (null on an error).
  const refreshTable = async () => {
    setIsEntryOpen(false);
    const { data: fresh, status } = await refetch();
    if (status === 404 && roundData && eventId) {
      await showGameGone(eventId, roundData.id);
    }
    return fresh;
  };
  const entry = useRoundEntry(
    roundData,
    (next) => {
      setRoundData(next);
      setIsEntryOpen(false);
    },
    refreshTable
  );
  const { previousScores, currentScores } = useScoreAnimation(
    roundData
      ? [roundData.me, roundData.right, roundData.opp, roundData.left]
      : null
  );

  // Refresh (and an undo) clears the round entry, as agreed with the owner;
  // closing it with the kanji keeps the choices. A stale entry clears itself
  // (useRoundEntry).
  const onRefresh = () => {
    entry.reset();
    return refreshTable();
  };

  const deltaPlayers = useMemo(
    () =>
      roundData
        ? [
            { ...roundData.me, position: 'bottom' as const },
            { ...roundData.right, position: 'right' as const },
            { ...roundData.opp, position: 'top' as const },
            { ...roundData.left, position: 'left' as const },
          ]
        : [],
    [roundData]
  );

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      onRefresh={onRefresh}
    >
      {roundData ? (
        // overflow-x-clip: while a round is entered, a side seat (rotated,
        // with its +/− marks) is 88px across in a 80px column and sticks out
        // 4px on a 320px phone; clipping the avatar's edge beats a page that
        // scrolls sideways. Moving the seat inward would hit the entry panel.
        <>
          <TableLayout
            className="overflow-x-clip"
            top={
              <CurrentGameSeat
                player={roundData.opp}
                entry={entry}
                isEntryOpen={isEntryOpen}
                previousScores={previousScores}
                currentScores={currentScores}
                onOpenDelta={openDelta}
              />
            }
            left={
              <CurrentGameSeat
                player={roundData.left}
                entry={entry}
                isEntryOpen={isEntryOpen}
                isSideSeat
                previousScores={previousScores}
                currentScores={currentScores}
                onOpenDelta={openDelta}
              />
            }
            right={
              <CurrentGameSeat
                player={roundData.right}
                entry={entry}
                isEntryOpen={isEntryOpen}
                isSideSeat
                previousScores={previousScores}
                currentScores={currentScores}
                onOpenDelta={openDelta}
              />
            }
            bottom={
              <CurrentGameSeat
                player={roundData.me}
                entry={entry}
                isEntryOpen={isEntryOpen}
                previousScores={previousScores}
                currentScores={currentScores}
                onOpenDelta={openDelta}
              />
            }
            center={
              // The center of the table: the round kanji, and the round entry
              // under it while it's open. It must fit without scrolling.
              <div className="relative z-10 flex h-full min-h-0 w-full flex-col items-center justify-center gap-2 px-1">
                <Button
                  size="lg"
                  variant="secondary"
                  className={
                    isEntryOpen ? 'h-11 px-3 text-3xl' : 'p-4 h-16 text-5xl'
                  }
                  aria-label={t('game.enterRound', {
                    round: roundData.round.kanji,
                  })}
                  aria-expanded={isEntryOpen}
                  onPress={() => setIsEntryOpen(!isEntryOpen)}
                >
                  {roundData.round.kanji}
                </Button>
                {isEntryOpen && <RoundEntryPanel entry={entry} />}
              </div>
            }
            refresh={
              <RefreshButton
                onRefresh={onRefresh}
                className="absolute z-1 left-0 bottom-0"
              />
            }
            menu={
              <ContextMenu
                isOpen={isContextMenuOpen}
                onOpen={() => setIsContextMenuOpen(true)}
                onClose={() => setIsContextMenuOpen(false)}
              >
                <TableMenu
                  gameId={roundData.id}
                  onRefetchNow={onRefresh}
                  onStale={async () => (await onRefresh())?.log ?? null}
                  log={roundData.log}
                />
              </ContextMenu>
            }
          />
          <Delta
            isOpen={isDeltaOpen}
            onClose={() => setIsDeltaOpen(false)}
            selectedPlayerId={deltaPlayerId}
            players={deltaPlayers}
            minHand={roundData.rules.min_hand}
            onSelectPlayer={setDeltaPlayerId}
          />
        </>
      ) : null}
    </LoadingTemplate>
  );
}
