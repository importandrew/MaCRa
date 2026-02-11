import {
  Button,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
} from '@heroui/react';

type ButtonColor =
  | 'default'
  | 'primary'
  | 'secondary'
  | 'success'
  | 'warning'
  | 'danger';

export interface ConfirmProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  prompt: string;
  headerText: string;
  cancelText: string;
  confirmText: string;
  confirmButtonColor: ButtonColor;
}

export default function Confirm(props: ConfirmProps) {
  return (
    <Modal isOpen={props.isOpen} onClose={props.onClose} placement="center">
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader className="flex flex-col gap-1">
              {props.headerText}
            </ModalHeader>
            <ModalBody>
              <p>{props.prompt}</p>
            </ModalBody>
            <ModalFooter>
              <Button color="default" variant="light" onPress={onClose}>
                {props.cancelText}
              </Button>
              <Button
                color={props.confirmButtonColor}
                onPress={props.onConfirm}
              >
                {props.confirmText}
              </Button>
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>
  );
}
