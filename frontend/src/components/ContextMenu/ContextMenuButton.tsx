import { Button } from '@heroui/react';
import { useCallback } from 'react';
import { BarsIcon } from '../Icons/BarsIcon';

export default function ContextMenuButton({
  onPress,
  className = 'absolute z-100',
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
      className={className}
    >
      <BarsIcon className="text-gray-400" />
    </Button>
  );
}
