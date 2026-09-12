import { Link } from 'react-router-dom';
import { MainButtonProps } from '../../models/Button';
import { Button } from '@heroui/react';

function MainButton({
  text,
  link = '/',
  onPress,
  color = 'primary',
}: MainButtonProps) {
  return (
    <Button
      {...(onPress ? { onPress } : { as: Link, to: link })}
      color={color}
      variant="flat"
      size="lg"
    >
      {text}
    </Button>
  );
}

export default MainButton;
