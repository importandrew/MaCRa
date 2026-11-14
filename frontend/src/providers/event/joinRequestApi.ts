import {
  DeleteRequest,
  GetRequest,
  PostRequest,
  PutRequest,
  authFetch,
  eventPath,
  fetchWithStatus,
} from '../../components/Request';
import { useFetch } from '../useFetch';
import type { JoinRequest, JoiningSettings } from '../../models/JoinRequest';

// updateJoining sets whether the event takes requests to join and its
// player limit (Manage → Members).
export async function updateJoining(
  eventId: number,
  settings: JoiningSettings
): Promise<boolean> {
  const result = await authFetch<null>(
    PutRequest(eventPath(eventId, '/joining'), settings)
  );
  return result !== null;
}

// The event's pending requests to join, oldest first, for its owner.
export const useFetchJoinRequests = (eventId: number | null) => {
  const f = useFetch<JoinRequest[]>(eventId, () =>
    fetchWithStatus<JoinRequest[]>(GetRequest(eventPath(eventId!, '/requests')))
  );
  return {
    requests: f.data,
    isLoading: f.isLoading,
    isError: f.isFailed,
    refetch: f.refetch,
  };
};

// answerJoinRequest accepts (the user becomes a player) or declines (they
// can't ask again) a user's request.
export async function answerJoinRequest(
  eventId: number,
  userId: number,
  answer: 'accept' | 'decline'
): Promise<boolean> {
  const result = await authFetch<null>(
    PostRequest(eventPath(eventId, `/requests/${userId}/${answer}`))
  );
  return result !== null;
}

// requestToJoin sends the signed-in user's request to join the event.
export async function requestToJoin(eventId: number): Promise<boolean> {
  const result = await authFetch<null>(
    PostRequest(eventPath(eventId, '/request'))
  );
  return result !== null;
}

// withdrawJoinRequest deletes the signed-in user's pending request.
export async function withdrawJoinRequest(eventId: number): Promise<boolean> {
  const result = await authFetch<null>(
    DeleteRequest(eventPath(eventId, '/request'))
  );
  return result !== null;
}
