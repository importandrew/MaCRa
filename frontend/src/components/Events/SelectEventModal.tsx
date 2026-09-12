import {
  Modal,
  ModalBody,
  ModalContent,
  ModalHeader,
  Spinner,
} from '@heroui/react';
import { useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { useFetchVisibleEvents } from '../../providers/event/eventApi';
import { useEvent } from '../../providers/event/eventConstants';

type SelectEventModalProps = {
  isOpen: boolean;
  onClose: () => void;
};

export default function SelectEventModal({
  isOpen,
  onClose,
}: SelectEventModalProps) {
  const navigate = useNavigate();
  const { selectedEventId, selectEvent } = useEvent();
  const { data, isLoading, isError, refetch } = useFetchVisibleEvents();

  useEffect(() => {
    if (isOpen) {
      refetch();
    }
  }, [isOpen, refetch]);

  const participatingEvents = useMemo(
    () => data?.filter((e) => e.is_participating) ?? [],
    [data]
  );

  const handleSelect = async (id: number) => {
    await selectEvent(id);
    onClose();
    navigate('/');
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} placement="center" size="lg">
      <ModalContent>
        <ModalHeader>Select event</ModalHeader>
        <ModalBody className="max-h-[70vh] overflow-y-auto gap-3 pb-6">
          {isLoading && (
            <div className="flex justify-center py-8">
              <Spinner label="Loading events..." />
            </div>
          )}
          {!isLoading && isError && (
            <p className="text-center text-danger py-6">
              Failed to load events.
            </p>
          )}
          {!isLoading && !isError && participatingEvents.length === 0 && (
            <p className="text-center text-default-500 py-6">
              You are not participating in any events.
            </p>
          )}
          {!isLoading &&
            !isError &&
            participatingEvents.map((event) => {
              const isSelected = selectedEventId === event.id;
              return (
                <button
                  key={event.id}
                  type="button"
                  onClick={() => handleSelect(event.id)}
                  className={`w-full text-left rounded-lg border px-4 py-3 transition-colors hover:bg-default-100 ${
                    isSelected
                      ? 'border-primary bg-primary-50'
                      : 'border-default-200'
                  }`}
                >
                  <div className="font-medium break-words">{event.name}</div>
                  {event.description && (
                    <p className="text-sm text-default-500 mt-1 line-clamp-3 break-words">
                      {event.description}
                    </p>
                  )}
                </button>
              );
            })}
        </ModalBody>
      </ModalContent>
    </Modal>
  );
}
