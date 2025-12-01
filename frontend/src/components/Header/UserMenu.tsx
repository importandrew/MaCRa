import {
  Avatar,
  Dropdown,
  DropdownItem,
  DropdownMenu,
  DropdownTrigger,
} from '@heroui/react';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';

export default function UserMenu({ signOut }: { signOut: () => void }) {
  const username =
    JSON.parse(localStorage.getItem(AUTH_USER_DATA) || '{}')?.name || 'Profile';
  const initials = username
    .trim()
    .split(' ')
    .map((item: string) => item[0].toLocaleUpperCase())
    .join('');

  return (
    <Dropdown placement="bottom-end">
      <DropdownTrigger>
        <Avatar
          as="button"
          className="transition-transform"
          color="default"
          size="md"
          name={initials}
        />
      </DropdownTrigger>
      <DropdownMenu
        disabledKeys={['profile']}
        aria-label="User menu"
        variant="flat"
      >
        <DropdownItem key="profile" className="h-14 gap-2">
          <p className="font-semibold">{username}</p>
        </DropdownItem>
        <DropdownItem key="logout" color="danger" onPress={signOut}>
          Log Out
        </DropdownItem>
      </DropdownMenu>
    </Dropdown>
  );
}
