import { Button, Spinner, type ButtonProps } from '@heroui/react';
import { ReactNode } from 'react';

export type PendingButtonProps = Omit<ButtonProps, 'children'> & {
  children: ReactNode;
  // icon shows before the label while the button isn't pending.
  icon?: ReactNode;
};

// PendingButton is a Button that shows a spinner before its label while
// isPending is set (HeroUI v3 leaves the pending look to the caller).
export default function PendingButton({
  children,
  icon,
  ...props
}: PendingButtonProps) {
  return (
    <Button {...props}>
      {({ isPending }) => (
        <>
          {isPending ? <Spinner color="current" size="sm" /> : icon}
          {children}
        </>
      )}
    </Button>
  );
}
