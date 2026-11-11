import { useTranslation } from 'react-i18next';
import { ReactNode } from 'react';
import ContextMenuButton from './ContextMenuButton';
import SideMenu from './SideMenu';

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
  const { t } = useTranslation();
  return (
    <>
      <ContextMenuButton
        ariaLabel={t('common.openMenu')}
        onPress={onOpen}
        className={buttonClassName}
      />
      <SideMenu isOpen={isOpen} onClose={onClose}>
        {children}
      </SideMenu>
    </>
  );
}
