import { Button } from '@heroui/react';
import { MainButtonProps } from '../../models/Button';
import { buttonTint } from '../buttonTint';
import { LinkButton } from '../LinkButton';

function MainButton({
  text,
  link = '/',
  onPress,
  tint = 'accent',
}: MainButtonProps) {
  const variant = tint === 'default' ? 'tertiary' : 'secondary';
  const className = tint === 'default' ? undefined : buttonTint(tint);

  if (onPress) {
    return (
      <Button
        variant={variant}
        size="lg"
        fullWidth
        className={className}
        onPress={onPress}
      >
        {text}
      </Button>
    );
  }
  return (
    <LinkButton
      to={link}
      variant={variant}
      size="lg"
      fullWidth
      className={className}
    >
      {text}
    </LinkButton>
  );
}

export default MainButton;
