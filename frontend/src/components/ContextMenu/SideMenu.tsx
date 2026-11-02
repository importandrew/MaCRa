import { Drawer, DrawerContent } from '@heroui/react';
import { ReactNode } from 'react';

// SideMenu sits below the default overlay layer (z-50) so modals and dropdowns
// opened from it still cover it, while ContextMenuButton (z-45) stays on top of
// it and can close it. On phones the menu takes the full screen width, so the
// button is the easiest way to close it.
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
    <Drawer
      size="3xl"
      isOpen={isOpen}
      onOpenChange={onOpenChange}
      classNames={{ wrapper: 'z-40', backdrop: 'z-40' }}
    >
      <DrawerContent>{children}</DrawerContent>
    </Drawer>
  );
}
