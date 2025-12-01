export interface Player {
  id: number;
  name: string;
  score?: number;
}

export interface PlayerSelectorArgs {
  wind: string;
  index: number;
  inputValue: string;
  suggestions: Player[];
  selectedKey: string;
  selectedKeys: string[];
  isFocused: boolean;
  onInputChange: (value: string) => void;
  onSelectionChange: (playerId: number | null) => void;
  onClear: () => void;
  onFocus: () => void;
  onBlur: () => void;
}