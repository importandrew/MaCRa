import { ReactNode } from 'react';
import ContextMenuButton from './ContextMenuButton';
import SideMenu from './SideMenu';

export default function ContextMenu({
  children,
  isOpen,
  onClose,
  setIsOpen,
}: {
  children: ReactNode;
  isOpen: boolean;
  onClose: () => void;
  setIsOpen: (value: boolean) => void;
}) {
  const onOpen = () => {
    setIsOpen(!isOpen);
  };

  return (
    <>
      <ContextMenuButton onPress={onOpen} />
      <SideMenu isOpen={isOpen} onOpenChange={onClose}>
        {children}
      </SideMenu>
    </>
  );
}
