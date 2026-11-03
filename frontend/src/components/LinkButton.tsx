import { buttonVariants } from '@heroui/styles';
import { Link, LinkProps } from 'react-router-dom';

type ButtonVariant =
  | 'primary'
  | 'secondary'
  | 'tertiary'
  | 'outline'
  | 'ghost'
  | 'danger'
  | 'danger-soft';

export type LinkButtonProps = LinkProps & {
  variant?: ButtonVariant;
  size?: 'sm' | 'md' | 'lg';
  fullWidth?: boolean;
  isIconOnly?: boolean;
};

// LinkButton is a react-router link styled as a HeroUI button. HeroUI v3
// Button has no `as` prop, and a real <a> keeps long-press and "open in new
// tab" working on phones.
export function LinkButton({
  variant = 'primary',
  size,
  fullWidth,
  isIconOnly,
  className,
  ...props
}: LinkButtonProps) {
  return (
    <Link
      className={buttonVariants({
        variant,
        size,
        fullWidth,
        isIconOnly,
        className: typeof className === 'string' ? className : undefined,
      })}
      {...props}
    />
  );
}
