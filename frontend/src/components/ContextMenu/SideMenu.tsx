import { Drawer } from '@heroui/react';
import { ReactNode } from 'react';

// SideMenu is the drawer behind the ≡ button. On phones it takes the full
// screen width (the logs and fans tables need the room), so there's no
// backdrop to tap: it closes with the ✕ button or by swiping its header to
// the right (built into the v3 Drawer). Escape closes it on desktop.
export default function SideMenu({
  isOpen,
  onClose,
  children,
}: {
  isOpen: boolean;
  onClose: () => void;
  children: ReactNode;
}) {
  return (
    <Drawer.Backdrop
      isOpen={isOpen}
      onOpenChange={(open) => {
        if (!open) onClose();
      }}
    >
      <Drawer.Content placement="right">
        <Drawer.Dialog
          aria-label="Menu"
          className="w-full max-w-full sm:w-full sm:max-w-3xl"
        >
          <Drawer.CloseTrigger aria-label="Close menu" />
          {children}
        </Drawer.Dialog>
      </Drawer.Content>
    </Drawer.Backdrop>
  );
}
