import { useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { AuthContext } from '../../providers/auth/authConstants';
import { CalendarIcon } from '../Icons/CalendarIcon';
import UserMenu from './UserMenu';
import VisitorMenu from './VisitorMenu';
import { LinkButton } from '../LinkButton';

// Header is the sticky top bar (HeroUI v3 has no Navbar): logo, Events and
// the user menu (a visitor's one with Log in and Settings when not signed in).
export default function Header() {
  const { t } = useTranslation();
  const { isAuthenticated, signOut } = useContext(AuthContext);

  return (
    <header className="sticky top-0 z-40 w-full h-16 bg-background/70 backdrop-blur-lg backdrop-saturate-150">
      <nav className="flex h-full min-w-0 items-center gap-2 px-4 sm:gap-4 sm:px-6">
        <Link
          to="/"
          className="flex shrink-0 items-center gap-2 text-foreground"
        >
          <img
            src="/makura.png"
            alt="MaCRa"
            width={40}
            height={40}
            className="h-10 w-10 shrink-0"
          />
          <span>MaCRa</span>
        </Link>
        <div className="flex min-w-0 grow items-center gap-2">
          <LinkButton to="/events" variant="ghost">
            <CalendarIcon className="w-6 h-6 shrink-0" />
            {t('header.events')}
          </LinkButton>
        </div>
        <div className="flex min-w-0 items-center gap-2">
          {isAuthenticated ? <UserMenu signOut={signOut} /> : <VisitorMenu />}
        </div>
      </nav>
    </header>
  );
}
