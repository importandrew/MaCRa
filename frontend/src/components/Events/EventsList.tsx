import {
  Button,
  Card,
  CardBody,
  CardHeader,
  DrawerBody,
  DrawerHeader,
} from '@heroui/react';
import { type ReactNode, useContext, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import type { Event } from '../../models/Event';
import type { EventRole } from '../../models/EventRole';
import { AuthContext } from '../../providers/auth/authConstants';
import { useFetchVisibleEvents } from '../../providers/event/eventApi';
import ContextMenuButton from '../ContextMenu/ContextMenuButton';
import SideMenu from '../ContextMenu/SideMenu';
import { CrownIcon } from '../Icons/CrownIcon';
import { EyeIcon } from '../Icons/EyeIcon';
import { CheckInCircleIcon } from '../Icons/CheckInCircleIcon';
import { MahjongTileIcon } from '../Icons/MahjongTileIcon';
import { ScalesIcon } from '../Icons/ScalesIcon';
import { TrophyIcon } from '../Icons/TrophyIcon';
import { UsersIcon } from '../Icons/UsersIcon';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';

const ROLE_ORDER: EventRole[] = ['player', 'referee', 'observer'];

const ROLE_META: Record<
  EventRole,
  { label: string; icon: (props: { className?: string }) => ReactNode }
> = {
  player: {
    label: 'Player',
    icon: (props) => <MahjongTileIcon {...props} />,
  },
  referee: {
    label: 'Referee',
    icon: (props) => <ScalesIcon {...props} />,
  },
  observer: {
    label: 'Observer',
    icon: (props) => <EyeIcon {...props} />,
  },
};

function MetaItem({
  icon,
  label,
}: {
  icon: ReactNode;
  label: string;
}) {
  return (
    <span
      className="inline-flex shrink-0"
      title={label}
      aria-label={label}
      role="img"
    >
      {icon}
    </span>
  );
}

function EventCardMeta({ event }: { event: Event }) {
  const isFinished = Boolean(event.finished_at);
  const roles = event.roles ?? [];
  const orderedRoles = ROLE_ORDER.filter((role) => roles.includes(role));
  const iconClass = 'w-6 h-6';

  return (
    <div
      className={`flex items-center gap-2 shrink-0 ${
        isFinished ? 'text-default-400' : 'text-default-500'
      }`}
    >
      {isFinished && (
        <MetaItem
          icon={<CheckInCircleIcon className={iconClass} />}
          label="Finished"
        />
      )}
      <MetaItem
        icon={
          event.type === 'tournament' ? (
            <TrophyIcon className={iconClass} />
          ) : (
            <UsersIcon className={iconClass} />
          )
        }
        label={event.type === 'tournament' ? 'Tournament' : 'Regular'}
      />
      {event.is_owner && (
        <MetaItem
          icon={<CrownIcon className={iconClass} />}
          label="Owner"
        />
      )}
      {orderedRoles.map((role) => {
        const meta = ROLE_META[role];
        return (
          <MetaItem
            key={role}
            icon={meta.icon({ className: iconClass })}
            label={meta.label}
          />
        );
      })}
    </div>
  );
}

export default function EventsList() {
  const navigate = useNavigate();
  const { isAuthenticated } = useContext(AuthContext);
  const { data, isLoading, isError, refetch } = useFetchVisibleEvents();
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  const closeDrawer = () => setIsDrawerOpen(false);

  const handleCreate = () => {
    closeDrawer();
    if (!isAuthenticated) {
      navigate('/login', { state: { from: { pathname: '/events/new' } } });
      return;
    }
    navigate('/events/new');
  };

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      onRefresh={refetch}
    >
      <div className="relative flex flex-col w-full max-w-full px-3 py-6 gap-4 min-w-0 overflow-x-hidden items-stretch">
        <div className="flex flex-col w-full max-w-[600px] mx-auto gap-4 min-w-0">
          <ContextMenuButton
            onPress={() => setIsDrawerOpen(true)}
            className="fixed bottom-4 right-4 z-100"
          />

          <SideMenu isOpen={isDrawerOpen} onOpenChange={closeDrawer}>
            <DrawerHeader>Events</DrawerHeader>
            <DrawerBody className="flex flex-col gap-3">
              <Button color="primary" onPress={handleCreate}>
                New event
              </Button>
            </DrawerBody>
          </SideMenu>

          <h1 className="text-2xl font-semibold">Events</h1>

          {!data?.length ? (
            <p className="text-center text-default-500 text-xl py-8">
              No events available.
            </p>
          ) : (
            data.map((event) => {
              const isFinished = Boolean(event.finished_at);
              return (
                <Card
                  key={event.id}
                  isPressable
                  onPress={() => navigate(`/event/${event.id}`)}
                  className={
                    isFinished
                      ? 'opacity-60 bg-default-100/60 border border-default-200 shadow-none'
                      : undefined
                  }
                >
                  <CardHeader className="flex justify-between items-center gap-2 min-w-0">
                    <span
                      className={`font-medium min-w-0 truncate text-lg ${
                        isFinished ? 'text-default-400' : ''
                      }`}
                      title={event.name}
                    >
                      {event.name}
                    </span>
                    <EventCardMeta event={event} />
                  </CardHeader>
                  {event.description && (
                    <CardBody className="pt-0">
                      <p
                        className={`text-sm truncate ${
                          isFinished
                            ? 'text-default-400'
                            : 'text-default-600'
                        }`}
                      >
                        {event.description}
                      </p>
                    </CardBody>
                  )}
                </Card>
              );
            })
          )}
        </div>
      </div>
    </LoadingTemplate>
  );
}
