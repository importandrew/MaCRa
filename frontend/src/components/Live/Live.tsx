import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';
import { useCallback, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useFetchLiveGameData } from '../../providers/game/gameDataApi';
import ContextMenu from '../ContextMenu/ContextMenu';
import RefreshButton from '../RefreshButton';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import AnimatedScore from '../Table/AnimatedScore';
import SeatName from '../Table/SeatName';
import TableLayout from '../Table/TableLayout';
import { useScoreAnimation } from '../Table/useScoreAnimation';
import LiveMenu from './Menu/LiveMenu';

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
  const { data, isLoading, isError, isForbidden, refetch } =
    useFetchLiveGameData(eventId, gameId);
  const [isContextMenuOpen, setIsContextMenuOpen] = useState(false);
  const [rotationOffset, setRotationOffset] = useState(0);
  const { previousScores, currentScores } = useScoreAnimation(
    data?.player_scores ?? null
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
        <TableLayout
          top={
            <>
              <SeatName player={top} />
              <span className="text-ellipsis overflow-hidden">
                <AnimatedScore
                  previousScores={previousScores}
                  currentScores={currentScores}
                  playerId={top.id}
                />
              </span>
            </>
          }
          left={
            <>
              <SeatName player={left} />
              <span className="text-ellipsis overflow-hidden">
                <AnimatedScore
                  previousScores={previousScores}
                  currentScores={currentScores}
                  playerId={left.id}
                />
              </span>
            </>
          }
          right={
            <>
              <SeatName player={right} />
              <span className="text-ellipsis overflow-hidden">
                <AnimatedScore
                  previousScores={previousScores}
                  currentScores={currentScores}
                  playerId={right.id}
                />
              </span>
            </>
          }
          bottom={
            <>
              <SeatName player={bot} />
              <span className="text-ellipsis overflow-hidden">
                <AnimatedScore
                  previousScores={previousScores}
                  currentScores={currentScores}
                  playerId={bot.id}
                />
              </span>
            </>
          }
          center={
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
          }
          refresh={
            <RefreshButton
              onRefresh={refetch}
              className="absolute z-1 left-0 bottom-0"
            />
          }
          menu={
            <ContextMenu
              isOpen={isContextMenuOpen}
              onOpen={() => setIsContextMenuOpen(true)}
              onClose={closeContextMenuOpen}
            >
              <LiveMenu
                gameId={data.id}
                eventId={eventId}
                onUndoSuccess={refetch}
                onStale={async () => (await refetch())?.log ?? null}
                log={data.log}
                onGameEnded={onGameEnded}
              />
            </ContextMenu>
          }
        />
      ) : null}
    </LoadingTemplate>
  );
}
