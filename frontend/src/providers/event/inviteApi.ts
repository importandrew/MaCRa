import {
  DeleteRequest,
  GetRequest,
  PostRequest,
  authFetch,
  eventPath,
  fetchWithStatus,
  publicFetch,
} from '../../components/Request';
import { useFetch } from '../useFetch';
import type { EventInvite } from '../../models/EventInvite';

// The invite link of an event, for its owner (Manage → Members): token is
// null while the link is off.
export const useFetchInviteToken = (eventId: number | null) => {
  const f = useFetch<{ token: string | null }>(eventId, () =>
    fetchWithStatus<{ token: string | null }>(
      GetRequest(eventPath(eventId!, '/invite'))
    )
  );
  return {
    token: f.data?.token ?? null,
    isLoading: f.isLoading,
    isError: f.isFailed,
    mutate: f.mutate,
  };
};

// newInviteToken turns the link on with a new token (the old link stops
// working) and returns it.
export async function newInviteToken(eventId: number): Promise<string | null> {
  const result = await authFetch<{ token: string }>(
    PostRequest(eventPath(eventId, '/invite'))
  );
  return result?.token ?? null;
}

export async function turnOffInvite(eventId: number): Promise<boolean> {
  const result = await authFetch<null>(
    DeleteRequest(eventPath(eventId, '/invite'))
  );
  return result !== null;
}

// inviteUrl is the link to share: /join/<token> on this site.
export function inviteUrl(token: string): string {
  return `${window.location.origin}/join/${token}`;
}

// The event an invite link leads to, for whoever opens it (a visitor too).
// A 404 is quiet: the page says the link doesn't work.
export const useFetchInvite = (token: string | null) => {
  const f = useFetch<EventInvite>(token, () =>
    publicFetch<EventInvite>(GetRequest(`/api/invite/${token}`), {
      quietStatuses: [404],
    })
  );
  return {
    invite: f.data,
    isLoading: f.isLoading,
    isNotFound: f.isNotFound,
    isError: f.isError,
    refetch: f.refetch,
  };
};

// joinByInvite makes the signed-in user a player of the link's event and
// returns its ID.
export async function joinByInvite(token: string): Promise<number | null> {
  const result = await authFetch<{ event_id: number }>(
    PostRequest(`/api/invite/${token}/join`)
  );
  return result?.event_id ?? null;
}
