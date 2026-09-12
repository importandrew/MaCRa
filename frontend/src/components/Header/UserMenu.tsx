import {
  Avatar,
  Dropdown,
  DropdownItem,
  DropdownMenu,
  DropdownSection,
  DropdownTrigger,
  useDisclosure,
} from '@heroui/react';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';
import { useEvent } from '../../providers/event/eventConstants';
import SelectEventModal from '../Events/SelectEventModal';

export default function UserMenu({ signOut }: { signOut: () => void }) {
  const username =
    JSON.parse(localStorage.getItem(AUTH_USER_DATA) || '{}')?.name || 'Profile';
  const { selectedEvent } = useEvent();
  const eventLabel = selectedEvent?.name || 'No event selected';
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <Dropdown placement="bottom-end">
        <DropdownTrigger>
          <Avatar
            as="button"
            className="transition-transform"
            style={{ cursor: 'pointer' }}
            color="default"
            size="md"
            name={username}
          />
        </DropdownTrigger>
        <DropdownMenu aria-label="User menu" variant="flat">
          <DropdownSection title={username} showDivider>
            <DropdownItem key="profile" href="/profile">
              Profile
            </DropdownItem>
          </DropdownSection>
          <DropdownSection
            title={eventLabel}
            classNames={{ heading: 'block max-w-[24ch] truncate' }}
            showDivider
          >
            <DropdownItem key="select-event" onPress={onOpen}>
              {selectedEvent ? 'Switch event' : 'Select event'}
            </DropdownItem>
          </DropdownSection>
          <DropdownSection>
            <DropdownItem
              key="logout"
              className="text-danger"
              onPress={signOut}
            >
              Log Out
            </DropdownItem>
          </DropdownSection>
        </DropdownMenu>
      </Dropdown>
      <SelectEventModal isOpen={isOpen} onClose={onClose} />
    </>
  );
}
