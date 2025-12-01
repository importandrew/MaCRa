import { Key } from 'react';
import { Autocomplete, AutocompleteItem } from '@heroui/react';
import { PlayerSelectorArgs } from '../../models/Player';
import PlayerId from './PlayerId';

export default function PlayerSelector({
  wind,
  inputValue,
  suggestions,
  selectedKey,
  selectedKeys,
  onInputChange,
  onSelectionChange,
  onClear,
  onFocus,
  onBlur,
}: PlayerSelectorArgs) {
  const handleSelectionChange = (playerId: Key | null) => {
    onSelectionChange(playerId ? Number(playerId) : null);
  };

  return (
    <Autocomplete
      label={wind}
      itemHeight={50}
      inputValue={inputValue}
      selectedKey={selectedKey}
      isClearable={true}
      disabledKeys={selectedKeys}
      className="w-full md:w-[300px]"
      onInputChange={(value) => onInputChange(value)}
      onSelectionChange={handleSelectionChange}
      clearButtonProps={{ onPress: onClear }}
      onFocus={onFocus}
      onBlur={onBlur}
      allowsCustomValue={true}
    >
      {suggestions?.map((player) => (
        <AutocompleteItem
          key={player.id}
          endContent={<PlayerId id={player.id} />}
        >
          {player.name}
        </AutocompleteItem>
      ))}
    </Autocomplete>
  );
}