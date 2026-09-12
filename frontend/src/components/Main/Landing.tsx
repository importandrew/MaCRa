import { Button, Spinner, useDisclosure } from '@heroui/react';
import { useContext } from 'react';
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

const iconClass = 'w-6 h-6 text-default-500';

export default function Landing() {
  const { isAuthenticated } = useContext(AuthContext);
  const { selectedEvent, selectedEventId, isLoading: eventLoading } =
    useEvent();
  const { data, isLoading } = useFetchPlayerStatus();
  const { hasCapability, isLoading: capsLoading } =
    useEventCapabilities(selectedEventId);
  const { isOpen, onOpen, onClose } = useDisclosure();

  const canCreateGame = hasCapability('create_game');
  const canRating = hasCapability('view_ratings');
  const canLiveGames = hasCapability('watch_live');
  const canFinishedGames = hasCapability('watch_results');

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
            className="text-sm text-default-500 font-sans max-w-md break-words"
            title={selectedEvent.name}
          >
            {selectedEvent.name}
          </Link>
        )}
        {showEventShortcuts &&
          (canRating || canLiveGames || canFinishedGames) && (
            <div className="flex items-center gap-1">
              {canRating && (
                <Button
                  as={Link}
                  to={`${eventBase}/rating`}
                  isIconOnly
                  variant="light"
                  aria-label="Rating"
                >
                  <RankingIcon className={iconClass} />
                </Button>
              )}
              {canLiveGames && (
                <Button
                  as={Link}
                  to={`${eventBase}/games/live`}
                  isIconOnly
                  variant="light"
                  aria-label="Live games"
                >
                  <EyeIcon className={iconClass} />
                </Button>
              )}
              {canFinishedGames && (
                <Button
                  as={Link}
                  to={`${eventBase}/games/finished`}
                  isIconOnly
                  variant="light"
                  aria-label="Finished games"
                >
                  <CheckInCircleIcon className={iconClass} />
                </Button>
              )}
            </div>
          )}
      </div>

      <div className="flex flex-col-reverse w-screen px-6 mb-6 gap-3">
        {!isAuthenticated && (
          <MainButton text="Login to play" link="/login" />
        )}

        {isAuthenticated && !selectedEventId && !eventLoading && (
          <MainButton text="Select event to play" onPress={onOpen} />
        )}

        {isAuthenticated &&
          selectedEventId &&
          (eventLoading || isLoading || capsLoading) && (
            <div className="flex justify-center">
              <Spinner label="Loading..." size="lg" />
            </div>
          )}

        {showAssistantShortcuts && (
          <>
            {data?.has_current && (
              <MainButton text="Current game" link="/now" />
            )}
            {!data?.has_current && canCreateGame && (
              <MainButton text="New game" link="/new" />
            )}
            {data?.has_last && (
              <MainButton text="Last game" color="secondary" link="/last" />
            )}
            {data?.has_last && (
              <MainButton
                text="Played games"
                color="default"
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
