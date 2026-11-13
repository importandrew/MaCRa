import { PlayerSelectorArgs } from '../../models/Player';
import PlayerAutocomplete from '../PlayerAutocomplete';

export default function PlayerSelector({
  wind,
  index,
  suggestions,
  selectedKeys,
  onChange,
}: PlayerSelectorArgs) {
  const selectedKey = selectedKeys[index] || null;

  const handleSelectionChange = (playerId: number | null) =>
    onChange(index, playerId ?? undefined);

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
