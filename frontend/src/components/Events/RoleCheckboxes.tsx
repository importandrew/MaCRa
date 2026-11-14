import { useTranslation } from 'react-i18next';
import { Checkbox, CheckboxGroup, Label } from '@heroui/react';
import type { EventRole } from '../../models/EventRole';
import { ALL_ROLES, ROLE_ICONS } from './roles';

// RoleCheckboxes picks a member's roles, each with its icon: when adding a
// member and in the role switcher.
export default function RoleCheckboxes({
  value,
  onChange,
}: {
  value: EventRole[];
  onChange: (roles: EventRole[]) => void;
}) {
  const { t } = useTranslation();
  return (
    <CheckboxGroup value={value} onChange={(v) => onChange(v as EventRole[])}>
      <Label>{t('members.roles')}</Label>
      {ALL_ROLES.map((role) => (
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
