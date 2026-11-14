import { useTranslation } from 'react-i18next';
import { Button, Modal } from '@heroui/react';
import { useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { addTime } from '../../providers/event/tournamentApi';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';
import ValueStepper from '../ValueStepper';

// Limits of the time added at once, mirroring the backend
// (models.MinExtraMinutes, MaxExtraMinutes).
const MIN_MINUTES = 1;
const MAX_MINUTES = 60;

// AddTimeModal adds minutes to a table of the live tour (a referee, from the
// desk). The minutes start at 1: usually a table needs only a little.
export default function AddTimeModal({
  eventId,
  table,
  onClose,
  onAdded,
}: {
  eventId: number;
  // table: the table to add time to; null while the modal is closed.
  table: { gameId: number; number: number } | null;
  onClose: () => void;
  onAdded: () => Promise<unknown>;
}) {
  const { t } = useTranslation();
  const [minutes, setMinutes] = useState(MIN_MINUTES);
  const { isPending, run } = useSubmitOnce();
  const isValid = minutes >= MIN_MINUTES && minutes <= MAX_MINUTES;

  const close = () => {
    setMinutes(MIN_MINUTES);
    onClose();
  };

  const add = () =>
    run(async () => {
      if (!table) return;
      if (await addTime(eventId, table.gameId, minutes)) {
        showToast({
          title: t('tours.timeAdded', { n: minutes, table: table.number }),
          variant: 'success',
        });
        await onAdded();
        close();
      }
    });

  return (
    <Modal.Backdrop
      isOpen={table != null}
      onOpenChange={(open) => {
        if (!open) close();
      }}
    >
      <Modal.Container placement="center">
        <Modal.Dialog>
          <Modal.CloseTrigger aria-label={t('common.close')} />
          <Modal.Header>
            <Modal.Heading>
              {t('tours.addTimeTitle', { table: table?.number ?? 0 })}
            </Modal.Heading>
          </Modal.Header>
          <Modal.Body className="flex flex-col gap-2">
            <ValueStepper
              label={t('tours.minutes')}
              value={minutes}
              min={MIN_MINUTES}
              onChange={setMinutes}
            />
            <p
              className={`text-center text-xs ${
                minutes > MAX_MINUTES ? 'text-danger' : 'text-muted'
              }`}
            >
              {t('tours.minutesRange', { min: MIN_MINUTES, max: MAX_MINUTES })}
            </p>
          </Modal.Body>
          <Modal.Footer>
            <Button slot="close" variant="tertiary">
              {t('common.cancel')}
            </Button>
            <PendingButton
              variant="primary"
              isPending={isPending}
              isDisabled={!isValid}
              onPress={add}
            >
              {t('tours.addMinutes', { n: minutes })}
            </PendingButton>
          </Modal.Footer>
        </Modal.Dialog>
      </Modal.Container>
    </Modal.Backdrop>
  );
}
