import { Button } from '@heroui/react';
import { MainButtonProps } from '../../models/Button';
import { LinkButton } from '../LinkButton';

function MainButton({
  text,
  link = '/',
  onPress,
  variant = 'primary',
}: MainButtonProps) {
  if (onPress) {
    return (
      <Button variant={variant} size="lg" fullWidth onPress={onPress}>
        {text}
      </Button>
    );
  }
  return (
    <LinkButton to={link} variant={variant} size="lg" fullWidth>
      {text}
    </LinkButton>
  );
}

export default MainButton;
