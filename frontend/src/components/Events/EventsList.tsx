import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';
import { cardVariants } from '@heroui/styles';
import { type ReactNode, useContext, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import type { Event } from '../../models/Event';
import type { EventRole } from '../../models/EventRole';
import { AuthContext } from '../../providers/auth/authConstants';
import { useFetchVisibleEvents } from '../../providers/event/eventApi';
import { SignIcon } from '../Icons/SignIcon';
import { CrownIcon } from '../Icons/CrownIcon';
import { EyeIcon } from '../Icons/EyeIcon';
import { CheckInCircleIcon } from '../Icons/CheckInCircleIcon';
import { MahjongTileIcon } from '../Icons/MahjongTileIcon';
import { ScalesIcon } from '../Icons/ScalesIcon';
import { TrophyIcon } from '../Icons/TrophyIcon';
import { UsersIcon } from '../Icons/UsersIcon';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import SegmentedPicker from '../SegmentedPicker';
import EventSearchField from './EventSearchField';
import FinishedEvents from './FinishedEvents';
import {
  type EventsTab,
  matchesQuery,
  readEventsTab,
  splitFinished,
  storeEventsTab,
} from './eventListFilters';

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

// EventCard is an event in the list: a link to it, with its name, type, the
// user's roles and the start of its description. Finished events are dimmed.
function EventCard({ event }: { event: Event }) {
  const isFinished = Boolean(event.finished_at);
  return (
    // A real link: long-press and "open in new tab" work on phones.
    <Link
      to={`/event/${event.id}`}
      className={cardStyles.base({
        className: isFinished
          ? 'opacity-60 bg-default/60 border border-border shadow-none'
          : undefined,
      })}
    >
      <div
        className={cardStyles.header({
          className: 'flex-row justify-between items-center gap-2 min-w-0',
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
          <p className="text-sm truncate text-muted">{event.description}</p>
        </div>
      )}
    </Link>
  );
}

// EventList is one tab of /events: its active events, then its finished
// ones collapsed (open while searching).
function EventList({
  events,
  isSearching,
}: {
  events: Event[];
  isSearching: boolean;
}) {
  const { active, finished } = splitFinished(events);
  return (
    <div className="flex flex-col gap-3">
      {active.map((event) => (
        <EventCard key={event.id} event={event} />
      ))}
      <FinishedEvents count={finished.length} isSearching={isSearching}>
        {finished.map((event) => (
          <EventCard key={event.id} event={event} />
        ))}
      </FinishedEvents>
    </div>
  );
}

// EventsList is /events: two tabs, "My events" (where the user is a member)
// and "Public events", each with its finished events collapsed, and a search
// by name over both (the tabs count the matches). A visitor sees the public
// events only, without tabs.
export default function EventsList() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { isAuthenticated } = useContext(AuthContext);
  const { data, isLoading, isError, refetch } = useFetchVisibleEvents();
  const [query, setQuery] = useState('');
  // chosenTab: the tab the user picked here or last time on this device.
  const [chosenTab, setChosenTab] = useState<EventsTab | null>(readEventsTab);

  // handleCreate: a visitor signs in first and comes back to the form.
  const handleCreate = () => {
    if (!isAuthenticated) {
      navigate('/login', { state: { from: { pathname: '/events/new' } } });
      return;
    }
    navigate('/events/new');
  };

  const events = data ?? [];
  const isSearching = query.trim() !== '';
  const found = events.filter((event) => matchesQuery(event, query));
  const mine = found.filter((event) => event.is_participating);
  const others = found.filter((event) => !event.is_participating);
  // Without a choice: "My events" if the user has any, else "Public".
  const hasOwn = events.some((event) => event.is_participating);
  const tab: EventsTab = !isAuthenticated
    ? 'public'
    : (chosenTab ?? (hasOwn ? 'mine' : 'public'));
  const pickTab = (next: EventsTab) => {
    setChosenTab(next);
    storeEventsTab(next);
  };
  const shown = tab === 'mine' ? mine : others;

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      onRefresh={refetch}
    >
      <div className="relative flex flex-col w-full max-w-full px-3 pt-6 pb-24 gap-4 min-w-0 overflow-x-hidden items-stretch">
        <div className="flex flex-col w-full max-w-[600px] mx-auto gap-4 min-w-0">
          {/* The one action of this page: a round (+) in the corner, no ≡
              menu (docs/user-scenarios.md, "Menus and overlays"). */}
          <Button
            isIconOnly
            variant="primary"
            aria-label={t('events.newEvent')}
            onPress={handleCreate}
            className="fixed bottom-4 right-4 z-10 size-14 rounded-full shadow-lg"
          >
            <SignIcon sign="plus" className="size-7 [&_path]:stroke-2" />
          </Button>

          <h1 className="text-2xl font-semibold">{t('events.title')}</h1>

          {events.length === 0 ? (
            <p className="text-center text-muted text-xl py-8">
              {t('events.noEvents')}
            </p>
          ) : (
            <>
              <EventSearchField value={query} onChange={setQuery} />
              {isAuthenticated && (
                <SegmentedPicker
                  ariaLabel={t('events.title')}
                  options={[
                    {
                      value: 'mine' as const,
                      label: t('events.myEventsCount', { n: mine.length }),
                    },
                    {
                      value: 'public' as const,
                      label: t('events.publicEventsCount', {
                        n: others.length,
                      }),
                    },
                  ]}
                  selectedKey={tab}
                  onSelectionChange={pickTab}
                />
              )}
              {shown.length === 0 ? (
                <p className="text-center text-muted py-6">
                  {isSearching
                    ? t('events.nothingFound', { query: query.trim() })
                    : tab === 'mine'
                      ? t('events.notParticipating')
                      : t('events.noPublicEvents')}
                </p>
              ) : (
                <EventList events={shown} isSearching={isSearching} />
              )}
            </>
          )}
        </div>
      </div>
    </LoadingTemplate>
  );
}
