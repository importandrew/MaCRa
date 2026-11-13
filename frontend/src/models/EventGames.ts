export type EventGameListStatus = 'live' | 'ended';

export type EventGameListPlayer = {
  id: number;
  name: string;
  wind?: string;
  gp: number;
};

export type EventGameListItem = {
  id: number;
  status: EventGameListStatus;
  round: { num: number; kanji: string } | null;
  // hands: the number of hands of a live game, for its progress.
  hands?: number;
  players: EventGameListPlayer[];
};

export type EventGamesListData = {
  games: EventGameListItem[];
  has_more: boolean;
};
