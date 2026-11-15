import {
  DeleteRequest,
  GetRequest,
  PostRequest,
  PutRequest,
  authFetch,
  eventGamePath,
  eventPath,
  fetchWithStatus,
} from '../../components/Request';
import { useFetch } from '../useFetch';
import type { TournamentParticipants } from '../../models/TournamentSettings';
import type { Seating, TourState } from '../../models/Tour';

// A tournament's players with their numbers, for its desk (staff).
export const useFetchParticipants = (eventId: number | null) => {
  const f = useFetch<TournamentParticipants>(eventId, () =>
    fetchWithStatus<TournamentParticipants>(
      GetRequest(eventPath(eventId!, '/participants'))
    )
  );
  return {
    data: f.data,
    isLoading: f.isLoading,
    isError: f.isFailed,
    isForbidden: f.isForbidden,
    refetch: f.refetch,
  };
};

// drawNumbers gives the players new random numbers (a referee, before the
// first tour starts).
export async function drawNumbers(eventId: number): Promise<boolean> {
  const result = await authFetch<null>(
    PostRequest(eventPath(eventId, '/draw'))
  );
  return result !== null;
}

// withdrawParticipant marks a player of a started tournament as withdrawn.
export async function withdrawParticipant(
  eventId: number,
  userId: number
): Promise<boolean> {
  const result = await authFetch<null>(
    PostRequest(eventPath(eventId, `/participants/${userId}/withdraw`))
  );
  return result !== null;
}

// The tournament's current tour with its tables, and what the next would be.
export const useFetchTours = (eventId: number | null) => {
  const f = useFetch<TourState>(eventId, () =>
    fetchWithStatus<TourState>(GetRequest(eventPath(eventId!, '/tours')))
  );
  return {
    data: f.data,
    isLoading: f.isLoading,
    isError: f.isFailed,
    refetch: f.refetch,
  };
};

// generateTour drafts the next tour's seating (seen by the staff only),
// replacing an earlier draft.
export async function generateTour(
  eventId: number,
  how: Seating
): Promise<boolean> {
  const result = await authFetch<unknown>(
    PostRequest(eventPath(eventId, '/tours'), how)
  );
  return result !== null;
}

// setTourPlan sets the seating a tour is to have.
export async function setTourPlan(
  eventId: number,
  tour: number,
  how: Seating
): Promise<boolean> {
  const result = await authFetch<null>(
    PutRequest(eventPath(eventId, `/tours/plan/${tour}`), how)
  );
  return result !== null;
}

// importSeatingFile stores tables for several tours from a CSV file's text,
// replacing the ones imported before; their tours are planned fixed.
export async function importSeatingFile(
  eventId: number,
  csv: string
): Promise<boolean> {
  const result = await authFetch<unknown>(
    PutRequest(eventPath(eventId, '/tours/file'), { csv })
  );
  return result !== null;
}

// deleteSeatingFile deletes the tables imported from a file.
export async function deleteSeatingFile(eventId: number): Promise<boolean> {
  const result = await authFetch<null>(
    DeleteRequest(eventPath(eventId, '/tours/file'))
  );
  return result !== null;
}

// cancelTour cancels the published seating of the tour not started yet.
export async function cancelTour(eventId: number): Promise<boolean> {
  const result = await authFetch<null>(
    DeleteRequest(eventPath(eventId, '/tours/current'))
  );
  return result !== null;
}

// tourAction publishes the drafted seating, starts the published tour, or
// finishes or aborts the live one.
export async function tourAction(
  eventId: number,
  action: 'publish' | 'start' | 'finish' | 'abort'
): Promise<boolean> {
  const result = await authFetch<null>(
    PostRequest(eventPath(eventId, `/tours/current/${action}`))
  );
  return result !== null;
}

// seatSitter seats a substitute at a seat of a tournament game (forfeit: the
// seat's player gets 0 / 0), or the seat's player back (sitter = seat), from
// the next log entry on.
export async function seatSitter(
  eventId: number,
  gameId: number,
  change: { seat: number; sitter: number; forfeit: boolean; log_size: number }
): Promise<boolean> {
  const result = await authFetch<null>(
    PostRequest(eventGamePath(eventId, gameId, '/seat'), change)
  );
  return result !== null;
}

// penalize enters a referee's penalty against a player of a live game (a
// late player back at the table).
export async function penalize(
  eventId: number,
  gameId: number,
  playerId: number,
  value: number,
  logSize: number
): Promise<boolean> {
  const result = await authFetch<unknown>(
    PostRequest(eventGamePath(eventId, gameId, '/penalty'), {
      loser: playerId,
      value,
      force_next_round: false,
      log_size: logSize,
    })
  );
  return result !== null;
}

// addTime adds minutes to a table of the live tour.
export async function addTime(
  eventId: number,
  gameId: number,
  minutes: number
): Promise<boolean> {
  const result = await authFetch<null>(
    PostRequest(eventGamePath(eventId, gameId, '/time'), { minutes })
  );
  return result !== null;
}
