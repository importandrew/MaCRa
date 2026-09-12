import {
  Button,
  Image,
  Link,
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
} from '@heroui/react';
import { useContext } from 'react';
import { AuthContext } from '../../providers/auth/authConstants';
import { CalendarIcon } from '../Icons/CalendarIcon';
import UserMenu from './UserMenu';

export default function Header() {
  const { isAuthenticated, signOut } = useContext(AuthContext);

  return (
    // Brand defaults to flex-grow/basis-0; override so logo keeps intrinsic size.
    // Start + end contents both grow and share the remaining width.
    <Navbar maxWidth="full" classNames={{ wrapper: 'min-w-0 gap-2 px-4 sm:gap-4 sm:px-6' }}>
      <NavbarBrand className="grow-0 basis-auto shrink-0 max-w-fit">
        <Link color="foreground" href="/" className="flex items-center gap-2">
          <Image
            src="/makura.png"
            alt="MaCRa"
            width={40}
            height={40}
            className="h-10 w-10"
            classNames={{ wrapper: 'h-10 w-10 shrink-0' }}
          />
          <NavbarItem>MaCRa</NavbarItem>
        </Link>
      </NavbarBrand>
      <NavbarContent justify="start" className="gap-2 min-w-0">
        <NavbarItem>
          <Button
            as={Link}
            href="/events"
            variant="light"
            startContent={<CalendarIcon className="w-6 h-6 shrink-0" />}
          >
            Events
          </Button>
        </NavbarItem>
      </NavbarContent>
      <NavbarContent justify="end" className="gap-2 min-w-0">
        {isAuthenticated ? (
          <UserMenu signOut={signOut} />
        ) : (
          <NavbarItem>
            <Button as={Link} href="/login" color="primary" variant="solid">
              Login
            </Button>
          </NavbarItem>
        )}
      </NavbarContent>
    </Navbar>
  );
}
