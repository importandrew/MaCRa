import {
  Dropdown,
  Header,
  Label,
  Separator,
  useOverlayState,
} from '@heroui/react';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';
import { useEvent } from '../../providers/event/eventConstants';
import SelectEventModal from '../Events/SelectEventModal';
import UserAvatar from '../UserAvatar';

export default function UserMenu({ signOut }: { signOut: () => void }) {
  const username =
    JSON.parse(localStorage.getItem(AUTH_USER_DATA) || '{}')?.name || 'Profile';
  const { selectedEvent } = useEvent();
  const eventLabel = selectedEvent?.name || 'No event selected';
  const { isOpen, open: onOpen, close: onClose } = useOverlayState();

  return (
    <>
      <Dropdown>
        <Dropdown.Trigger className="rounded-full" aria-label="User menu">
          <UserAvatar name={username} />
        </Dropdown.Trigger>
        <Dropdown.Popover placement="bottom end">
          <Dropdown.Menu aria-label="User menu">
            <Dropdown.Section>
              <Header>{username}</Header>
              <Dropdown.Item id="profile" textValue="Profile" href="/profile">
                <Label>Profile</Label>
              </Dropdown.Item>
            </Dropdown.Section>
            <Separator />
            <Dropdown.Section>
              <Header className="block max-w-[24ch] truncate">
                {eventLabel}
              </Header>
              <Dropdown.Item
                id="select-event"
                textValue="Select event"
                onAction={onOpen}
              >
                <Label>{selectedEvent ? 'Switch event' : 'Select event'}</Label>
              </Dropdown.Item>
            </Dropdown.Section>
            <Separator />
            <Dropdown.Section>
              <Dropdown.Item
                id="logout"
                textValue="Log Out"
                variant="danger"
                onAction={signOut}
              >
                <Label>Log Out</Label>
              </Dropdown.Item>
            </Dropdown.Section>
          </Dropdown.Menu>
        </Dropdown.Popover>
      </Dropdown>
      <SelectEventModal isOpen={isOpen} onClose={onClose} />
    </>
  );
}
