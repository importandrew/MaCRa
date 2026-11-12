import { useTranslation } from 'react-i18next';
import { Button, Drawer, useOverlayState } from '@heroui/react';
import { useCallback, useMemo, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { Outlet, useNavigate } from 'react-router-dom';
import { useRouteIds } from '../../hooks/useRouteIds';
import { useEventCapabilities } from '../../providers/event/useEventCapabilities';
import { finishEvent, useFetchEvent } from '../../providers/event/eventApi';
import {
  EventLayoutContext,
  EventLayoutContextValue,
} from '../../providers/event/eventLayoutContext';
import Confirm from '../Confirm';
import ContextMenu from '../ContextMenu/ContextMenu';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import { showToast } from '../Toast';
import PendingButton from '../PendingButton';

export default function EventLayout() {
  const { eventId } = useRouteIds();
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const { isPending: isFinishing, run } = useSubmitOnce();
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
    isNotFound: eventNotFound,
    refetch: refetchEvent,
  } = useFetchEvent(eventId);
  const {
    me,
    isLoading: meLoading,
    isError: meError,
    isForbidden: meForbidden,
    isNotFound: meNotFound,
    refetch: refetchMe,
    hasCapability,
  } = useEventCapabilities(eventId);

  const isLoading = eventLoading || meLoading;
  const isError = eventError || meError;
  const isForbidden = eventForbidden || meForbidden;
  // The server says the event doesn't exist (a bad id is handled below).
  const isNotFound = eventNotFound || meNotFound;

  const refetch = useCallback(() => {
    refetchEvent();
    refetchMe();
  }, [refetchEvent, refetchMe]);

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

  const menuItems = [
    { label: t('events.manage'), path: '/manage', isAllowed: canManage },
    { label: t('events.rating'), path: '/rating', isAllowed: canRating },
    {
      label: t('events.liveGames'),
      path: '/games/live',
      isAllowed: canLiveGames,
    },
    {
      label: t('events.finishedGames'),
      path: '/games/finished',
      isAllowed: canFinishedGames,
    },
  ];

  const goTo = (path: string) => {
    closeDrawer();
    navigate(path);
  };

  const handleFinishEvent = async () => {
    if (eventId == null) return;

    const ok = await finishEvent(eventId);

    if (!ok) return;

    showToast({
      title: t('events.finishedToast'),
      variant: 'success',
    });
    refetch();
    closeDrawer();
  };

  if (eventId == null) {
    return (
      <LoadingTemplate
        isLoading={false}
        isError={false}
        notFound="event"
        onRefresh={refetch}
        onGoHome={() => navigate('/events')}
      >
        {null}
      </LoadingTemplate>
    );
  }

  return (
    <EventLayoutContext.Provider value={contextValue}>
      <LoadingTemplate
        isLoading={isLoading}
        isError={isError}
        isForbidden={isForbidden}
        notFound={isNotFound ? 'event' : null}
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
                {event?.name ?? t('events.event')}
              </Drawer.Heading>
            </Drawer.Header>
            <Drawer.Body className="flex flex-col gap-2 flex-1">
              {menuItems
                .filter((item) => item.isAllowed)
                .map((item) => (
                  <Button
                    key={item.path}
                    variant="tertiary"
                    fullWidth
                    // Tertiary text is muted inside overlays; keep it readable.
                    className="text-foreground"
                    onPress={() => goTo(`/event/${eventId}${item.path}`)}
                  >
                    {item.label}
                  </Button>
                ))}

              {canFinishEvent && (
                <PendingButton
                  variant="danger"
                  className="mt-auto w-1/2 self-center"
                  isPending={isFinishing}
                  onPress={onFinishConfirmOpen}
                >
                  {t('events.finishEvent')}
                </PendingButton>
              )}
            </Drawer.Body>
          </ContextMenu>
          <Confirm
            isOpen={isFinishConfirmOpen}
            onClose={onFinishConfirmClose}
            onConfirm={() => {
              onFinishConfirmClose();
              void run(handleFinishEvent);
            }}
            headerText={t('events.finishConfirmTitle')}
            cancelText={t('common.cancel')}
            confirmText={t('events.finishEvent')}
            prompt={t('events.finishConfirmPrompt')}
            confirmVariant="danger"
          />
          <Outlet />
        </div>
      </LoadingTemplate>
    </EventLayoutContext.Provider>
  );
}
