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
import UserMenu from './UserMenu';

export default function Header() {
  const { isAuthenticated, signOut } = useContext(AuthContext);

  return (
    <Navbar maxWidth="full">
      <NavbarBrand className="gap-3">
        <Link color="foreground" href="/">
          <Image src="../makura.png" className="h-10" />
          <NavbarItem>MaCRa</NavbarItem>
        </Link>
      </NavbarBrand>
      <NavbarContent justify="end">
        {isAuthenticated ? (
          <UserMenu signOut={signOut} />
        ) : (
          <>
            <NavbarItem>
              <Button as={Link} color="primary" href="/signup" variant="light">
                Sign Up
              </Button>
            </NavbarItem>
            <NavbarItem>
              <Button as={Link} href="/login" color="primary" variant="solid">
                Login
              </Button>
            </NavbarItem>
          </>
        )}
      </NavbarContent>
    </Navbar>
  );
}
