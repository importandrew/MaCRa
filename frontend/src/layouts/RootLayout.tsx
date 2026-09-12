import { Outlet } from 'react-router-dom';
import Header from '../components/Header/Header';

export default function RootLayout() {
  return (
    <div className="flex flex-col h-dvh w-full max-w-full min-h-[550px]">
      <div className="w-full shrink-0">
        <Header />
      </div>
      <Outlet />
    </div>
  );
}
