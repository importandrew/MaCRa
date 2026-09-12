import { Outlet } from 'react-router-dom';

/** Full-width shell without horizontal overflow for event browse/manage pages. */
export default function EventPagesLayout() {
  return (
    <div className="flex flex-col grow w-full min-w-0 overflow-x-hidden">
      <main className="flex flex-col grow w-full min-w-0 light text-foreground bg-background">
        <Outlet />
      </main>
    </div>
  );
}
