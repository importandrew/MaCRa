import {
  GetRequest,
  PostRequest,
  authFetch,
  eventPath,
  fetchWithStatus,
} from '../../components/Request';
import { useFetch } from '../useFetch';
import type { TournamentParticipants } from '../../models/TournamentSettings';

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
