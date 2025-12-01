import { Button } from '@heroui/react';
import { useCallback } from 'react';
import { BarsIcon } from '../Icons/BarsIcon';

export default function ContextMenuButton({
  onPress,
}: {
  onPress: () => void;
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
      className="absolute z-100"
    >
      <BarsIcon className="text-gray-400" />
    </Button>
  );
}
