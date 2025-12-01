import { ReactNode } from 'react';
import Header from './components/Header/Header';

const Layout = ({ children }: { children: ReactNode }) => {
  return (
    <div className="flex flex-col h-dvh text-5xl items-center my-auto">
      <div className="flex w-screen">
        <Header />
      </div>
      <div className="flex grow">{children}</div>
    </div>
  );
};

export default Layout;
