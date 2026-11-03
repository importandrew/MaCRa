import {
  CloseButton,
  ComboBox,
  Input,
  Label,
  ListBox,
  type Key,
} from '@heroui/react';
import { Player } from '../models/Player';

function PlayerId({ id }: { id: number }) {
  return <span className="text-xs text-muted">[{id}]</span>;
}

export type PlayerAutocompleteProps = {
  label: string;
  items: Player[];
  selectedKey: string | null;
  onSelectionChange: (id: number | null) => void;
  disabledKeys?: string[];
};

// PlayerAutocomplete picks a player by typing part of the name: tap the field
// and type, as with HeroUI v2 Autocomplete.
export default function PlayerAutocomplete({
  label,
  items,
  selectedKey,
  onSelectionChange,
  disabledKeys = [],
}: PlayerAutocompleteProps) {
  const handleSelectionChange = (key: Key | null) => {
    onSelectionChange(key == null ? null : Number(key));
  };

  return (
    <ComboBox
      selectedKey={selectedKey}
      disabledKeys={disabledKeys}
      className="w-full md:w-[300px]"
      onSelectionChange={handleSelectionChange}
    >
      <Label>{label}</Label>
      <ComboBox.InputGroup>
        {/* The trigger arrow is absolutely positioned at the end of the field,
            so the clear button sits just before it and the input reserves
            room for both. */}
        <Input className={selectedKey != null ? 'pe-14' : undefined} />
        {selectedKey != null && (
          <CloseButton
            aria-label={`Clear ${label}`}
            onPress={() => handleSelectionChange(null)}
            className="absolute end-7 top-1/2 z-10 -translate-y-1/2"
          />
        )}
        <ComboBox.Trigger />
      </ComboBox.InputGroup>
      <ComboBox.Popover>
        <ListBox>
          {items.map((player) => (
            <ListBox.Item
              key={player.id}
              id={String(player.id)}
              textValue={player.name}
            >
              <span className="flex w-full items-center justify-between gap-2">
                {player.name}
                <PlayerId id={player.id} />
              </span>
            </ListBox.Item>
          ))}
        </ListBox>
      </ComboBox.Popover>
    </ComboBox>
  );
}
