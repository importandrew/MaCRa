import { Button } from '@heroui/react';
import { useCallback } from 'react';
import { BarsIcon } from '../Icons/BarsIcon';

// ContextMenuButton is layered above SideMenu (z-40) and below modals and
// dropdowns (z-50). className sets only its position.
export default function ContextMenuButton({
  onPress,
  className = 'absolute',
}: {
  onPress: () => void;
  className?: string;
}) {
  const handleOnPress = useCallback(() => {
    onPress();
  }, [onPress]);

  return (
    <Button
      variant="light"
      radius="full"
      size="lg"
      isIconOnly
      onPress={handleOnPress}
      className={`z-45 ${className}`}
    >
      <BarsIcon className="text-gray-400" />
    </Button>
  );
}
