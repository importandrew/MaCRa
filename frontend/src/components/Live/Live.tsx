import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';
import UserAvatar from '../UserAvatar';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useFetchLiveGameData } from '../../providers/game/gameDataApi';
import ContextMenu from '../ContextMenu/ContextMenu';
import RefreshButton from '../RefreshButton';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import AnimatedScore from '../CurrentGame/AnimatedScore';
import '../CurrentGame/CurrentGame.css';
import LiveMenu from './Menu/LiveMenu';
import { useTableLogs } from '../CurrentGame/useTableLogs';

// onGameEnded is called once this user has ended the game (Finish): the game
// page reloads its status and shows the result.
export default function Live({ onGameEnded }: { onGameEnded: () => void }) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { eventId: eventIdParam, gameId: gameIdParam } = useParams<{
    eventId: string;
    gameId: string;
  }>();
  const eventId = eventIdParam ? parseInt(eventIdParam, 10) : 0;
  const gameId = gameIdParam ? parseInt(gameIdParam, 10) : 0;
  const { data, playersList, isLoading, isError, isForbidden, refetch } =
    useFetchLiveGameData(eventId, gameId);
  const logs = useTableLogs(data, refetch, eventId);
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

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      isForbidden={isForbidden}
      onRefresh={onRefresh}
      onGoHome={() => navigate('/')}
    >
      {data && bot && right && top && left ? (
        <div className="flex justify-center h-full w-screen">
          <div className="grid w-screen grid-cols-[25%_50%_25%] grid-rows-[20%_60%_20%] text-center">
            <div></div>
            <div className="relative self-center h-full w-full">
              <div className="table-top-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  <span className="shrink-0">{top.wind}</span>
                  <UserAvatar name={top.name} className="shrink-0" />
                  <span className="min-w-0 truncate">{top.name}</span>
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
                  <span className="shrink-0">{left.wind}</span>
                  <UserAvatar name={left.name} className="shrink-0" />
                  <span className="min-w-0 truncate">{left.name}</span>
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
                    variant="tertiary"
                    type="button"
                    aria-label={t('game.rotateClockwise')}
                    onPress={rotateClockwise}
                    className="text-xl"
                  >
                    ↻
                  </Button>
                  <Button
                    size="sm"
                    variant="tertiary"
                    type="button"
                    aria-label={t('game.rotateCounterclockwise')}
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
                  <span className="shrink-0">{right.wind}</span>
                  <UserAvatar name={right.name} className="shrink-0" />
                  <span className="min-w-0 truncate">{right.name}</span>
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
              <RefreshButton
                onRefresh={refetch}
                className="absolute z-1 left-0 bottom-0"
              />
            </div>
            <div className="relative self-center h-full w-full">
              <div className="table-bottom-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
                <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
                  <span className="shrink-0">{bot.wind}</span>
                  <UserAvatar name={bot.name} className="shrink-0" />
                  <span className="min-w-0 truncate">{bot.name}</span>
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
                onOpen={() => setIsContextMenuOpen(true)}
                onClose={closeContextMenuOpen}
              >
                <LiveMenu
                  gameId={data.id}
                  eventId={eventId}
                  onUndoSuccess={refetch}
                  logs={logs}
                  logSize={data.log_size}
                  onGameEnded={onGameEnded}
                />
              </ContextMenu>
            </div>
          </div>
        </div>
      ) : null}
    </LoadingTemplate>
  );
}
