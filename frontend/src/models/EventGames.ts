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
  // hidden_tour: the live tour whose tables the caller can't see yet (they
  // are seen by the staff and the players whose own game ended); missing
  // when none is hidden.
  hidden_tour?: number;
};
