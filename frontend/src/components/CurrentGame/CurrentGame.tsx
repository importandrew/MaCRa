import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';
import { useEffect, useMemo, useRef, useState } from 'react';
import type { GameData } from '../../models/GameData';
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
import { getStoredEventId } from '../../providers/event/eventConstants';
import './CurrentGame.css';

export default function CurrentGame() {
  const { t } = useTranslation();
  const { data, isLoading, isError, refetch } = useFetchGameData();
  const [roundData, setRoundData] = useState<GameData | null>(null);
  const [isContextMenuOpen, setIsContextMenuOpen] = useState(false);
  const [isEntryOpen, setIsEntryOpen] = useState(false);
  const [isDeltaOpen, setIsDeltaOpen] = useState(false);
  const [deltaPlayerId, setDeltaPlayerId] = useState<number | null>(null);
  const openDelta = (playerId: number) => {
    setDeltaPlayerId(playerId);
    setIsDeltaOpen(true);
  };
  const entry = useRoundEntry(
    roundData,
    (next) => {
      setRoundData(next);
      setIsEntryOpen(false);
    },
    () => onRefresh()
  );
  const [currentScores, setCurrentScores] = useState<Record<number, number>>(
    {}
  );
  const currentScoresRef = useRef<Record<number, number>>({});
  const [previousScores, setPreviousScores] = useState<Record<number, number>>(
    {}
  );

  // The round entry clears after a submit (in useRoundEntry) and on refresh,
  // including undo, as agreed with the owner. Closing it with the kanji keeps
  // the choices.
  // A 404 for a table on screen means its game has ended or been aborted:
  // open the result (or home) instead of a loading error.
  const showGameGone = useShowGameGone();
  const onRefresh = async () => {
    entry.reset();
    setIsEntryOpen(false);
    const { data: fresh, status } = await refetch();
    const eventId = getStoredEventId();
    if (status === 404 && roundData && eventId) {
      await showGameGone(eventId, roundData.id);
    }
    return fresh;
  };

  useEffect(() => {
    if (data) {
      setRoundData(data);
    }
  }, [data]);

  // Update currentScores when roundData changes, and save previous scores
  useEffect(() => {
    if (roundData) {
      // Save current scores to previous before updating
      if (currentScoresRef.current) {
        setPreviousScores(currentScoresRef.current);
      }

      // Update current scores with new data
      const newScores = {
        [roundData.me.id]: roundData.me.score,
        [roundData.right.id]: roundData.right.score,
        [roundData.opp.id]: roundData.opp.score,
        [roundData.left.id]: roundData.left.score,
      };
      currentScoresRef.current = newScores;
      setCurrentScores(newScores);
    }
  }, [roundData]);

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
        <div className="flex justify-center h-full w-screen overflow-x-clip">
          <div className="grid w-screen grid-cols-[25%_50%_25%] grid-rows-[20%_60%_20%] text-center">
            <div></div>
            <div className="relative self-center h-full w-full">
              <div className="table-top-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
                <CurrentGameSeat
                  player={roundData.opp}
                  entry={entry}
                  isEntryOpen={isEntryOpen}
                  previousScores={previousScores}
                  currentScores={currentScores}
                  onOpenDelta={openDelta}
                />
              </div>
            </div>
            <div></div>
            <div className="relative self-center h-full w-full">
              <div className="table-right-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
                <CurrentGameSeat
                  player={roundData.left}
                  entry={entry}
                  isEntryOpen={isEntryOpen}
                  isSideSeat
                  previousScores={previousScores}
                  currentScores={currentScores}
                  onOpenDelta={openDelta}
                />
              </div>
            </div>

            {/* The center of the table: the round kanji, and the round entry
                under it while it's open. It must fit without scrolling. */}
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

            <div className="relative self-center h-full w-full">
              <div className="table-left-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
                <CurrentGameSeat
                  player={roundData.right}
                  entry={entry}
                  isEntryOpen={isEntryOpen}
                  isSideSeat
                  previousScores={previousScores}
                  currentScores={currentScores}
                  onOpenDelta={openDelta}
                />
              </div>
            </div>
            <div>
              <RefreshButton
                onRefresh={onRefresh}
                className="absolute z-1 left-0 bottom-0"
              />
            </div>
            <div className="relative self-center h-full w-full">
              <div className="table-bottom-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
                <CurrentGameSeat
                  player={roundData.me}
                  entry={entry}
                  isEntryOpen={isEntryOpen}
                  previousScores={previousScores}
                  currentScores={currentScores}
                  onOpenDelta={openDelta}
                />
              </div>
            </div>
            <div className="flex justify-end items-end">
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
            </div>
          </div>
          <Delta
            isOpen={isDeltaOpen}
            onClose={() => setIsDeltaOpen(false)}
            selectedPlayerId={deltaPlayerId}
            players={deltaPlayers}
            onSelectPlayer={setDeltaPlayerId}
          />
        </div>
      ) : null}
    </LoadingTemplate>
  );
}
