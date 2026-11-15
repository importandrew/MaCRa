import type { ReactNode } from 'react';
import type { EventRole } from '../../models/EventRole';
import { EyeIcon } from '../Icons/EyeIcon';
import { MahjongTileIcon } from '../Icons/MahjongTileIcon';
import { ScalesIcon } from '../Icons/ScalesIcon';
import { SwapIcon } from '../Icons/SwapIcon';

// ALL_ROLES: the member roles, in the order screens show them. The
// substitute role exists only in tournaments (rolesForEvent).
export const ALL_ROLES: EventRole[] = [
  'player',
  'substitute',
  'referee',
  'observer',
];

// rolesForEvent: the roles a member of an event of this type may get.
export function rolesForEvent(type: 'regular' | 'tournament'): EventRole[] {
  return type === 'tournament'
    ? ALL_ROLES
    : ALL_ROLES.filter((role) => role !== 'substitute');
}

// EXCLUSIVE_ROLES: roles that can't go together; checking one unchecks the
// other.
const EXCLUSIVE_ROLES: [EventRole, EventRole][] = [['player', 'substitute']];

// withoutConflicts keeps the roles just checked and drops the roles they
// exclude.
export function withoutConflicts(
  previous: EventRole[],
  next: EventRole[]
): EventRole[] {
  const added = next.filter((role) => !previous.includes(role));
  return next.filter(
    (role) =>
      !EXCLUSIVE_ROLES.some(
        ([a, b]) =>
          (role === a && added.includes(b)) || (role === b && added.includes(a))
      )
  );
}

export const ROLE_ICONS: Record<
  EventRole,
  (props: { className?: string }) => ReactNode
> = {
  player: (props) => <MahjongTileIcon {...props} />,
  referee: (props) => <ScalesIcon {...props} />,
  observer: (props) => <EyeIcon {...props} />,
  substitute: (props) => <SwapIcon {...props} />,
};
