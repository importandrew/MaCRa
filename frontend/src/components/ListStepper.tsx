import { Button } from '@heroui/react';

// ListStepper picks one of a fixed, ordered list of numbers (e.g. the
// official penalty amounts) with − and + buttons that move to the previous
// or next one, so game screens need no keyboard. It looks like ValueStepper:
// the number above, the buttons under it, to fit the center of the table on
// a 320px phone. At either end of the list its button is disabled.
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

  return (
    <div className="flex flex-col items-center gap-1">
      {/* The number is read out when it changes; the group names it. */}
      <output
        aria-label={label}
        aria-live="polite"
        className="flex h-9 w-24 items-center justify-center rounded-xl bg-field text-2xl font-semibold tabular-nums shadow-field"
      >
        {value}
      </output>
      <div className="flex items-center justify-center gap-1">
        <Button
          size="sm"
          variant="tertiary"
          className="min-w-14 md:h-10 md:min-w-16 md:text-base"
          aria-label={`${label} −`}
          isDisabled={index <= 0}
          onPress={() => step(-1)}
        >
          −
        </Button>
        <Button
          size="sm"
          variant="tertiary"
          className="min-w-14 md:h-10 md:min-w-16 md:text-base"
          aria-label={`${label} +`}
          isDisabled={index < 0 || index >= values.length - 1}
          onPress={() => step(1)}
        >
          +
        </Button>
      </div>
    </div>
  );
}
