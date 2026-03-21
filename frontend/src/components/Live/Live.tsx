import { Avatar, Button } from '@heroui/react';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useFetchLiveGameData } from '../../providers/game/gameDataApi';
import ContextMenu from '../ContextMenu/ContextMenu';
import { RefreshIcon } from '../Icons/RefreshIcon';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import AnimatedScore from '../Now/AnimatedScore';
import '../Now/Now.css';
import LiveMenu from './Menu/LiveMenu';

export default function Live() {
  const { id } = useParams<{ id: string }>();
  const gameId = id ? parseInt(id, 10) : 0;
  const { data, playersList, isLoading, isError, refetch } =
    useFetchLiveGameData(gameId);
  const [isContextMenuOpen, setIsContextMenuOpen] = useState(false);
  const [rotationOffset, setRotationOffset] = useState(0);
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

  const closeContextMenuOpen = () => {
    setIsContextMenuOpen(false);
  };

  const onRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

  const getRotatedIndex = (index: number): number => {
    const rotatedIndex = (index - rotationOffset + 4) % 4;
    return rotatedIndex;
  };

  const rotateCounterClockwise = () => {
    setRotationOffset((prev) => (prev + 1 + 4) % 4);
  };

  const rotateClockwise = () => {
    setRotationOffset((prev) => (prev - 1 + 4) % 4);
  };

  // Update currentScores when roundData changes, and save previous scores
  useEffect(() => {
    if (data) {
      // Save current scores to previous before updating
      if (currentScoresRef.current) {
        setPreviousScores(currentScoresRef.current);
      }

      // Update current scores with new data
      const newScores = {
        [data.player_scores[0].id]: data.player_scores[0].score,
        [data.player_scores[1].id]: data.player_scores[1].score,
        [data.player_scores[2].id]: data.player_scores[2].score,
        [data.player_scores[3].id]: data.player_scores[3].score,
      };
      currentScoresRef.current = newScores;
      setCurrentScores(newScores);
    }
  }, [data, rotationOffset]);

  const bot = data && data.player_scores[getRotatedIndex(0)];
  const right = data && data.player_scores[getRotatedIndex(1)];
  const top = data && data.player_scores[getRotatedIndex(2)];
  const left = data && data.player_scores[getRotatedIndex(3)];
  const hasNoData =
    !data ||
    !data.player_scores ||
    data.player_scores.length !== 4 ||
    playersList.length !== 4;

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      hasNoData={hasNoData}
      onRefresh={onRefresh}
    >
      {data && bot && right && top && left ? (
        <div className="flex justify-center h-full w-screen text-2xl">
          <div className="grid w-screen grid-cols-[25%_50%_25%] grid-rows-[20%_60%_20%] text-center">
            <div></div>
            <div className="relative self-center h-full w-full">
              <div className="table-top-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  {top.wind}{' '}
                  <Avatar
                    className="transition-transform"
                    color="default"
                    size="md"
                    name={top.name}
                  />{' '}
                  {top.name}
                </span>
                <span className="text-ellipsis overflow-hidden">
                  <AnimatedScore
                    previousScores={previousScores}
                    currentScores={currentScores}
                    playerId={top.id}
                  />
                </span>
              </div>
            </div>
            <div></div>
            <div className="relative self-center h-full w-full">
              <div className="table-right-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  {left.wind}{' '}
                  <Avatar
                    className="transition-transform"
                    color="default"
                    size="md"
                    name={left.name}
                  />{' '}
                  {left.name}
                </span>
                <span className="text-ellipsis overflow-hidden">
                  <AnimatedScore
                    previousScores={previousScores}
                    currentScores={currentScores}
                    playerId={left.id}
                  />
                </span>
              </div>
            </div>

            <div className="relative z-10 py-3 justify-self-center self-center items-center w-full text-5xl">
              <div className="flex flex-col items-center gap-2">
                <div className="p-4 h-16 text-5xl">{data.round.kanji}</div>
                <div className="flex gap-2">
                  <Button
                    size="sm"
                    variant="flat"
                    color="default"
                    type="button"
                    onPress={rotateClockwise}
                    className="text-xl"
                  >
                    ↻
                  </Button>
                  <Button
                    size="sm"
                    variant="flat"
                    color="default"
                    type="button"
                    onPress={rotateCounterClockwise}
                    className="text-xl"
                  >
                    ↺
                  </Button>
                </div>
              </div>
            </div>

            <div className="relative self-center h-full w-full">
              <div className="table-left-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  {right.wind}{' '}
                  <Avatar
                    className="transition-transform"
                    color="default"
                    size="md"
                    name={right.name}
                  />{' '}
                  {right.name}
                </span>
                <span className="text-ellipsis overflow-hidden">
                  <AnimatedScore
                    previousScores={previousScores}
                    currentScores={currentScores}
                    playerId={right.id}
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
                  {bot.wind}{' '}
                  <Avatar
                    className="transition-transform"
                    color="default"
                    size="md"
                    name={bot.name}
                  />{' '}
                  {bot.name}
                </span>
                <span className="text-ellipsis overflow-hidden">
                  <AnimatedScore
                    previousScores={previousScores}
                    currentScores={currentScores}
                    playerId={bot.id}
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
                <LiveMenu gameId={data.id} />
              </ContextMenu>
            </div>
          </div>
        </div>
      ) : null}
    </LoadingTemplate>
  );
}
