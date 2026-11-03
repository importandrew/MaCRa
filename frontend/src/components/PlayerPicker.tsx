import { Label, ToggleButton, ToggleButtonGroup } from '@heroui/react';
import { useId } from 'react';
import { PlayerOrder } from '../models/GameData';

// PlayerPicker picks one of the table's players with a 2×2 grid of equal
// tiles: one tap, and long names wrap to a second line instead of pushing the
// layout apart.
export default function PlayerPicker({
  label,
  players,
  selectedId,
  onChange,
  disabledId,
}: {
  label: string;
  players: PlayerOrder[];
  selectedId: string;
  onChange: (id: string) => void;
  disabledId?: string;
}) {
  const labelId = useId();

  return (
    <div className="flex w-full flex-col gap-1">
      <Label id={labelId}>{label}</Label>
      <ToggleButtonGroup
        aria-labelledby={labelId}
        selectionMode="single"
        disallowEmptySelection
        isDetached
        className="grid w-full grid-cols-2 items-stretch gap-2"
        selectedKeys={selectedId ? [selectedId] : []}
        onSelectionChange={(keys) => {
          const [key] = keys;
          if (key != null) onChange(String(key));
        }}
      >
        {players.map((player) => (
          <ToggleButton
            key={player.id}
            id={String(player.id)}
            isDisabled={String(player.id) === disabledId}
            className="h-auto min-h-11 w-full min-w-0 rounded-xl px-2 py-1.5"
          >
            <span className="line-clamp-2 whitespace-normal break-words text-center leading-tight">
              {player.name}
            </span>
          </ToggleButton>
        ))}
      </ToggleButtonGroup>
    </div>
  );
}
