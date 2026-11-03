import { AlertDialog, Button } from '@heroui/react';

export interface ConfirmProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  prompt: string;
  headerText: string;
  cancelText: string;
  confirmText: string;
  confirmVariant: 'primary' | 'danger';
}

// Confirm is an alert dialog: a tap on the backdrop or ESC doesn't dismiss it,
// so a stray tap can't silently cancel actions like "Abort game".
export default function Confirm(props: ConfirmProps) {
  return (
    <AlertDialog.Backdrop
      isOpen={props.isOpen}
      onOpenChange={(isOpen) => {
        if (!isOpen) props.onClose();
      }}
    >
      <AlertDialog.Container placement="center">
        <AlertDialog.Dialog className="sm:max-w-[400px]">
          <AlertDialog.Header>
            <AlertDialog.Icon
              status={props.confirmVariant === 'danger' ? 'danger' : 'accent'}
            />
            <AlertDialog.Heading>{props.headerText}</AlertDialog.Heading>
          </AlertDialog.Header>
          <AlertDialog.Body>
            <p>{props.prompt}</p>
          </AlertDialog.Body>
          <AlertDialog.Footer>
            <Button slot="close" variant="tertiary">
              {props.cancelText}
            </Button>
            <Button variant={props.confirmVariant} onPress={props.onConfirm}>
              {props.confirmText}
            </Button>
          </AlertDialog.Footer>
        </AlertDialog.Dialog>
      </AlertDialog.Container>
    </AlertDialog.Backdrop>
  );
}
