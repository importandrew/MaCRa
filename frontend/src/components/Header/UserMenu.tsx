import {
  Avatar,
  Dropdown,
  DropdownItem,
  DropdownMenu,
  DropdownSection,
  DropdownTrigger,
} from '@heroui/react';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';

export default function UserMenu({ signOut }: { signOut: () => void }) {
  const username =
    JSON.parse(localStorage.getItem(AUTH_USER_DATA) || '{}')?.name || 'Profile';

  return (
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
        <DropdownSection>
          <DropdownItem key="logout" className="text-danger" onPress={signOut}>
            Log Out
          </DropdownItem>
        </DropdownSection>
      </DropdownMenu>
    </Dropdown>
  );
}
