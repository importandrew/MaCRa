import { useTranslation } from 'react-i18next';
import { Card, useOverlayState } from '@heroui/react';
import { useRef } from 'react';
import type { TourState } from '../../models/Tour';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import {
  deleteSeatingPlan,
  generateSeatingPlan,
  importSeatingPlan,
} from '../../providers/event/tournamentApi';
import Confirm from '../Confirm';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';
import { getDateTimeFormatted } from '../../utils/date-tools';

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

// SeatingPlanCard is the seating prepared for all tours on the tournament's
// desk: what it is, and for a referee generating it (for the remaining
// tours), importing it from a CSV file, or deleting it. A tour takes its
// tables from it with the "fixed" method (ToursCard).
export default function SeatingPlanCard({
  eventId,
  state,
  isDrawn,
  canRun,
  onChanged,
}: {
  eventId: number;
  state: TourState;
  // isDrawn: the participant numbers are drawn for the current players.
  isDrawn: boolean;
  canRun: boolean;
  onChanged: () => Promise<unknown>;
}) {
  const { t } = useTranslation();
  const { isPending: isGenerating, run: runGenerate } = useSubmitOnce();
  const { isPending: isImporting, run: runImport } = useSubmitOnce();
  const { isPending: isDeleting, run: runDelete } = useSubmitOnce();
  const replaceConfirm = useOverlayState();
  const deleteConfirm = useOverlayState();
  const fileInput = useRef<HTMLInputElement>(null);
  const plan = state.plan;

  const generate = () =>
    runGenerate(async () => {
      replaceConfirm.close();
      if (await generateSeatingPlan(eventId)) {
        showToast({ title: t('seatingPlan.generated'), variant: 'success' });
      }
      await onChanged();
    });

  const importFile = (file: File) =>
    runImport(async () => {
      if (file.size > maxFileBytes) {
        showToast({ title: t('seatingPlan.tooLarge'), variant: 'danger' });
        return;
      }
      if (await importSeatingPlan(eventId, await file.text())) {
        showToast({ title: t('seatingPlan.imported'), variant: 'success' });
      }
      await onChanged();
    });

  const remove = () =>
    runDelete(async () => {
      deleteConfirm.close();
      if (await deleteSeatingPlan(eventId)) {
        showToast({ title: t('seatingPlan.deleted'), variant: 'success' });
      }
      await onChanged();
    });

  // Why the plan can't be generated now, if it can't (an import needs
  // neither: its numbers are checked when a tour is seated).
  const blocker = !isDrawn
    ? t('tours.needsDraw')
    : state.missing > 0
      ? t('tours.missing', { count: state.missing })
      : null;

  return (
    <Card>
      <Card.Content className="gap-3">
        <h2 className="font-medium">{t('seatingPlan.title')}</h2>
        {plan ? (
          <div className="text-sm">
            <p>
              {t(`seatingPlan.source.${plan.source}`, {
                date: getDateTimeFormatted(plan.created_at),
              })}
            </p>
            <p>
              {t('seatingPlan.summary', {
                tours: tourRanges(plan.tours),
                count: plan.tours.length,
                tables: plan.tables,
                players: plan.players,
              })}
            </p>
            <p className={plan.repeats > 0 ? 'text-warning' : 'text-muted'}>
              {plan.repeats > 0
                ? t('seatingPlan.repeats', {
                    count: plan.repeats,
                    maxMet: plan.max_met,
                  })
                : t('tours.noRepeats')}
            </p>
          </div>
        ) : (
          <p className="text-sm text-muted">{t('seatingPlan.none')}</p>
        )}

        {canRun && (
          <>
            <p className="text-sm text-muted">
              {blocker ?? t('seatingPlan.generateHint')}
            </p>
            <PendingButton
              variant={plan ? 'secondary' : 'primary'}
              isPending={isGenerating}
              isDisabled={blocker != null}
              // A new plan replaces the old one: confirm it.
              onPress={plan ? replaceConfirm.open : generate}
              className="w-3/4 self-center"
            >
              {t('seatingPlan.generate')}
            </PendingButton>
            <p className="text-sm text-muted">{t('seatingPlan.fileHint')}</p>
            <input
              ref={fileInput}
              type="file"
              accept=".csv,.txt,text/csv,text/plain"
              className="hidden"
              aria-label={t('seatingPlan.import')}
              onChange={(e) => {
                const file = e.target.files?.[0];
                // The same file may be picked again after a fix.
                e.target.value = '';
                if (file) importFile(file);
              }}
            />
            <PendingButton
              variant="secondary"
              isPending={isImporting}
              onPress={() => fileInput.current?.click()}
              className="w-3/4 self-center"
            >
              {t('seatingPlan.import')}
            </PendingButton>
            {plan && (
              <PendingButton
                variant="danger-soft"
                isPending={isDeleting}
                onPress={deleteConfirm.open}
                className="w-3/4 self-center"
              >
                {t('seatingPlan.delete')}
              </PendingButton>
            )}
          </>
        )}
      </Card.Content>
      <Confirm
        isOpen={replaceConfirm.isOpen}
        onClose={replaceConfirm.close}
        onConfirm={generate}
        headerText={t('seatingPlan.replaceTitle')}
        prompt={t('seatingPlan.replacePrompt')}
        cancelText={t('common.cancel')}
        confirmText={t('seatingPlan.generate')}
        confirmVariant="primary"
      />
      <Confirm
        isOpen={deleteConfirm.isOpen}
        onClose={deleteConfirm.close}
        onConfirm={remove}
        headerText={t('seatingPlan.deleteTitle')}
        prompt={t('seatingPlan.deletePrompt')}
        cancelText={t('common.cancel')}
        confirmText={t('seatingPlan.delete')}
        confirmVariant="danger"
      />
    </Card>
  );
}
