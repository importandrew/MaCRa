import { PlayerSelectorArgs } from '../../models/Player';
import PlayerAutocomplete from '../PlayerAutocomplete';

export default function PlayerSelector({
  wind,
  index,
  players,
  suggestions,
  selectedKeys,
  handleSuggestionsChange,
}: PlayerSelectorArgs) {
  const selectedKey = selectedKeys[index] || null;

  const handleSelectionChange = (playerId: number | null) => {
    if (playerId == null) {
      delete players[index];
    } else {
      players[index] = playerId;
    }
    handleSuggestionsChange();
  };

  return (
    <PlayerAutocomplete
      label={wind}
      items={suggestions}
      selectedKey={selectedKey}
      disabledKeys={selectedKeys.filter(Boolean)}
      onSelectionChange={handleSelectionChange}
    />
  );
}
