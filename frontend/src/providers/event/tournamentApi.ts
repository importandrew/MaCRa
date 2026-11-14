import {
  DeleteRequest,
  GetRequest,
  PostRequest,
  authFetch,
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

// publishTour publishes the next tour's seating: its games are created and
// each player sees their table.
export async function publishTour(
  eventId: number,
  seating: SeatingMethod
): Promise<boolean> {
  const result = await authFetch<unknown>(
    PostRequest(eventPath(eventId, '/tours'), { seating })
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
