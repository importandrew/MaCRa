import { Modal } from '@heroui/react';
import { useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { useFetchVisibleEvents } from '../../providers/event/eventApi';
import { useEvent } from '../../providers/event/eventConstants';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';

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
    <Modal.Backdrop
      isOpen={isOpen}
      onOpenChange={(open) => {
        if (!open) onClose();
      }}
    >
      <Modal.Container placement="center" size="lg" scroll="inside">
        <Modal.Dialog>
          <Modal.CloseTrigger />
          <Modal.Header>
            <Modal.Heading>Select event</Modal.Heading>
          </Modal.Header>
          <Modal.Body className="flex flex-col gap-3 pb-6">
            {isLoading && (
              <div className="flex justify-center py-8">
                <LoadingSpinner label="Loading events..." size="md" />
              </div>
            )}
            {!isLoading && isError && (
              <p className="text-center text-danger py-6">
                Failed to load events.
              </p>
            )}
            {!isLoading && !isError && participatingEvents.length === 0 && (
              <p className="text-center text-muted py-6">
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
                    className={`w-full text-left rounded-lg border px-4 py-3 transition-colors hover:bg-default ${
                      isSelected
                        ? 'border-accent bg-accent-soft'
                        : 'border-border'
                    }`}
                  >
                    <div className="font-medium break-words">{event.name}</div>
                    {event.description && (
                      <p className="text-sm text-muted mt-1 line-clamp-3 break-words">
                        {event.description}
                      </p>
                    )}
                  </button>
                );
              })}
          </Modal.Body>
        </Modal.Dialog>
      </Modal.Container>
    </Modal.Backdrop>
  );
}
