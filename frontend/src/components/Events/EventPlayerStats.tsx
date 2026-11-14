import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router-dom';
import { useEventLayout } from '../../providers/event/eventLayoutContext';
import { useFetchPlayerStats } from '../../providers/event/statsApi';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import PlayerStatsView from '../Stats/PlayerStatsView';

// EventPlayerStats is /event/:eventId/player/:playerId: a player's statistics
// in the event, opened from the rating. Others' statistics follow the
// rating's reach; the server answers 403 otherwise.
export default function EventPlayerStats() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { eventId, event } = useEventLayout();
  const { playerId: playerParam } = useParams();
  const playerId = Number(playerParam) || null;
  const { stats, isLoading, isError, isForbidden, isNotFound, refetch } =
    useFetchPlayerStats(eventId, playerId);

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      isForbidden={isForbidden}
      notFound={isNotFound || playerId == null ? 'player' : null}
      onRefresh={refetch}
      onGoHome={() => navigate(`/event/${eventId}/rating`)}
    >
      {stats && event && (
        <div className="flex flex-col w-full max-w-[600px] mx-auto px-3 md:px-0 md:my-3 gap-4 min-w-0">
          <div className="min-w-0">
            <h1 className="text-2xl font-semibold truncate" title={stats.name}>
              {stats.name}
            </h1>
            <p className="text-muted truncate">
              {t('stats.inEvent', { name: event.name })}
            </p>
          </div>
          <PlayerStatsView stats={stats} />
        </div>
      )}
    </LoadingTemplate>
  );
}
