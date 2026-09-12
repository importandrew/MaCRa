import { Chip, Divider } from '@heroui/react';
import { type ReactNode, useMemo } from 'react';
import type { EventMember } from '../../models/EventRole';
import { useFetchEventMembers } from '../../providers/event/eventApi';
import { useEventLayout } from '../../providers/event/eventLayoutContext';
import { getDateTimeFormatted } from '../../utils/date-tools';
import { CalendarIcon } from '../Icons/CalendarIcon';
import { CheckInCircleIcon } from '../Icons/CheckInCircleIcon';
import { CrownIcon } from '../Icons/CrownIcon';
import { EyeIcon } from '../Icons/EyeIcon';
import { ScalesIcon } from '../Icons/ScalesIcon';
import { TrophyIcon } from '../Icons/TrophyIcon';
import { UsersIcon } from '../Icons/UsersIcon';
import { MahjongTileIcon } from '../Icons/MahjongTileIcon';

const iconClass = 'w-6 h-6';

function namesWithRole(
  members: EventMember[] | null,
  role: 'referee' | 'observer' | 'player'
): string[] {
  if (!members) return [];
  return members
    .filter((m) => m.roles.includes(role))
    .map((m) => m.name)
    .sort((a, b) => a.localeCompare(b));
}

function DetailRow({
  icon,
  label,
  children,
}: {
  icon?: ReactNode;
  label: string;
  children: ReactNode;
}) {
  return (
    <div>
      <dt className="flex items-center gap-1.5 font-medium text-default-500">
        {icon != null && (
          <span className="shrink-0" aria-hidden>
            {icon}
          </span>
        )}
        {label}
      </dt>
      <dd className="mt-0.5">{children}</dd>
    </div>
  );
}

export default function EventOverview() {
  const { eventId, event } = useEventLayout();
  const { data: members } = useFetchEventMembers(eventId);

  const ownerName = useMemo(() => {
    if (!event) return null;
    const owner = members?.find((m) => m.user_id === event.owner_user_id);
    return owner?.name ?? null;
  }, [event, members]);

  const referees = useMemo(
    () => namesWithRole(members, 'referee'),
    [members]
  );
  const observers = useMemo(
    () => namesWithRole(members, 'observer'),
    [members]
  );
  const players = useMemo(
    () => namesWithRole(members, 'player'),
    [members]
  );

  if (!event) return null;

  const isFinished = Boolean(event.finished_at);
  const typeLabel = event.type === 'tournament' ? 'Tournament' : 'Regular';

  return (
    <div className="flex flex-col w-full max-w-[600px] mx-auto px-3 py-6 gap-4 min-w-0">
      <div className="flex flex-col items-center gap-2 text-center">
        <h1
          className={`text-2xl font-semibold break-words w-full ${
            isFinished ? 'text-default-500' : ''
          }`}
        >
          {event.name}
        </h1>
        <div className="flex gap-1 flex-wrap justify-center">
          <Chip
            size="sm"
            variant="flat"
            startContent={
              event.type === 'tournament' ? (
                <TrophyIcon className={iconClass} />
              ) : (
                <UsersIcon className={iconClass} />
              )
            }
          >
            {typeLabel}
          </Chip>
          {isFinished && (
            <Chip
              size="sm"
              color="default"
              variant="flat"
              startContent={<CheckInCircleIcon className={iconClass} />}
            >
              Finished
            </Chip>
          )}
        </div>
      </div>

      {event.description && (
        <p className="text-default-600 break-words">{event.description}</p>
      )}

      <Divider />

      <dl className="flex flex-col gap-3 text-default-600">
        <DetailRow icon={<CrownIcon className={iconClass} />} label="Owner">
          {ownerName ?? `User #${event.owner_user_id}`}
        </DetailRow>
        <DetailRow
          icon={<MahjongTileIcon className={iconClass} />}
          label="Players"
        >
          {players.length}
        </DetailRow>  
        <DetailRow
          icon={<ScalesIcon className={iconClass} />}
          label="Referees"
        >
          {referees.length > 0 ? referees.join(', ') : '—'}
        </DetailRow>
        <DetailRow
          icon={<EyeIcon className={iconClass} />}
          label="Observers"
        >
          {observers.length > 0 ? observers.join(', ') : '—'}
        </DetailRow>
        <DetailRow
          icon={<CalendarIcon className={iconClass} />}
          label="Created"
        >
          {getDateTimeFormatted(event.created_at)}
        </DetailRow>
        {event.finished_at && (
          <DetailRow
            icon={<CheckInCircleIcon className={iconClass} />}
            label="Finished"
          >
            {getDateTimeFormatted(event.finished_at)}
          </DetailRow>
        )}
      </dl>
    </div>
  );
}
