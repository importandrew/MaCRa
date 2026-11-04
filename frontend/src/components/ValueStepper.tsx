import { Button, Input, TextField } from '@heroui/react';
import { useEffect, useRef, useState } from 'react';

const LONG_PRESS_MS = 500;
const STEPS = [-10, -1, 1, 10];

// ValueStepper sets a number with −10/−1/+1/+10 buttons, so game screens
// don't need the phone keyboard. A long press on the number (or Enter/Space
// on it) opens a numeric field for values far from the start, e.g. 88.
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
  const [isEditing, setIsEditing] = useState(false);
  const [draft, setDraft] = useState('');
  const timer = useRef<number | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const cancelLongPress = () => {
    if (timer.current !== null) {
      window.clearTimeout(timer.current);
      timer.current = null;
    }
  };

  useEffect(() => cancelLongPress, []);

  // Focus the field when it opens, so the numeric keyboard comes up.
  useEffect(() => {
    if (isEditing) {
      inputRef.current?.focus();
      inputRef.current?.select();
    }
  }, [isEditing]);

  const startEditing = () => {
    cancelLongPress();
    setDraft(String(value));
    setIsEditing(true);
  };

  const commit = () => {
    const typed = Number(draft);
    if (draft !== '' && Number.isFinite(typed)) {
      onChange(Math.max(min, typed));
    }
    setIsEditing(false);
  };

  // The number sits above the buttons, not between them: the stepper must
  // fit the center of the table on a 360px phone.
  return (
    <div className="flex flex-col items-center gap-1">
      <span className="text-xs text-muted">{label}</span>
      {isEditing ? (
        <TextField
          aria-label={label}
          className="w-24"
          value={draft}
          onChange={(text) => setDraft(text.replace(/\D/g, ''))}
          onBlur={commit}
          onKeyDown={(e) => {
            if (e.key === 'Enter') commit();
          }}
        >
          <Input
            ref={inputRef}
            inputMode="numeric"
            pattern="\d*"
            autoComplete="off"
            className="h-10 text-center text-2xl font-semibold"
          />
        </TextField>
      ) : (
        // A plain button, not HeroUI's: React Aria's press handling has no
        // long press, and its pointer capture would swallow the timer.
        <button
          type="button"
          aria-label={`${label} ${value}. Long press to type a value`}
          className="h-10 w-24 select-none touch-manipulation text-center text-foreground text-4xl leading-none font-semibold tabular-nums [-webkit-touch-callout:none]"
          onPointerDown={() => {
            cancelLongPress();
            timer.current = window.setTimeout(startEditing, LONG_PRESS_MS);
          }}
          onPointerUp={cancelLongPress}
          onPointerLeave={cancelLongPress}
          onPointerCancel={cancelLongPress}
          onContextMenu={(e) => e.preventDefault()}
          onKeyDown={(e) => {
            if (e.key === 'Enter' || e.key === ' ') {
              e.preventDefault();
              startEditing();
            }
          }}
        >
          {value}
        </button>
      )}
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
            {step > 0 ? `+${step}` : `\u2212${-step}`}
          </Button>
        ))}
      </div>
    </div>
  );
}
