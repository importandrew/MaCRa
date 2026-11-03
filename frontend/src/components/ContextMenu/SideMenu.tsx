import { Drawer } from '@heroui/react';
import { ReactNode } from 'react';
import ContextMenuButton from './ContextMenuButton';

// ButtonRect is where the ≡ button sits, as distances from the right and
// bottom edges of the viewport.
export type ButtonRect = {
  right: number;
  bottom: number;
  width: number;
  height: number;
};

// SideMenu is the drawer behind the ≡ button. While it's open, React Aria makes
// everything outside it inert, so the ≡ button that opened it can't be tapped.
// SideMenu draws its own ≡ button at the same screen position (closeButtonRect)
// so a second tap on ≡ still closes the menu. On phones the menu takes the full
// screen width, which leaves no backdrop to tap.
// The panel is pinned to the right and bottom edges of the viewport, so the
// close button is positioned from those edges. position: fixed wouldn't work:
// the panel has a transform, which makes it the containing block.
export default function SideMenu({
  isOpen,
  onClose,
  closeButtonRect,
  children,
}: {
  isOpen: boolean;
  onClose: () => void;
  closeButtonRect: ButtonRect | null;
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
          className="relative w-full max-w-full sm:w-full sm:max-w-3xl"
        >
          {children}
          {closeButtonRect && (
            <ContextMenuButton
              ariaLabel="Close menu"
              onPress={onClose}
              className="absolute z-10"
              style={{
                right: closeButtonRect.right,
                bottom: closeButtonRect.bottom,
                width: closeButtonRect.width,
                height: closeButtonRect.height,
              }}
            />
          )}
        </Drawer.Dialog>
      </Drawer.Content>
    </Drawer.Backdrop>
  );
}
