import { useTranslation } from 'react-i18next';
import { Button, Card, Modal, useOverlayState } from '@heroui/react';
import { useRef, useState } from 'react';
import type { Seating, TourPlan, TourState } from '../../models/Tour';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import {
  deleteSeatingFile,
  importSeatingFile,
  setTourPlan,
} from '../../providers/event/tournamentApi';
import Confirm from '../Confirm';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';
import SeatingFields from './SeatingFields';
import { seatingLabel } from './seatingLabel';

// maxFileBytes bounds a seating file read in the browser (the server takes
// up to 256 KB of JSON).
const maxFileBytes = 200 * 1024;

// tourRanges writes tour numbers as ranges: [1, 2, 3, 5] → "1–3, 5".
function tourRanges(tours: number[]): string {
  const ranges: string[] = [];
  for (let i = 0; i < tours.length;) {
    let j = i;
    while (j + 1 < tours.length && tours[j + 1] === tours[j] + 1) j++;
    ranges.push(i === j ? `${tours[i]}` : `${tours[i]}–${tours[j]}`);
    i = j + 1;
  }
  return ranges.join(', ');
}

// TourPlanCard is the tournament's plan on its desk: each tour's seating
// method (the next tour's seating starts from it), and the tables imported
// from a file that the "fixed" method takes. A referee changes a tour's
// method until its seating is published, imports a file or deletes it.
export default function TourPlanCard({
  eventId,
  state,
  canRun,
  onChanged,
}: {
  eventId: number;
  state: TourState;
  canRun: boolean;
  onChanged: () => Promise<unknown>;
}) {
  const { t } = useTranslation();
  const { isPending: isSaving, run: runSave } = useSubmitOnce();
  const { isPending: isImporting, run: runImport } = useSubmitOnce();
  const { isPending: isDeleting, run: runDelete } = useSubmitOnce();
  const deleteConfirm = useOverlayState();
  const fileInput = useRef<HTMLInputElement>(null);
  // editing: the tour whose plan the modal changes, and the seating picked
  // there; null while it is closed.
  const [editing, setEditing] = useState<{
    tour: TourPlan;
    seating: Seating;
  } | null>(null);
  const tableCount = Math.max(1, Math.floor(state.active / 4));
  const file = state.file;

  const edit = (tour: TourPlan) =>
    setEditing({
      tour,
      seating: {
        seating: tour.seating,
        step: tour.seating_step ?? undefined,
        avoid_repeats: tour.avoid_repeats,
      },
    });

  const save = () =>
    runSave(async () => {
      if (!editing) return;
      const { seating } = editing;
      const step =
        seating.seating === 'interval' &&
        seating.step != null &&
        seating.step < tableCount
          ? seating.step
          : undefined;
      if (
        await setTourPlan(eventId, editing.tour.number, { ...seating, step })
      ) {
        showToast({
          title: t('tourPlan.saved', { n: editing.tour.number }),
          variant: 'success',
        });
        setEditing(null);
      }
      await onChanged();
    });

  const importFile = (picked: File) =>
    runImport(async () => {
      if (picked.size > maxFileBytes) {
        showToast({ title: t('tourPlan.tooLarge'), variant: 'danger' });
        return;
      }
      if (await importSeatingFile(eventId, await picked.text())) {
        showToast({ title: t('tourPlan.imported'), variant: 'success' });
      }
      await onChanged();
    });

  const remove = () =>
    runDelete(async () => {
      deleteConfirm.close();
      if (await deleteSeatingFile(eventId)) {
        showToast({ title: t('tourPlan.deleted'), variant: 'success' });
      }
      await onChanged();
    });

  return (
    <Card>
      <Card.Content className="gap-3">
        <h2 className="font-medium">{t('tourPlan.title')}</h2>
        <p className="text-sm text-muted">{t('tourPlan.hint')}</p>
        <ul aria-label={t('tourPlan.title')} className="flex flex-col">
          {state.plan.map((tour) => (
            <li
              key={tour.number}
              className="flex min-h-9 items-center gap-2 border-b border-border last:border-b-0"
            >
              <span className="w-16 shrink-0 text-sm">
                {t('tourPlan.tour', { n: tour.number })}
              </span>
              <span
                className={`min-w-0 flex-1 text-sm ${tour.seated ? 'text-muted' : ''}`}
              >
                {seatingLabel(
                  t,
                  tour.seating,
                  tour.seating_step,
                  tour.avoid_repeats
                )}
                {tour.seated && ` · ${t('tourPlan.seated')}`}
              </span>
              {canRun && !tour.seated && (
                <Button
                  size="sm"
                  variant="tertiary"
                  aria-label={t('tourPlan.editTour', { n: tour.number })}
                  onPress={() => edit(tour)}
                >
                  {t('tourPlan.edit')}
                </Button>
              )}
            </li>
          ))}
        </ul>

        <h3 className="text-sm font-medium">{t('tourPlan.fileTitle')}</h3>
        {file ? (
          <div className="text-sm">
            <p>
              {t('tourPlan.fileSummary', {
                tours: tourRanges(file.tours),
                count: file.tours.length,
                tables: file.tables,
                players: file.players,
              })}
            </p>
            <p className={file.repeats > 0 ? 'text-warning' : 'text-muted'}>
              {file.repeats > 0
                ? t('tourPlan.fileRepeats', {
                    count: file.repeats,
                    maxMet: file.max_met,
                  })
                : t('tours.noRepeats')}
            </p>
          </div>
        ) : null}
        <p className="text-sm text-muted">{t('tourPlan.fileHint')}</p>
        {canRun && (
          <>
            <input
              ref={fileInput}
              type="file"
              accept=".csv,.txt,text/csv,text/plain"
              className="hidden"
              aria-label={t('tourPlan.import')}
              onChange={(e) => {
                const picked = e.target.files?.[0];
                // The same file may be picked again after a fix.
                e.target.value = '';
                if (picked) importFile(picked);
              }}
            />
            <PendingButton
              variant="secondary"
              isPending={isImporting}
              onPress={() => fileInput.current?.click()}
              className="w-3/4 self-center"
            >
              {file ? t('tourPlan.reimport') : t('tourPlan.import')}
            </PendingButton>
            {file && (
              <PendingButton
                variant="danger-soft"
                isPending={isDeleting}
                onPress={deleteConfirm.open}
                className="w-3/4 self-center"
              >
                {t('tourPlan.delete')}
              </PendingButton>
            )}
          </>
        )}
      </Card.Content>

      <Modal.Backdrop
        isOpen={editing != null}
        onOpenChange={(open) => {
          if (!open) setEditing(null);
        }}
      >
        <Modal.Container placement="center">
          <Modal.Dialog>
            <Modal.CloseTrigger aria-label={t('common.close')} />
            <Modal.Header>
              <Modal.Heading>
                {t('tourPlan.editTitle', { n: editing?.tour.number ?? 0 })}
              </Modal.Heading>
            </Modal.Header>
            <Modal.Body className="flex flex-col gap-3">
              {editing && (
                <SeatingFields
                  value={editing.seating}
                  onChange={(seating) => setEditing({ ...editing, seating })}
                  tableCount={tableCount}
                  hasFile={editing.tour.file_tables > 0}
                />
              )}
            </Modal.Body>
            <Modal.Footer>
              <Button slot="close" variant="tertiary">
                {t('common.cancel')}
              </Button>
              <PendingButton
                variant="primary"
                isPending={isSaving}
                onPress={save}
              >
                {t('common.save')}
              </PendingButton>
            </Modal.Footer>
          </Modal.Dialog>
        </Modal.Container>
      </Modal.Backdrop>
      <Confirm
        isOpen={deleteConfirm.isOpen}
        onClose={deleteConfirm.close}
        onConfirm={remove}
        headerText={t('tourPlan.deleteTitle')}
        prompt={t('tourPlan.deletePrompt')}
        cancelText={t('common.cancel')}
        confirmText={t('tourPlan.delete')}
        confirmVariant="danger"
      />
    </Card>
  );
}
