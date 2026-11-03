import { Button, Drawer, Spinner, useOverlayState } from '@heroui/react';
import { useCallback, useMemo, useState } from 'react';
import { Outlet, useNavigate, useParams } from 'react-router-dom';
import { hasCapability as checkCapability } from '../../providers/event/useEventCapabilities';
import {
  finishEvent,
  useFetchEvent,
  useFetchEventMe,
} from '../../providers/event/eventApi';
import {
  EventLayoutContext,
  EventLayoutContextValue,
} from '../../providers/event/eventLayoutContext';
import Confirm from '../Confirm';
import ContextMenu from '../ContextMenu/ContextMenu';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import { showToast } from '../Toast';

export default function EventLayout() {
  const { eventId: eventIdParam } = useParams();
  const eventId = eventIdParam ? parseInt(eventIdParam, 10) : null;
  const navigate = useNavigate();
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const [isFinishing, setIsFinishing] = useState(false);
  const {
    isOpen: isFinishConfirmOpen,
    open: onFinishConfirmOpen,
    close: onFinishConfirmClose,
  } = useOverlayState();

  const {
    data: event,
    isLoading: eventLoading,
    isError: eventError,
    isForbidden: eventForbidden,
    refetch: refetchEvent,
  } = useFetchEvent(eventId);
  const {
    data: me,
    isLoading: meLoading,
    isError: meError,
    isForbidden: meForbidden,
    refetch: refetchMe,
  } = useFetchEventMe(eventId);

  const isLoading = eventLoading || meLoading;
  const isError = eventError || meError;
  const isForbidden = eventForbidden || meForbidden;

  const refetch = useCallback(() => {
    refetchEvent();
    refetchMe();
  }, [refetchEvent, refetchMe]);

  const hasCapability = useCallback(
    (cap: Parameters<EventLayoutContextValue['hasCapability']>[0]) =>
      checkCapability(me, cap),
    [me]
  );

  const contextValue = useMemo<EventLayoutContextValue | null>(() => {
    if (eventId == null) return null;
    return {
      eventId,
      event,
      me,
      isLoading,
      isError,
      isForbidden,
      refetch,
      hasCapability,
    };
  }, [
    eventId,
    event,
    me,
    isLoading,
    isError,
    isForbidden,
    refetch,
    hasCapability,
  ]);

  const canManage =
    hasCapability('manage_event') || hasCapability('manage_members');
  const canManageEvent = hasCapability('manage_event');
  const canFinishEvent = canManageEvent && !event?.finished_at;
  const canRating = hasCapability('view_ratings');
  const canLiveGames = hasCapability('watch_live');
  const canFinishedGames = hasCapability('watch_results');

  const closeDrawer = () => setIsDrawerOpen(false);

  const goTo = (path: string) => {
    closeDrawer();
    navigate(path);
  };

  const handleFinishEvent = async () => {
    if (eventId == null) return;

    setIsFinishing(true);
    const ok = await finishEvent(eventId);
    setIsFinishing(false);

    if (!ok) return;

    showToast({
      title: 'Event finished',
      variant: 'warning',
    });
    refetch();
    closeDrawer();
  };

  if (eventId == null || Number.isNaN(eventId)) {
    return null;
  }

  return (
    <EventLayoutContext.Provider value={contextValue}>
      <LoadingTemplate
        isLoading={isLoading}
        isError={isError}
        isForbidden={isForbidden}
        onRefresh={refetch}
        onGoHome={() => navigate('/events')}
      >
        <div className="relative flex flex-col w-full max-w-full min-h-full min-w-0 overflow-x-hidden">
          <ContextMenu
            isOpen={isDrawerOpen}
            onOpen={() => setIsDrawerOpen(true)}
            onClose={closeDrawer}
            buttonClassName="fixed bottom-4 right-4 z-10"
          >
            <Drawer.Header>
              <Drawer.Heading className="break-words">
                {event?.name ?? 'Event'}
              </Drawer.Heading>
            </Drawer.Header>
            <Drawer.Body className="flex flex-col gap-2 flex-1">
              {canManage && (
                <Button
                  variant="tertiary"
                  fullWidth
                  // Tertiary text is muted inside overlays; keep it readable.
                  className="text-foreground"
                  onPress={() => goTo(`/event/${eventId}/manage`)}
                >
                  Manage
                </Button>
              )}
              {canRating && (
                <Button
                  variant="tertiary"
                  fullWidth
                  // Tertiary text is muted inside overlays; keep it readable.
                  className="text-foreground"
                  onPress={() => goTo(`/event/${eventId}/rating`)}
                >
                  Rating
                </Button>
              )}
              {canLiveGames && (
                <Button
                  variant="tertiary"
                  fullWidth
                  // Tertiary text is muted inside overlays; keep it readable.
                  className="text-foreground"
                  onPress={() => goTo(`/event/${eventId}/games/live`)}
                >
                  Live games
                </Button>
              )}
              {canFinishedGames && (
                <Button
                  variant="tertiary"
                  fullWidth
                  // Tertiary text is muted inside overlays; keep it readable.
                  className="text-foreground"
                  onPress={() => goTo(`/event/${eventId}/games/finished`)}
                >
                  Finished games
                </Button>
              )}

              {canFinishEvent && (
                <Button
                  variant="danger"
                  className="mt-auto w-1/2 self-center"
                  isPending={isFinishing}
                  onPress={onFinishConfirmOpen}
                >
                  {({ isPending }) => (
                    <>
                      {isPending && <Spinner color="current" size="sm" />}
                      Finish event
                    </>
                  )}
                </Button>
              )}
            </Drawer.Body>
          </ContextMenu>
          <Confirm
            isOpen={isFinishConfirmOpen}
            onClose={onFinishConfirmClose}
            onConfirm={() => {
              onFinishConfirmClose();
              handleFinishEvent();
            }}
            headerText="Confirm finish event"
            cancelText="Cancel"
            confirmText="Finish event"
            prompt="Are you sure you want to finish this event? This action cannot be undone."
            confirmVariant="danger"
          />
          <Outlet />
        </div>
      </LoadingTemplate>
    </EventLayoutContext.Provider>
  );
}
