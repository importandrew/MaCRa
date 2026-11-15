import { useTranslation } from 'react-i18next';
import { Checkbox, CheckboxGroup, Label } from '@heroui/react';
import type { EventRole } from '../../models/EventRole';
import { ROLE_ICONS, rolesForEvent, withoutConflicts } from './roles';

// RoleCheckboxes picks a member's roles, each with its icon: when adding a
// member and in the role switcher. The roles are those of the event's type;
// checking player or substitute unchecks the other.
export default function RoleCheckboxes({
  value,
  onChange,
  eventType,
}: {
  value: EventRole[];
  onChange: (roles: EventRole[]) => void;
  eventType: 'regular' | 'tournament';
}) {
  const { t } = useTranslation();
  return (
    <CheckboxGroup
      value={value}
      onChange={(next) =>
        onChange(withoutConflicts(value, next as EventRole[]))
      }
    >
      <Label>{t('members.roles')}</Label>
      {rolesForEvent(eventType).map((role) => (
        <Checkbox key={role} value={role}>
          <Checkbox.Content>
            <Checkbox.Control>
              <Checkbox.Indicator />
            </Checkbox.Control>
            <span className="flex items-center gap-2 leading-none">
              {ROLE_ICONS[role]({ className: 'w-6 h-6 shrink-0' })}
              {t(`roles.${role}`)}
            </span>
          </Checkbox.Content>
        </Checkbox>
      ))}
    </CheckboxGroup>
  );
}
