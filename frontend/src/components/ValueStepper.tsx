import { Button, FieldError, Input, TextField } from '@heroui/react';
import { useState } from 'react';

const STEPS = [-10, -1, 1, 10];

// ValueStepper sets a number with −10/−1/+1/+10 buttons, so game screens
// rarely need the phone keyboard; the number itself is a numeric field for
// values far from the start, e.g. 88. A TextField, not NumberField:
// NumberField commits only on blur, so tapping OK right after typing would
// submit the old value. Values below min are shown as an error, not clamped;
// the caller keeps its submit disabled meanwhile.
export default function ValueStepper({
  label,
  value,
  min,
  onChange,
}: {
  label: string;
  value: number;
  min: number;
  onChange: (value: number) => void;
}) {
  // draft is the field's text: it may be empty while typing, which value
  // (a number) can't show. It follows value when that changes from outside
  // (the buttons, a reset).
  const [draft, setDraft] = useState(String(value));
  const [shownValue, setShownValue] = useState(value);
  const [isFocused, setIsFocused] = useState(false);

  if (value !== shownValue) {
    setShownValue(value);
    setDraft(String(value));
  }

  const type = (text: string) => {
    const digits = text.replace(/\D/g, '');
    const typed = digits === '' ? 0 : Number(digits);
    setDraft(digits);
    setShownValue(typed);
    onChange(typed);
  };

  // The error waits for blur, so typing 16 doesn't flash it at "1".
  const isTooLow = value < min;

  // The number sits above the buttons, not between them, and has no visible
  // caption (the field's aria-label names it): the stepper must fit the
  // center of the table on a 360px phone.
  return (
    <div className="flex flex-col items-center gap-1">
      <TextField
        aria-label={label}
        className="flex flex-col items-center gap-1"
        value={draft}
        onChange={type}
        onFocus={() => setIsFocused(true)}
        onBlur={() => setIsFocused(false)}
        isInvalid={isTooLow && !isFocused}
      >
        <Input
          inputMode="numeric"
          pattern="\d*"
          autoComplete="off"
          onFocus={(e) => e.currentTarget.select()}
          className="h-9 w-24 text-center text-2xl font-semibold tabular-nums"
        />
        <FieldError className="text-center">
          {`${label} should be ${min} or higher`}
        </FieldError>
      </TextField>
      <div className="flex items-center justify-center gap-1">
        {STEPS.map((step) => (
          <Button
            key={step}
            size="sm"
            variant="tertiary"
            className="min-w-10 px-1.5 md:h-10 md:min-w-12 md:text-base"
            aria-label={`${label} ${step > 0 ? '+' : ''}${step}`}
            isDisabled={step < 0 && value <= min}
            onPress={() => onChange(Math.max(min, value + step))}
          >
            {step > 0 ? `+${step}` : `−${-step}`}
          </Button>
        ))}
      </div>
    </div>
  );
}
