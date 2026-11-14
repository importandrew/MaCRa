import { useTranslation } from 'react-i18next';
import {
  CloseButton,
  ComboBox,
  Description,
  Input,
  Label,
  ListBox,
  type Key,
} from '@heroui/react';
import { useRef, useState } from 'react';
import { Player } from '../models/Player';

function PlayerId({ id }: { id: number }) {
  return <span className="text-xs text-muted">[{id}]</span>;
}

// matchesPlayer: a number is a user ID and matches exactly; anything else is
// part of a name, ignoring case (as the server's search does).
function matchesPlayer(player: Player, query: string): boolean {
  const q = query.trim();
  if (/^\d+$/.test(q)) return player.id === Number(q);
  return player.name.toLocaleLowerCase().includes(q.toLocaleLowerCase());
}

export type PlayerAutocompleteProps = {
  label: string;
  items: Player[];
  selectedKey: string | null;
  onSelectionChange: (id: number | null) => void;
  disabledKeys?: string[];
  // onSearch: the items come from a server search (Manage → add a member),
  // called with the text as the user types; without it the field filters
  // items itself (picking players from the event's members).
  onSearch?: (query: string) => void;
  description?: string;
};

// PlayerAutocomplete picks a player by typing part of the name or the user
// ID: tap the field and type. Each item shows the ID, so namesakes can be told
// apart.
export default function PlayerAutocomplete({
  label,
  items,
  selectedKey,
  onSelectionChange,
  disabledKeys = [],
  onSearch,
  description,
}: PlayerAutocompleteProps) {
  const { t } = useTranslation();
  const inputRef = useRef<HTMLInputElement>(null);
  const [inputValue, setInputValue] = useState('');
  // picked: the player last picked here, kept for its name: a server search
  // may no longer list them.
  const [picked, setPicked] = useState<Player | null>(null);
  const selected =
    items.find((player) => String(player.id) === selectedKey) ??
    (picked && String(picked.id) === selectedKey ? picked : null);

  // The text follows the selection when it changes from outside (Clear,
  // Shuffle, a reset after adding).
  const [shownKey, setShownKey] = useState<string | null>(selectedKey);
  if (selectedKey !== shownKey) {
    setShownKey(selectedKey);
    setInputValue(selected?.name ?? '');
  }

  // Items are filtered here (controlled items get no filtering from the
  // ComboBox); the selected player's own name lists everyone.
  const shown =
    onSearch || (selected && inputValue === selected.name) || !inputValue.trim()
      ? items
      : items.filter((player) => matchesPlayer(player, inputValue));

  const handleSelectionChange = (key: Key | null) => {
    const player = items.find((p) => String(p.id) === String(key)) ?? null;
    setPicked(player);
    setShownKey(key == null ? null : String(key));
    setInputValue(player?.name ?? '');
    onSelectionChange(key == null ? null : Number(key));
    // React Aria keeps focus in the input after a pick (it returns focus there
    // when the list closes), which leaves the phone keyboard open. Drop the
    // focus on the next frame, after the list has closed.
    if (key != null) requestAnimationFrame(() => inputRef.current?.blur());
  };

  const handleInputChange = (text: string) => {
    setInputValue(text);
    onSearch?.(text);
  };

  return (
    <ComboBox
      items={shown}
      inputValue={inputValue}
      onInputChange={handleInputChange}
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
        <Input
          ref={inputRef}
          className={selectedKey != null ? 'pe-14' : undefined}
        />
        {selectedKey != null && (
          <CloseButton
            aria-label={t('common.clearField', { field: label })}
            onPress={() => handleSelectionChange(null)}
            className="absolute end-7 top-1/2 z-10 -translate-y-1/2"
          />
        )}
        <ComboBox.Trigger />
      </ComboBox.InputGroup>
      {description && <Description>{description}</Description>}
      <ComboBox.Popover>
        <ListBox>
          {(player: Player) => (
            <ListBox.Item id={String(player.id)} textValue={player.name}>
              <span className="flex w-full items-center justify-between gap-2">
                {player.name}
                <PlayerId id={player.id} />
              </span>
            </ListBox.Item>
          )}
        </ListBox>
      </ComboBox.Popover>
    </ComboBox>
  );
}
