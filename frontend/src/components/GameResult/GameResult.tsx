import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';
import { useCallback, useContext, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useRouteIds } from '../../hooks/useRouteIds';
import { AuthContext } from '../../providers/auth/authConstants';
import { useEvent } from '../../providers/event/eventConstants';
import { useEventCapabilities } from '../../providers/event/useEventCapabilities';
import {
  useFetchGameLogsData,
  useFetchGameResultData,
} from '../../providers/game/gameDataApi';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import LogsMenu from '../GameMenu/LogsMenu';
import GameLogs from './GameLogs';
import GameResultTable from './GameResultTable';

export default function GameResult() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { isAuthenticated } = useContext(AuthContext);
  // The result of a game in the URL (/event/N/game/M), or of the user's last
  // game in the selected event (/last).
  const route = useRouteIds();
  const { selectedEventId } = useEvent();
  const isEventGameRoute = route.eventId != null && route.gameId != null;
  const eventId = isEventGameRoute ? route.eventId : selectedEventId;
  const { hasCapability } = useEventCapabilities(eventId);
  const canCreateGame = hasCapability('create_game');

  const { data, isError, isForbidden, isLoading, fetchData } =
    useFetchGameResultData(
      isEventGameRoute && route.eventId && route.gameId
        ? { eventId: route.eventId, gameId: route.gameId }
        : undefined
    );
  const canUndoLastRound =
    hasCapability('undo_round') && Boolean(data?.can_undo_last_round);
  const [isChartVisible, setIsChartVisible] = useState(false);
  const [isLogsVisible, setIsLogsVisible] = useState(false);
  // The log is loaded once, when the chart or the log is first shown, and
  // serves both: a finished game's log doesn't change (an undo leaves the
  // page).
  const [isLogWanted, setIsLogWanted] = useState(false);
  const logs = useFetchGameLogsData(
    isLogWanted && data ? data.id : null,
    eventId
  );

  const showChart = () => {
    setIsChartVisible(!isChartVisible);
    setIsLogsVisible(false);
    setIsLogWanted(true);
  };
  const showLogs = () => {
    setIsLogsVisible(!isLogsVisible);
    setIsChartVisible(false);
    setIsLogWanted(true);
  };

  const onRefresh = useCallback(() => {
    fetchData();
  }, [fetchData]);

  const playAgain = () => {
    if (data) {
      const playerIds = data.player_results.map((player) => player.id);
      navigate('/new', { state: { playerIds } });
    }
  };

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      isForbidden={isForbidden}
      subject="game"
      onRefresh={onRefresh}
      onGoHome={() => navigate('/')}
    >
      {data ? (
        <div className="flex flex-col grow justify-center w-screen sm:max-w-[600px] h-full px-3 gap-3">
          <GameResultTable data={data} />
          <div className="flex flex-wrap gap-3 [&>*]:flex-1 [&>*]:basis-[calc(50%-0.375rem)] sm:[&>*]:basis-0">
            <Button
              type="button"
              onPress={showChart}
              variant="tertiary"
              fullWidth
            >
              {isChartVisible ? t('result.hideChart') : t('result.showChart')}
            </Button>
            <Button
              type="button"
              onPress={showLogs}
              variant="tertiary"
              fullWidth
            >
              {isLogsVisible ? t('result.hideLogs') : t('result.showLogs')}
            </Button>
            {!isEventGameRoute && isAuthenticated && canCreateGame && (
              <Button
                type="button"
                onPress={playAgain}
                variant="primary"
                fullWidth
              >
                {t('result.playAgain')}
              </Button>
            )}
          </div>
          {isLogsVisible && (
            <LogsMenu
              gameId={data.id}
              eventId={eventId}
              canUndoLastRound={canUndoLastRound}
              onUndoSuccess={() => {
                navigate('/now');
              }}
              logs={logs}
              onStale={logs.fetchData}
            />
          )}
          {isChartVisible && (
            <GameLogs isChartVisible={true} isLogVisible={false} logs={logs} />
          )}
        </div>
      ) : null}
    </LoadingTemplate>
  );
}
