import { Outlet, useLocation } from 'react-router-dom';
import Header from '../components/Header/Header';
import ErrorBoundary from '../components/ErrorBoundary/ErrorBoundary';

export default function RootLayout() {
  const { pathname } = useLocation();

  return (
    <div className="flex flex-col h-dvh w-full max-w-full min-h-[550px]">
      <div className="w-full shrink-0">
        <Header />
      </div>
      {/* Below the header, so the header keeps working when a page crashes;
          a new path clears the error. */}
      <ErrorBoundary resetKey={pathname}>
        <Outlet />
      </ErrorBoundary>
    </div>
  );
}
