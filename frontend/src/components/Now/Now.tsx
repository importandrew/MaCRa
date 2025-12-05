import { Button, Spinner } from '@heroui/react';
import { useCallback, useEffect, useState } from 'react';
import { useFetchGameData } from '../../providers/game/gameDataApi';
import ContextMenu from '../ContextMenu/ContextMenu';
import PlayerMenu from './Menu/PlayerMenu';
import TableMenu from './Menu/TableMenu';
import './Now.css';
import ScroreInput from './ScoreInput';

export default function Now() {
  const { data, playersList, isLoading, isError, refetch } =
    useFetchGameData(true);
  const [finedPlayerId, setFinedPlayerId] = useState<number | null>(null);
  const [isContextMenuOpen, setIsContextMenuOpen] = useState(false);

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
    if (finedPlayerId !== null) {
      openContextMenuOpen();
    }

    if (finedPlayerId === null && isContextMenuOpen) {
      closeContextMenuOpen();
      refetch();
    }
  }, [finedPlayerId]);

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

  if (!data) {
    return (
      <div className="flex flex-col h-full w-screen justify-center items-center gap-3">
        <p className="mb-4 text-center text-3xl">No data available</p>
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
              {data.opp.wind}{' '}
              {
                <PlayerMenu
                  player={playersList.find(
                    (player) => player.name === data.opp.name
                  )}
                  setFinedPlayerId={setFinedPlayerId}
                />
              }{' '}
              {data.opp.name}
            </span>
            <span className="text-ellipsis overflow-hidden">
              {data.opp.score}
            </span>
          </div>
        </div>
        <div></div>
        <div className="relative self-center h-full w-full">
          <div className="table-right-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
            <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
              {data.left.wind}{' '}
              {
                <PlayerMenu
                  player={playersList.find(
                    (player) => player.name === data.left.name
                  )}
                  setFinedPlayerId={setFinedPlayerId}
                />
              }{' '}
              {data.left.name}
            </span>
            <span className="text-ellipsis overflow-hidden">
              {data.left.score}
            </span>
          </div>
        </div>

        <ScroreInput data={data} playersList={playersList} refetch={refetch} />

        <div className="relative self-center h-full w-full">
          <div className="table-left-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
            <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
              {data.right.wind}{' '}
              {
                <PlayerMenu
                  player={playersList.find(
                    (player) => player.name === data.right.name
                  )}
                  setFinedPlayerId={setFinedPlayerId}
                />
              }{' '}
              {data.right.name}
            </span>
            <span className="text-ellipsis overflow-hidden">
              {data.right.score}
            </span>
          </div>
        </div>
        <div></div>
        <div className="relative self-center h-full w-full">
          <div className="table-bottom-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
            <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
              {data.me.wind}{' '}
              {
                <PlayerMenu
                  player={playersList.find(
                    (player) => player.name === data.me.name
                  )}
                  setFinedPlayerId={setFinedPlayerId}
                />
              }{' '}
              {data.me.name}
            </span>
            <span className="text-ellipsis overflow-hidden">
              {data.me.score}
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
              gameId={data.id}
              setFinedPlayerId={setFinedPlayerId}
            />
          </ContextMenu>
        </div>
      </div>
    </div>
  );
}
