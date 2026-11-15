import { useTranslation } from 'react-i18next';
import { Switch } from '@heroui/react';
import { SEATING_METHODS, type Seating } from '../../models/Tour';
import SegmentedPicker from '../SegmentedPicker';
import ListStepper from '../ListStepper';
import { defaultAvoidRepeats } from './seatingLabel';

// SeatingFields picks how a tour is seated: the method (fixed only when the
// tour has tables from a file), an interval's step, and for Swiss and
// interval whether to avoid repeat meetings. Used for the next tour and in
// the tournament's plan.
export default function SeatingFields({
  value,
  onChange,
  tableCount,
  hasFile,
}: {
  value: Seating;
  onChange: (value: Seating) => void;
  // tableCount: the tables of the current participants, for the step.
  tableCount: number;
  // hasFile: the tour has tables from a file (the "fixed" method).
  hasFile: boolean;
}) {
  const { t } = useTranslation();
  const method = value.seating;

  return (
    <>
      <SegmentedPicker
        ariaLabel={t('tours.methodLabel')}
        // Four one-word labels don't fit a row at 320 px in Russian.
        orientation="vertical"
        options={SEATING_METHODS.filter(
          (option) => option !== 'fixed' || hasFile
        ).map((option) => ({
          value: option,
          label: t(`tours.method.${option}`),
        }))}
        selectedKey={method}
        onSelectionChange={(seating) =>
          onChange({ seating, avoid_repeats: defaultAvoidRepeats(seating) })
        }
      />
      {method === 'interval' && tableCount > 1 && (
        <>
          <p className="text-sm">{t('tours.step')}</p>
          <ListStepper
            label={t('tours.step')}
            values={Array.from({ length: tableCount }, (_, i) => i + 1)}
            value={Math.min(value.step ?? tableCount, tableCount)}
            onChange={(step) =>
              onChange({
                ...value,
                step: step === tableCount ? undefined : step,
              })
            }
          />
        </>
      )}
      {(method === 'swiss' || method === 'interval') && (
        <Switch
          size="sm"
          isSelected={value.avoid_repeats}
          onChange={(avoid_repeats) => onChange({ ...value, avoid_repeats })}
        >
          <Switch.Content>
            <Switch.Control>
              <Switch.Thumb />
            </Switch.Control>
            {t('tours.avoidRepeats')}
          </Switch.Content>
        </Switch>
      )}
    </>
  );
}
