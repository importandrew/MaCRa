import { Outlet } from 'react-router-dom';

/** Centered / w-screen-friendly shell for game UI and landing. */
export default function AssistantLayout() {
  return (
    <div className="flex flex-col grow items-center my-auto min-h-0 w-full">
      <main className="flex grow light text-foreground bg-background">
        <Outlet />
      </main>
    </div>
  );
}
