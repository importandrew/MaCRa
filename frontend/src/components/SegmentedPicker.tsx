import { ToggleButton, ToggleButtonGroup } from '@heroui/react';
import { useId } from 'react';

export type SegmentedPickerOption<T extends string> = {
  value: T;
  label: string;
  isDisabled?: boolean;
};

// SegmentedPicker picks one of a few options with a row of toggle buttons
// (event type, visibility, ...). One tap, and the choice is always visible.
export default function SegmentedPicker<T extends string>({
  label,
  ariaLabel,
  options,
  selectedKey,
  onSelectionChange,
  isOptionDisabled,
}: {
  label?: string;
  ariaLabel?: string;
  options: SegmentedPickerOption<T>[];
  // selectedKey is null until the user picks an option.
  selectedKey: T | null;
  onSelectionChange: (key: T) => void;
  isOptionDisabled?: (value: T) => boolean;
}) {
  const labelId = useId();

  return (
    <div className="flex flex-col gap-1.5">
      {label && (
        <p id={labelId} className="text-sm">
          {label}
        </p>
      )}
      <ToggleButtonGroup
        aria-labelledby={label ? labelId : undefined}
        aria-label={label ? undefined : ariaLabel}
        selectionMode="single"
        disallowEmptySelection
        fullWidth
        className="items-stretch"
        selectedKeys={selectedKey ? [selectedKey] : []}
        onSelectionChange={(keys) => {
          const [key] = keys;
          if (key != null) onSelectionChange(String(key) as T);
        }}
      >
        {options.map((opt, index) => (
          <ToggleButton
            key={opt.value}
            id={opt.value}
            isDisabled={opt.isDisabled || isOptionDisabled?.(opt.value)}
            className="min-w-0 h-auto min-h-9 py-1.5 not-data-selected:text-foreground"
          >
            {index > 0 && <ToggleButtonGroup.Separator />}
            <span className="whitespace-normal text-center leading-tight">
              {opt.label}
            </span>
          </ToggleButton>
        ))}
      </ToggleButtonGroup>
    </div>
  );
}
