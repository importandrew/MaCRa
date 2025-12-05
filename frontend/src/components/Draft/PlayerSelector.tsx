import { Key } from 'react';
import { Autocomplete, AutocompleteItem } from '@heroui/react';
import { PlayerSelectorArgs } from '../../models/Player';
import PlayerId from './PlayerId';

export default function PlayerSelector({
  wind,
  index,
  players,
  suggestions,
  selectedKeys,
  handleSuggestionsChange,
}: PlayerSelectorArgs) {
  const handleSuggestionClick = (playerId: Key | null, index: number) => {
    if (!playerId) return;

    players[index] = Number(playerId);
    handleSuggestionsChange();
  };

  const onClear = () => {
    delete players[index];
    handleSuggestionsChange();
  };

  return (
    <Autocomplete
      label={wind}
      itemHeight={50}
      selectedKey={selectedKeys[index]}
      isClearable={true}
      disabledKeys={selectedKeys}
      className="w-full md:w-[300px]"
      onSelectionChange={(id) => handleSuggestionClick(id, index)}
      clearButtonProps={{ onPress: onClear }}
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
