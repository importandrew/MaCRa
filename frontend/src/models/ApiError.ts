// ApiError is the body of an API error response (backend httperr.Error): a
// stable code and, for some codes, data the message needs.
export type ApiError = {
  error: string;
  // name is the player a message is about, e.g. for player_in_another_game.
  name?: string;
};
