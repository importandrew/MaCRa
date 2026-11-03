import { Button } from '@heroui/react';
import { BarsIcon } from '../Icons/BarsIcon';

// ContextMenuButton is the round ≡ button that opens a SideMenu.
// className sets only its position.
export default function ContextMenuButton({
  onPress,
  className = '',
  ariaLabel,
}: {
  onPress: () => void;
  className?: string;
  ariaLabel: string;
}) {
  return (
    <Button
      variant="ghost"
      size="lg"
      isIconOnly
      aria-label={ariaLabel}
      onPress={onPress}
      className={`rounded-full ${className}`}
    >
      <BarsIcon className="text-gray-400" />
    </Button>
  );
}
