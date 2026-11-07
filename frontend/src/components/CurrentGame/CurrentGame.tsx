import { Button } from '@heroui/react';
import { useEffect, useRef, useState } from 'react';
import { GameData } from '../../models/GameData';
import { useFetchGameData } from '../../providers/game/gameDataApi';
import ContextMenu from '../ContextMenu/ContextMenu';
import RefreshButton from '../RefreshButton';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import AnimatedScore from './AnimatedScore';
import UserAvatar from '../UserAvatar';
import TableMenu from './Menu/TableMenu';
import { RoundEntryPanel, SeatMarks } from './RoundEntry';
import { useRoundEntry } from './useRoundEntry';
import './CurrentGame.css';

export default function CurrentGame() {
  const { data, playersList, isLoading, isError, refetch } = useFetchGameData();
  const [roundData, setRoundData] = useState<GameData | null>(null);
  const [isContextMenuOpen, setIsContextMenuOpen] = useState(false);
  const [isEntryOpen, setIsEntryOpen] = useState(false);
  const entry = useRoundEntry(
    roundData,
    (next) => {
      setRoundData(next);
      setIsEntryOpen(false);
    },
    () => void onRefresh()
  );
  const [currentScores, setCurrentScores] = useState<Record<number, number>>(
    {}
  );
  const currentScoresRef = useRef<Record<number, number>>({});
  const [previousScores, setPreviousScores] = useState<Record<number, number>>(
    () => {
      // Initialize with player IDs as keys and zeros as values
      if (playersList.length == 0) {
        return {};
      }
      return playersList.reduce(
        (acc, p) => ({ ...acc, [p.id]: 0 }),
        {} as Record<number, number>
      );
    }
  );

  // The round entry clears after a submit (in useRoundEntry) and on refresh,
  // including undo, as agreed with the owner. Closing it with the kanji keeps
  // the choices.
  const onRefresh = () => {
    entry.reset();
    setIsEntryOpen(false);
    return refetch();
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

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      onRefresh={onRefresh}
    >
      {roundData ? (
        <div className="flex justify-center h-full w-screen">
          <div className="grid w-screen grid-cols-[25%_50%_25%] grid-rows-[20%_60%_20%] text-center">
            <div></div>
            <div className="relative self-center h-full w-full">
              <div className="table-top-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
                {/* Toward the center, so the name stays at the screen edge. */}
                {isEntryOpen && (
                  <SeatMarks entry={entry} player={roundData.opp} />
                )}
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  <span className="shrink-0">{roundData.opp.wind}</span>
                  <UserAvatar name={roundData.opp.name} className="shrink-0" />
                  <span className="min-w-0 truncate">{roundData.opp.name}</span>
                </span>
                <span
                  className={`text-ellipsis overflow-hidden ${isEntryOpen ? 'hidden' : ''}`}
                >
                  <AnimatedScore
                    previousScores={previousScores}
                    currentScores={currentScores}
                    playerId={roundData.opp.id}
                  />
                </span>
              </div>
            </div>
            <div></div>
            <div className="relative self-center h-full w-full">
              <div className="table-right-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
                {/* Toward the center, so the name stays at the screen edge. */}
                {isEntryOpen && (
                  <SeatMarks entry={entry} player={roundData.left} isSideSeat />
                )}
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  <span className="shrink-0">{roundData.left.wind}</span>
                  <UserAvatar name={roundData.left.name} className="shrink-0" />
                  <span className="min-w-0 truncate">
                    {roundData.left.name}
                  </span>
                </span>
                <span
                  className={`text-ellipsis overflow-hidden ${isEntryOpen ? 'hidden' : ''}`}
                >
                  <AnimatedScore
                    previousScores={previousScores}
                    currentScores={currentScores}
                    playerId={roundData.left.id}
                  />
                </span>
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
                aria-label={`Enter result of round ${roundData.round.kanji}`}
                aria-expanded={isEntryOpen}
                onPress={() => setIsEntryOpen(!isEntryOpen)}
              >
                {roundData.round.kanji}
              </Button>
              {isEntryOpen && <RoundEntryPanel entry={entry} />}
            </div>

            <div className="relative self-center h-full w-full">
              <div className="table-left-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
                {/* Toward the center, so the name stays at the screen edge. */}
                {isEntryOpen && (
                  <SeatMarks
                    entry={entry}
                    player={roundData.right}
                    isSideSeat
                  />
                )}
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  <span className="shrink-0">{roundData.right.wind}</span>
                  <UserAvatar
                    name={roundData.right.name}
                    className="shrink-0"
                  />
                  <span className="min-w-0 truncate">
                    {roundData.right.name}
                  </span>
                </span>
                <span
                  className={`text-ellipsis overflow-hidden ${isEntryOpen ? 'hidden' : ''}`}
                >
                  <AnimatedScore
                    previousScores={previousScores}
                    currentScores={currentScores}
                    playerId={roundData.right.id}
                  />
                </span>
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
                {/* Toward the center, so the name stays at the screen edge. */}
                {isEntryOpen && (
                  <SeatMarks entry={entry} player={roundData.me} />
                )}
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  <span className="shrink-0">{roundData.me.wind}</span>
                  <UserAvatar name={roundData.me.name} className="shrink-0" />
                  <span className="min-w-0 truncate">{roundData.me.name}</span>
                </span>
                <span
                  className={`text-ellipsis overflow-hidden ${isEntryOpen ? 'hidden' : ''}`}
                >
                  <AnimatedScore
                    previousScores={previousScores}
                    currentScores={currentScores}
                    playerId={roundData.me.id}
                  />
                </span>
              </div>
            </div>
            <div className="flex justify-end items-end">
              <ContextMenu
                isOpen={isContextMenuOpen}
                onOpen={() => setIsContextMenuOpen(true)}
                onClose={() => setIsContextMenuOpen(false)}
              >
                <TableMenu gameId={roundData.id} onRefetchNow={onRefresh} />
              </ContextMenu>
            </div>
          </div>
        </div>
      ) : null}
    </LoadingTemplate>
  );
}
