import { useTranslation } from 'react-i18next';
import { Button, Drawer } from '@heroui/react';
import { cardVariants } from '@heroui/styles';
import { type ReactNode, useContext, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import type { Event } from '../../models/Event';
import type { EventRole } from '../../models/EventRole';
import { AuthContext } from '../../providers/auth/authConstants';
import { useFetchVisibleEvents } from '../../providers/event/eventApi';
import ContextMenu from '../ContextMenu/ContextMenu';
import { CrownIcon } from '../Icons/CrownIcon';
import { EyeIcon } from '../Icons/EyeIcon';
import { CheckInCircleIcon } from '../Icons/CheckInCircleIcon';
import { MahjongTileIcon } from '../Icons/MahjongTileIcon';
import { ScalesIcon } from '../Icons/ScalesIcon';
import { TrophyIcon } from '../Icons/TrophyIcon';
import { UsersIcon } from '../Icons/UsersIcon';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';

const cardStyles = cardVariants();

const ROLE_ORDER: EventRole[] = ['player', 'referee', 'observer'];

const ROLE_ICONS: Record<
  EventRole,
  (props: { className?: string }) => ReactNode
> = {
  player: (props) => <MahjongTileIcon {...props} />,
  referee: (props) => <ScalesIcon {...props} />,
  observer: (props) => <EyeIcon {...props} />,
};

function MetaItem({ icon, label }: { icon: ReactNode; label: string }) {
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
  const { t } = useTranslation();
  const isFinished = Boolean(event.finished_at);
  const roles = event.roles ?? [];
  const orderedRoles = ROLE_ORDER.filter((role) => roles.includes(role));
  const iconClass = 'w-6 h-6';

  return (
    <div
      className={`flex items-center gap-2 shrink-0 ${
        isFinished ? 'text-muted/60' : 'text-muted'
      }`}
    >
      {isFinished && (
        <MetaItem
          icon={<CheckInCircleIcon className={iconClass} />}
          label={t('events.finished')}
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
        label={t(`eventType.${event.type}`)}
      />
      {event.is_owner && (
        <MetaItem
          icon={<CrownIcon className={iconClass} />}
          label={t('roles.owner')}
        />
      )}
      {orderedRoles.map((role) => (
        <MetaItem
          key={role}
          icon={ROLE_ICONS[role]({ className: iconClass })}
          label={t(`roles.${role}`)}
        />
      ))}
    </div>
  );
}

export default function EventsList() {
  const { t } = useTranslation();
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
          <ContextMenu
            isOpen={isDrawerOpen}
            onOpen={() => setIsDrawerOpen(true)}
            onClose={closeDrawer}
            buttonClassName="fixed bottom-4 right-4 z-10"
          >
            <Drawer.Header>
              <Drawer.Heading>{t('events.title')}</Drawer.Heading>
            </Drawer.Header>
            <Drawer.Body className="flex flex-col gap-3">
              <Button variant="primary" fullWidth onPress={handleCreate}>
                {t('events.newEvent')}
              </Button>
            </Drawer.Body>
          </ContextMenu>

          <h1 className="text-2xl font-semibold">{t('events.title')}</h1>

          {!data?.length ? (
            <p className="text-center text-muted text-xl py-8">
              {t('events.noEvents')}
            </p>
          ) : (
            data.map((event) => {
              const isFinished = Boolean(event.finished_at);
              return (
                // A real link: long-press and "open in new tab" work on phones.
                <Link
                  key={event.id}
                  to={`/event/${event.id}`}
                  className={cardStyles.base({
                    className: isFinished
                      ? 'opacity-60 bg-default/60 border border-border shadow-none'
                      : undefined,
                  })}
                >
                  <div
                    className={cardStyles.header({
                      className:
                        'flex-row justify-between items-center gap-2 min-w-0',
                    })}
                  >
                    <span
                      className={`font-medium min-w-0 truncate text-lg ${
                        isFinished ? 'text-muted' : ''
                      }`}
                      title={event.name}
                    >
                      {event.name}
                    </span>
                    <EventCardMeta event={event} />
                  </div>
                  {event.description && (
                    <div className={cardStyles.content()}>
                      <p className="text-sm truncate text-muted">
                        {event.description}
                      </p>
                    </div>
                  )}
                </Link>
              );
            })
          )}
        </div>
      </div>
    </LoadingTemplate>
  );
}
