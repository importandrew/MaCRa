import {
  Dropdown,
  Header,
  Label,
  Separator,
  useOverlayState,
} from '@heroui/react';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../../providers/auth/authConstants';
import { useEvent } from '../../providers/event/eventConstants';
import SelectEventModal from '../Events/SelectEventModal';
import UserAvatar from '../UserAvatar';

export default function UserMenu({ signOut }: { signOut: () => void }) {
  const { t } = useTranslation();
  const { user } = useAuth();
  const username = user?.name || t('header.profile');
  const { selectedEvent } = useEvent();
  const eventLabel = selectedEvent?.name || t('header.noEventSelected');
  const { isOpen, open: onOpen, close: onClose } = useOverlayState();

  return (
    <>
      <Dropdown>
        <Dropdown.Trigger
          className="rounded-full"
          aria-label={t('header.userMenu')}
        >
          <UserAvatar name={username} />
        </Dropdown.Trigger>
        <Dropdown.Popover placement="bottom end">
          <Dropdown.Menu aria-label={t('header.userMenu')}>
            <Dropdown.Section>
              <Header>{username}</Header>
              <Dropdown.Item
                id="profile"
                textValue={t('header.profile')}
                href="/profile"
              >
                <Label>{t('header.profile')}</Label>
              </Dropdown.Item>
              <Dropdown.Item
                id="settings"
                textValue={t('header.settings')}
                href="/settings"
              >
                <Label>{t('header.settings')}</Label>
              </Dropdown.Item>
            </Dropdown.Section>
            <Separator />
            <Dropdown.Section>
              <Header className="block max-w-[24ch] truncate">
                {eventLabel}
              </Header>
              <Dropdown.Item
                id="select-event"
                textValue={t('header.selectEvent')}
                onAction={onOpen}
              >
                <Label>
                  {selectedEvent
                    ? t('header.switchEvent')
                    : t('header.selectEvent')}
                </Label>
              </Dropdown.Item>
            </Dropdown.Section>
            <Separator />
            <Dropdown.Section>
              <Dropdown.Item
                id="logout"
                textValue={t('header.logOut')}
                variant="danger"
                onAction={signOut}
              >
                <Label>{t('header.logOut')}</Label>
              </Dropdown.Item>
            </Dropdown.Section>
          </Dropdown.Menu>
        </Dropdown.Popover>
      </Dropdown>
      <SelectEventModal isOpen={isOpen} onClose={onClose} />
    </>
  );
}
