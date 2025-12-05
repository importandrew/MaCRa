export interface Player {
  id: number;
  name: string;
  score?: number;
}

export interface PlayerSelectorArgs {
  wind: string;
  index: number;
  players: (number | undefined)[];
  suggestions: Player[];
  selectedKeys: string[];
  handleSuggestionsChange: () => void;
}
