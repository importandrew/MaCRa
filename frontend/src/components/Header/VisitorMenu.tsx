import { Avatar, Dropdown, Label } from '@heroui/react';
import { useTranslation } from 'react-i18next';
import { UserIcon } from '../Icons/UserIcon';

// VisitorMenu is the user menu of a visitor who isn't signed in: log in, and
// the settings (the language above all), which don't need an account.
export default function VisitorMenu() {
  const { t } = useTranslation();

  return (
    <Dropdown>
      <Dropdown.Trigger
        className="rounded-full"
        aria-label={t('header.userMenu')}
      >
        <Avatar>
          <Avatar.Fallback>
            <UserIcon className="w-6 h-6 text-muted" />
          </Avatar.Fallback>
        </Avatar>
      </Dropdown.Trigger>
      <Dropdown.Popover placement="bottom end">
        <Dropdown.Menu aria-label={t('header.userMenu')}>
          <Dropdown.Item id="login" textValue={t('header.logIn')} href="/login">
            <Label>{t('header.logIn')}</Label>
          </Dropdown.Item>
          <Dropdown.Item
            id="settings"
            textValue={t('header.settings')}
            href="/settings"
          >
            <Label>{t('header.settings')}</Label>
          </Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown.Popover>
    </Dropdown>
  );
}
