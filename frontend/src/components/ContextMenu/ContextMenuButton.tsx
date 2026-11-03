import { Button } from '@heroui/react';
import { CSSProperties, Ref } from 'react';
import { BarsIcon } from '../Icons/BarsIcon';

// ContextMenuButton is the round ≡ button that opens and closes a SideMenu.
// className sets only its position.
export default function ContextMenuButton({
  onPress,
  className = '',
  style,
  ref,
  ariaLabel,
}: {
  onPress: () => void;
  className?: string;
  style?: CSSProperties;
  ref?: Ref<HTMLButtonElement>;
  ariaLabel: string;
}) {
  return (
    <Button
      ref={ref}
      variant="ghost"
      size="lg"
      isIconOnly
      aria-label={ariaLabel}
      onPress={onPress}
      className={`rounded-full ${className}`}
      style={style}
    >
      <BarsIcon className="text-gray-400" />
    </Button>
  );
}
