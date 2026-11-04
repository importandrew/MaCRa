import { Button } from '@heroui/react';
import { useCallback, useEffect, useRef, useState } from 'react';
import { GameData } from '../../models/GameData';
import { useFetchGameData } from '../../providers/game/gameDataApi';
import ContextMenu from '../ContextMenu/ContextMenu';
import RefreshButton from '../RefreshButton';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import AnimatedScore from './AnimatedScore';
import PlayerMenu from './Menu/PlayerMenu';
import TableMenu from './Menu/TableMenu';
import RoundCard from './RoundCard';
import './CurrentGame.css';

export default function CurrentGame() {
  const { data, playersList, isLoading, isError, refetch } = useFetchGameData();
  const [roundData, setRoundData] = useState<GameData | null>(null);
  const [isContextMenuOpen, setIsContextMenuOpen] = useState(false);
  const [isRoundCardOpen, setIsRoundCardOpen] = useState(false);
  // roundCardKey remounts the round card, which clears its choices: after
  // a submit and on refresh (including undo), as agreed with the owner.
  const [roundCardKey, setRoundCardKey] = useState(0);
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

  const onRefresh = useCallback(() => {
    setRoundCardKey((key) => key + 1);
    return refetch();
  }, [refetch]);

  const onRoundSubmitted = (next: GameData) => {
    setRoundData(next);
    setIsRoundCardOpen(false);
    setRoundCardKey((key) => key + 1);
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
        <div className="flex justify-center h-full w-screen text-lg">
          <div className="grid w-screen grid-cols-[25%_50%_25%] grid-rows-[20%_60%_20%] text-center">
            <div></div>
            <div className="relative self-center h-full w-full">
              <div className="table-top-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  <span className="shrink-0">{roundData.opp.wind}</span>
                  {
                    <PlayerMenu
                      player={playersList.find(
                        (player) => player.id === roundData.opp.id
                      )}
                    />
                  }
                  <span className="min-w-0 truncate">{roundData.opp.name}</span>
                </span>
                <span className="text-ellipsis overflow-hidden">
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
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  <span className="shrink-0">{roundData.left.wind}</span>
                  {
                    <PlayerMenu
                      player={playersList.find(
                        (player) => player.id === roundData.left.id
                      )}
                    />
                  }
                  <span className="min-w-0 truncate">
                    {roundData.left.name}
                  </span>
                </span>
                <span className="text-ellipsis overflow-hidden">
                  <AnimatedScore
                    previousScores={previousScores}
                    currentScores={currentScores}
                    playerId={roundData.left.id}
                  />
                </span>
              </div>
            </div>

            <div className="relative z-10 py-3 justify-self-center self-center">
              <Button
                size="lg"
                variant="secondary"
                className="p-4 h-16 text-5xl"
                aria-label={`Enter result of round ${roundData.round.kanji}`}
                onPress={() => setIsRoundCardOpen(true)}
              >
                {roundData.round.kanji}
              </Button>
              <RoundCard
                key={roundCardKey}
                game={roundData}
                isOpen={isRoundCardOpen}
                onClose={() => setIsRoundCardOpen(false)}
                onSubmitted={onRoundSubmitted}
              />
            </div>

            <div className="relative self-center h-full w-full">
              <div className="table-left-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  <span className="shrink-0">{roundData.right.wind}</span>
                  {
                    <PlayerMenu
                      player={playersList.find(
                        (player) => player.id === roundData.right.id
                      )}
                    />
                  }
                  <span className="min-w-0 truncate">
                    {roundData.right.name}
                  </span>
                </span>
                <span className="text-ellipsis overflow-hidden">
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
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  <span className="shrink-0">{roundData.me.wind}</span>
                  {
                    <PlayerMenu
                      player={playersList.find(
                        (player) => player.id === roundData.me.id
                      )}
                    />
                  }
                  <span className="min-w-0 truncate">{roundData.me.name}</span>
                </span>
                <span className="text-ellipsis overflow-hidden">
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
