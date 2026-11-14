import { Button } from '@heroui/react';

// ListStepper picks one of a fixed, ordered list of numbers (e.g. the
// official penalty amounts) with − and + buttons that move to the previous
// or next one, so game screens need no keyboard. At either end of the list
// its button is disabled.
export default function ListStepper({
  label,
  values,
  value,
  onChange,
}: {
  label: string;
  values: number[];
  value: number;
  onChange: (value: number) => void;
}) {
  const index = values.indexOf(value);
  const step = (by: -1 | 1) => {
    const next = values[index + by];
    if (next != null) onChange(next);
  };

  // One line, "− 10 +": the number is plain text, not a field, since it
  // can't be typed. 44 + 48 + 44 px and two 8 px gaps fit the 152 px center
  // of the table on a 320px phone.
  return (
    <div className="flex items-center justify-center gap-2">
      <Button
        size="sm"
        variant="tertiary"
        className="min-w-11 px-0 text-lg md:h-10 md:min-w-12"
        aria-label={`${label} −`}
        isDisabled={index <= 0}
        onPress={() => step(-1)}
      >
        −
      </Button>
      {/* The number is read out when it changes; the label names it. */}
      <output
        aria-label={label}
        aria-live="polite"
        className="min-w-12 text-center text-2xl font-semibold tabular-nums"
      >
        {value}
      </output>
      <Button
        size="sm"
        variant="tertiary"
        className="min-w-11 px-0 text-lg md:h-10 md:min-w-12"
        aria-label={`${label} +`}
        isDisabled={index < 0 || index >= values.length - 1}
        onPress={() => step(1)}
      >
        +
      </Button>
    </div>
  );
}
