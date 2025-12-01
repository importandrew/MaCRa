import { Drawer, DrawerContent } from '@heroui/react';
import { ReactNode } from 'react';

export default function SideMenu({
  isOpen,
  children,
  onOpenChange,
}: {
  isOpen: boolean;
  children: ReactNode;
  onOpenChange: () => void;
}) {
  return (
    <Drawer size="3xl" isOpen={isOpen} onOpenChange={onOpenChange}>
      <DrawerContent>{children}</DrawerContent>
    </Drawer>
  );
}
