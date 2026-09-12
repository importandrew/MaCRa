import { Key } from 'react';
import { Autocomplete, AutocompleteItem } from '@heroui/react';
import { Player } from '../models/Player';

function PlayerId({ id }: { id: number }) {
  return <span className="text-xs">[{id}]</span>;
}

export type PlayerAutocompleteProps = {
  label: string;
  items: Player[];
  selectedKey: string | null;
  onSelectionChange: (id: number | null) => void;
  disabledKeys?: string[];
};

export default function PlayerAutocomplete({
  label,
  items,
  selectedKey,
  onSelectionChange,
  disabledKeys = [],
}: PlayerAutocompleteProps) {
  const handleSelectionChange = (key: Key | null) => {
    if (key == null) {
      onSelectionChange(null);
      return;
    }
    onSelectionChange(Number(key));
  };

  return (
    <Autocomplete
      label={label}
      itemHeight={50}
      selectedKey={selectedKey}
      isClearable={true}
      disabledKeys={disabledKeys}
      className="w-full md:w-[300px]"
      onSelectionChange={handleSelectionChange}
    >
      {items.map((player) => (
        <AutocompleteItem
          key={String(player.id)}
          textValue={player.name}
          endContent={<PlayerId id={player.id} />}
        >
          {player.name}
        </AutocompleteItem>
      ))}
    </Autocomplete>
  );
}
