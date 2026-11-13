import { useOverlayState } from '@heroui/react';
import { useContext, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { AuthContext } from '../../providers/auth/authConstants';
import { useEvent } from '../../providers/event/eventConstants';
import { useEventCapabilities } from '../../providers/event/useEventCapabilities';
import { useFetchPlayerStatus } from '../../providers/game/gameDataApi';
import SelectEventModal from '../Events/SelectEventModal';
import { CheckInCircleIcon } from '../Icons/CheckInCircleIcon';
import { EyeIcon } from '../Icons/EyeIcon';
import { RankingIcon } from '../Icons/RankingIcon';
import MainButton from './MainButton';
import { LinkButton } from '../LinkButton';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';

const iconClass = 'w-6 h-6 text-muted';

export default function Landing() {
  const { t } = useTranslation();
  const { isAuthenticated } = useContext(AuthContext);
  const {
    selectedEvent,
    selectedEventId,
    isLoading: eventLoading,
    dropInaccessibleEvent,
  } = useEvent();
  const { data, isLoading } = useFetchPlayerStatus();
  const {
    me,
    hasCapability,
    isLoading: capsLoading,
    isForbidden: capsForbidden,
  } = useEventCapabilities(selectedEventId);
  const { isOpen, open: onOpen, close: onClose } = useOverlayState();

  const canCreateGame = hasCapability('create_game');
  const canRating = hasCapability('view_ratings');
  const canLiveGames = hasCapability('watch_live');
  const canFinishedGames = hasCapability('watch_results');

  // The user has lost access to the selected event (e.g. removed from a
  // private one): clear it, so "Select event to play" shows instead of an
  // empty page.
  useEffect(() => {
    if (capsForbidden && selectedEventId != null) dropInaccessibleEvent();
  }, [capsForbidden, selectedEventId, dropInaccessibleEvent]);

  // noGameReason says why there are no game buttons, in their place: nothing
  // to do in the selected event.
  const hasGameButton = Boolean(
    data?.has_current || canCreateGame || data?.has_last
  );
  let noGameReason: string | null = null;
  if (!hasGameButton && selectedEvent) {
    if (selectedEvent.finished_at) {
      noGameReason = t('landing.eventFinished');
    } else if (me && !me.roles.includes('player')) {
      noGameReason = t('landing.notAPlayer');
    }
  }

  const showAssistantShortcuts =
    isAuthenticated && selectedEvent && !isLoading && !capsLoading && data;

  const showEventShortcuts =
    isAuthenticated && selectedEvent && selectedEventId != null && !capsLoading;

  const eventBase = selectedEventId != null ? `/event/${selectedEventId}` : '';

  return (
    <div className="flex flex-col grow h-full">
      <div className="flex flex-col grow w-screen px-16 font-mono justify-center items-center text-center gap-4">
        <span className={'text-5xl'}>MaCRa</span>
        {selectedEvent && selectedEventId != null && (
          <Link
            to={eventBase}
            className="text-sm text-muted font-sans max-w-md break-words"
            title={selectedEvent.name}
          >
            {selectedEvent.name}
          </Link>
        )}
        {showEventShortcuts &&
          (canRating || canLiveGames || canFinishedGames) && (
            <div className="flex items-center gap-1">
              {canRating && (
                <LinkButton
                  to={`${eventBase}/rating`}
                  isIconOnly
                  variant="ghost"
                  aria-label={t('landing.rating')}
                >
                  <RankingIcon className={iconClass} />
                </LinkButton>
              )}
              {canLiveGames && (
                <LinkButton
                  to={`${eventBase}/games/live`}
                  isIconOnly
                  variant="ghost"
                  aria-label={t('landing.liveGames')}
                >
                  <EyeIcon className={iconClass} />
                </LinkButton>
              )}
              {canFinishedGames && (
                <LinkButton
                  to={`${eventBase}/games/finished`}
                  isIconOnly
                  variant="ghost"
                  aria-label={t('landing.finishedGames')}
                >
                  <CheckInCircleIcon className={iconClass} />
                </LinkButton>
              )}
            </div>
          )}
      </div>

      <div className="flex flex-col-reverse w-screen px-6 mb-6 gap-3">
        {!isAuthenticated && (
          <MainButton text={t('landing.loginToPlay')} link="/login" />
        )}

        {isAuthenticated && !selectedEventId && !eventLoading && (
          <MainButton text={t('landing.selectEventToPlay')} onPress={onOpen} />
        )}

        {isAuthenticated &&
          selectedEventId &&
          (eventLoading || isLoading || capsLoading) && (
            <div className="flex justify-center">
              <LoadingSpinner label={t('landing.loading')} />
            </div>
          )}

        {showAssistantShortcuts && (
          <>
            {noGameReason && (
              <p className="text-center text-muted">{noGameReason}</p>
            )}
            {data?.has_current && (
              <MainButton text={t('landing.currentGame')} link="/now" />
            )}
            {!data?.has_current && canCreateGame && (
              <MainButton text={t('landing.newGame')} link="/new" />
            )}
            {data?.has_last && (
              <MainButton
                text={t('landing.lastGame')}
                variant="secondary"
                link="/last"
              />
            )}
            {data?.has_last && (
              <MainButton
                text={t('landing.playedGames')}
                variant="tertiary"
                link="/history"
              />
            )}
          </>
        )}
      </div>

      <SelectEventModal isOpen={isOpen} onClose={onClose} />
    </div>
  );
}
