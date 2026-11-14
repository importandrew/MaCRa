import { useTranslation } from 'react-i18next';
import { Chip, Separator } from '@heroui/react';
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
import { BookIcon } from '../Icons/BookIcon';
import { UserIcon } from '../Icons/UserIcon';
import { settingsDifferences, rulesDifferences } from './rulesText';

const iconClass = 'w-6 h-6';

// namesWithRole lists the staff by name. Non-members get only the staff from
// /members, so it is never used for players (see player_count).
function namesWithRole(
  members: EventMember[] | null,
  role: 'referee' | 'observer'
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
      <dt className="flex items-center gap-1.5 font-medium text-muted">
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
  const { t } = useTranslation();
  const { eventId, event } = useEventLayout();
  const { data: members } = useFetchEventMembers(eventId);

  const ownerName = useMemo(() => {
    if (!event) return null;
    const owner = members?.find((m) => m.user_id === event.owner_user_id);
    return owner?.name ?? null;
  }, [event, members]);

  const referees = useMemo(() => namesWithRole(members, 'referee'), [members]);
  const observers = useMemo(
    () => namesWithRole(members, 'observer'),
    [members]
  );

  if (!event) return null;

  const isFinished = Boolean(event.finished_at);
  const differences = rulesDifferences(event.rules, t);
  const settingsLine = settingsDifferences(event.settings, t);
  const typeLabel = t(`eventType.${event.type}`);

  return (
    <div className="flex flex-col w-full max-w-[600px] mx-auto px-3 py-6 gap-4 min-w-0">
      <div className="flex flex-col items-center gap-2 text-center">
        <h1
          className={`text-2xl font-semibold break-words w-full ${
            isFinished ? 'text-muted' : ''
          }`}
        >
          {event.name}
        </h1>
        <div className="flex gap-1 flex-wrap justify-center">
          <Chip size="sm" variant="soft">
            {event.type === 'tournament' ? (
              <TrophyIcon className={iconClass} />
            ) : (
              <UsersIcon className={iconClass} />
            )}
            {typeLabel}
          </Chip>
          {isFinished && (
            <Chip size="sm" variant="soft">
              <CheckInCircleIcon className={iconClass} />
              {t('events.finished')}
            </Chip>
          )}
        </div>
      </div>

      {event.description && (
        <p className="text-muted break-words">{event.description}</p>
      )}

      <Separator />

      <dl className="flex flex-col gap-3 text-muted">
        <DetailRow
          icon={<CrownIcon className={iconClass} />}
          label={t('roles.owner')}
        >
          {ownerName ?? t('events.userNumber', { id: event.owner_user_id })}
        </DetailRow>
        <DetailRow
          icon={<MahjongTileIcon className={iconClass} />}
          label={t('events.players')}
        >
          {event.player_count ?? '—'}
        </DetailRow>
        <DetailRow
          icon={<BookIcon className={iconClass} />}
          label={t('events.rules')}
        >
          {differences.length > 0
            ? differences.join(', ')
            : t('rules.presets.standard')}
        </DetailRow>
        {settingsLine.length > 0 && (
          <DetailRow
            icon={<UserIcon className={iconClass} />}
            label={t('eventSettings.title')}
          >
            {settingsLine.join(', ')}
          </DetailRow>
        )}
        <DetailRow
          icon={<ScalesIcon className={iconClass} />}
          label={t('events.referees')}
        >
          {referees.length > 0 ? referees.join(', ') : '—'}
        </DetailRow>
        <DetailRow
          icon={<EyeIcon className={iconClass} />}
          label={t('events.observers')}
        >
          {observers.length > 0 ? observers.join(', ') : '—'}
        </DetailRow>
        <DetailRow
          icon={<CalendarIcon className={iconClass} />}
          label={t('events.createdAt')}
        >
          {getDateTimeFormatted(event.created_at)}
        </DetailRow>
        {event.finished_at && (
          <DetailRow
            icon={<CheckInCircleIcon className={iconClass} />}
            label={t('events.finished')}
          >
            {getDateTimeFormatted(event.finished_at)}
          </DetailRow>
        )}
      </dl>
    </div>
  );
}
