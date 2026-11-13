export interface Player {
  id: number;
  name: string;
  score?: number;
}

export interface PlayerSelectorArgs {
  wind: string;
  index: number;
  suggestions: Player[];
  selectedKeys: string[];
  // onChange picks (or, with undefined, clears) the player of seat index.
  onChange: (index: number, playerId: number | undefined) => void;
}
