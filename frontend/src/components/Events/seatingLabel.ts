import type { TFunction } from 'i18next';
import type { SeatingMethod } from '../../models/Tour';

// defaultAvoidRepeats: a Swiss seating avoids repeat meetings unless told
// otherwise, an interval one keeps to its scheme (the owner's choice).
export function defaultAvoidRepeats(method: SeatingMethod): boolean {
  return method === 'swiss';
}

// seatingLabel words a tour's seating: "Interval, step 4, avoiding repeats".
export function seatingLabel(
  t: TFunction,
  seating: SeatingMethod,
  step: number | null | undefined,
  avoidRepeats: boolean
): string {
  let label = t(`tours.method.${seating}`);
  if (step != null) label += t('tours.stepPart', { step });
  if (avoidRepeats) label += t('tours.avoidPart');
  return label;
}
