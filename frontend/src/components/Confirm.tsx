import {
  Button,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
} from '@heroui/react';

export interface ConfirmProps {
  isOpen: boolean;
  onClose: any;
  onConfirm: () => void;
  prompt: string;
  cancelText: string;
  confirmText: string;
}

export default function Confirm(props: ConfirmProps) {
  return (
    <Modal isOpen={props.isOpen} onClose={props.onClose} placement="center">
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader className="flex flex-col gap-1">
              Confirm Abort Game
            </ModalHeader>
            <ModalBody>
              <p>{props.prompt}</p>
            </ModalBody>
            <ModalFooter>
              <Button color="default" variant="light" onPress={onClose}>
                {props.cancelText}
              </Button>
              <Button color="danger" onPress={props.onConfirm}>
                {props.confirmText}
              </Button>
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>
  );
}
