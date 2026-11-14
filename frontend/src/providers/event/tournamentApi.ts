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
import type { SeatingMethod, TourState } from '../../models/Tour';

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
  seating: SeatingMethod,
  // step: an interval seating's step; undefined: the number of tables.
  step?: number
): Promise<boolean> {
  const result = await authFetch<unknown>(
    PostRequest(eventPath(eventId, '/tours'), { seating, step })
  );
  return result !== null;
}

// generateSeatingPlan prepares the seating for all the remaining tours,
// replacing an earlier plan. The search takes a few seconds.
export async function generateSeatingPlan(eventId: number): Promise<boolean> {
  const result = await authFetch<unknown>(
    PostRequest(eventPath(eventId, '/tours/plan'))
  );
  return result !== null;
}

// importSeatingPlan stores a seating for several tours from a CSV file's
// text, replacing an earlier plan.
export async function importSeatingPlan(
  eventId: number,
  csv: string
): Promise<boolean> {
  const result = await authFetch<unknown>(
    PutRequest(eventPath(eventId, '/tours/plan'), { csv })
  );
  return result !== null;
}

// deleteSeatingPlan deletes the seating prepared for all tours.
export async function deleteSeatingPlan(eventId: number): Promise<boolean> {
  const result = await authFetch<null>(
    DeleteRequest(eventPath(eventId, '/tours/plan'))
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
