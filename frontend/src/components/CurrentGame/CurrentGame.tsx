import { Button } from '@heroui/react';
import { useCallback, useEffect, useRef, useState } from 'react';
import { GameData } from '../../models/GameData';
import { useFetchGameData } from '../../providers/game/gameDataApi';
import ContextMenu from '../ContextMenu/ContextMenu';
import { RefreshIcon } from '../Icons/RefreshIcon';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import AnimatedScore from './AnimatedScore';
import PlayerMenu from './Menu/PlayerMenu';
import TableMenu from './Menu/TableMenu';
import './CurrentGame.css';
import ScroreInput from './ScoreInput';

export default function CurrentGame() {
  const { data, playersList, isLoading, isError, refetch } = useFetchGameData();
  const [roundData, setRoundData] = useState<GameData | null>(null);
  const [finedPlayerId, setFinedPlayerId] = useState<number | null>(null);
  const [isContextMenuOpen, setIsContextMenuOpen] = useState(false);
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

  const openContextMenuOpen = () => {
    setIsContextMenuOpen(true);
  };

  const closeContextMenuOpen = () => {
    setIsContextMenuOpen(false);
    setFinedPlayerId(null);
  };

  const onRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

  useEffect(() => {
    if (data) {
      setRoundData(data);
    }

    if (finedPlayerId !== null) {
      openContextMenuOpen();
    }

    if (finedPlayerId === null && isContextMenuOpen) {
      closeContextMenuOpen();
      refetch();
    }
  }, [finedPlayerId, data]);

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
      hasNoData={!roundData}
      onRefresh={onRefresh}
    >
      {roundData ? (
        <div className="flex justify-center h-full w-screen text-2xl">
          <div className="grid w-screen grid-cols-[25%_50%_25%] grid-rows-[20%_60%_20%] text-center">
            <div></div>
            <div className="relative self-center h-full w-full">
              <div className="table-top-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  {roundData.opp.wind}{' '}
                  {
                    <PlayerMenu
                      player={playersList.find(
                        (player) => player.name === roundData.opp.name
                      )}
                      setFinedPlayerId={setFinedPlayerId}
                    />
                  }{' '}
                  {roundData.opp.name}
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
                  {roundData.left.wind}{' '}
                  {
                    <PlayerMenu
                      player={playersList.find(
                        (player) => player.name === roundData.left.name
                      )}
                      setFinedPlayerId={setFinedPlayerId}
                    />
                  }{' '}
                  {roundData.left.name}
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

            <ScroreInput
              kanji={roundData.round.kanji}
              playersList={playersList}
              setRoundData={setRoundData}
            />

            <div className="relative self-center h-full w-full">
              <div className="table-left-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  {roundData.right.wind}{' '}
                  {
                    <PlayerMenu
                      player={playersList.find(
                        (player) => player.name === roundData.right.name
                      )}
                      setFinedPlayerId={setFinedPlayerId}
                    />
                  }{' '}
                  {roundData.right.name}
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
              <Button
                variant="light"
                radius="full"
                size="lg"
                isIconOnly
                onPress={refetch}
                className="absolute z-1 left-0 bottom-0"
              >
                <RefreshIcon className="text-gray-400" />
              </Button>
            </div>
            <div className="relative self-center h-full w-full">
              <div className="table-bottom-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  {roundData.me.wind}{' '}
                  {
                    <PlayerMenu
                      player={playersList.find(
                        (player) => player.name === roundData.me.name
                      )}
                      setFinedPlayerId={setFinedPlayerId}
                    />
                  }{' '}
                  {roundData.me.name}
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
                setIsOpen={setIsContextMenuOpen}
                onClose={closeContextMenuOpen}
              >
                <TableMenu
                  players={playersList}
                  finedPlayerId={finedPlayerId}
                  gameId={roundData.id}
                  setFinedPlayerId={setFinedPlayerId}
                  onRefetchNow={refetch}
                />
              </ContextMenu>
            </div>
          </div>
        </div>
      ) : null}
    </LoadingTemplate>
  );
}
