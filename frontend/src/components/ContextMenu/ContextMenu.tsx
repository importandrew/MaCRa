import { ReactNode, useLayoutEffect, useRef, useState } from 'react';
import ContextMenuButton from './ContextMenuButton';
import SideMenu, { ButtonRect } from './SideMenu';

// ContextMenu is the ≡ button plus the SideMenu it opens. Every page with a ≡
// menu uses it; buttonClassName positions the button on that page.
export default function ContextMenu({
  children,
  isOpen,
  onOpen,
  onClose,
  buttonClassName = 'absolute',
}: {
  children: ReactNode;
  isOpen: boolean;
  onOpen: () => void;
  onClose: () => void;
  buttonClassName?: string;
}) {
  const buttonRef = useRef<HTMLButtonElement>(null);
  const [buttonRect, setButtonRect] = useState<ButtonRect | null>(null);

  // Remember where the ≡ button is whenever the menu opens, whether the button
  // or the page opened it, so SideMenu can draw its close button on top of it.
  useLayoutEffect(() => {
    if (!isOpen) return;
    const rect = buttonRef.current?.getBoundingClientRect();
    if (!rect) return;
    setButtonRect({
      right: window.innerWidth - rect.right,
      bottom: window.innerHeight - rect.bottom,
      width: rect.width,
      height: rect.height,
    });
  }, [isOpen]);

  return (
    <>
      <ContextMenuButton
        ref={buttonRef}
        ariaLabel="Open menu"
        onPress={onOpen}
        className={buttonClassName}
      />
      <SideMenu isOpen={isOpen} onClose={onClose} closeButtonRect={buttonRect}>
        {children}
      </SideMenu>
    </>
  );
}
