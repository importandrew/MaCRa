import { Avatar, Button, Spinner } from '@heroui/react';
import { useCallback, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useFetchLiveGameData } from '../../providers/game/gameDataApi';
import ContextMenu from '../ContextMenu/ContextMenu';
import { RefreshIcon } from '../Icons/RefreshIcon';
import AnimatedScore from '../Now/AnimatedScore';
import '../Now/Now.css';
import LiveMenu from './Menu/LiveMenu';

export default function Live() {
  const { id } = useParams<{ id: string }>();
  const gameId = id ? parseInt(id, 10) : 0;
  const { data, isLoading, isError, refetch } = useFetchLiveGameData(gameId);
  const [isContextMenuOpen, setIsContextMenuOpen] = useState(false);
  const [rotationOffset, setRotationOffset] = useState(0);

  const closeContextMenuOpen = () => {
    setIsContextMenuOpen(false);
  };

  const onRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

  const getPlayer = (index: number) => {
    if (!data || !data.player_scores || data.player_scores.length !== 4) {
      return null;
    }

    const rotatedIndex = (index - rotationOffset + 4) % 4;
    return data.player_scores[rotatedIndex];
  };

  const rotateCounterClockwise = () => {
    setRotationOffset((prev) => (prev + 1 + 4) % 4);
  };

  const rotateClockwise = () => {
    setRotationOffset((prev) => (prev - 1 + 4) % 4);
  };

  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-screen my-auto justify-center items-center">
        <Spinner label="Loading..." size="lg" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="flex w-screen h-full justify-center items-center">
        <Button type="button" onPress={onRefresh}>
          Refresh
        </Button>
      </div>
    );
  }

  if (!data || !data.player_scores || data.player_scores.length !== 4) {
    return (
      <div className="flex flex-col h-full w-screen justify-center items-center gap-3">
        <p className="mb-4 text-center text-3xl">No data available</p>
        <Button type="button" onPress={onRefresh}>
          Refresh
        </Button>
      </div>
    );
  }

  const east = getPlayer(0);
  const south = getPlayer(1);
  const west = getPlayer(2);
  const north = getPlayer(3);

  if (!east || !south || !west || !north) {
    return (
      <div className="flex flex-col h-full w-screen justify-center items-center gap-3">
        <p className="mb-4 text-center text-3xl">Invalid game data</p>
        <Button type="button" onPress={onRefresh}>
          Refresh
        </Button>
      </div>
    );
  }

  return (
    <div className="flex justify-center h-full w-screen text-2xl">
      <div className="grid w-screen grid-cols-[25%_50%_25%] grid-rows-[20%_60%_20%] text-center">
        <div></div>
        <div className="relative self-center h-full w-full">
          <div className="table-top-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
            <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
              {west.wind}{' '}
              <Avatar
                className="transition-transform"
                color="default"
                size="md"
                name={west.name}
              />{' '}
              {west.name}
            </span>
            <span className="text-ellipsis overflow-hidden">
              <AnimatedScore score={west.score} />
            </span>
          </div>
        </div>
        <div></div>
        <div className="relative self-center h-full w-full">
          <div className="table-right-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
            <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
              {north.wind}{' '}
              <Avatar
                className="transition-transform"
                color="default"
                size="md"
                name={north.name}
              />{' '}
              {north.name}
            </span>
            <span className="text-ellipsis overflow-hidden">
              <AnimatedScore score={north.score} />
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
              {south.wind}{' '}
              <Avatar
                className="transition-transform"
                color="default"
                size="md"
                name={south.name}
              />{' '}
              {south.name}
            </span>
            <span className="text-ellipsis overflow-hidden">
              <AnimatedScore score={south.score} />
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
              {east.wind}{' '}
              <Avatar
                className="transition-transform"
                color="default"
                size="md"
                name={east.name}
              />{' '}
              {east.name}
            </span>
            <span className="text-ellipsis overflow-hidden">
              <AnimatedScore score={east.score} />
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
  );
}
