import { useCallback, useEffect, useRef, useState } from 'react';
import { GetRequest, authFetch } from '../../components/Request';
import type { Player } from '../../models/Player';

// SEARCH_DELAY_MS: how long typing pauses before a search is sent.
const SEARCH_DELAY_MS = 250;

// isSearchable mirrors the backend's rules (docs/user-scenarios.md, "Finding
// players"): a user ID, or at least 2 characters of a name.
export function isSearchable(query: string): boolean {
  const q = query.trim();
  return /^\d+$/.test(q) || [...q].length >= 2;
}

// searchPlayers asks the server for the players matching query: a user ID
// exactly, or names containing it (at most 20, best matches first).
export async function searchPlayers(query: string): Promise<Player[]> {
  const request = GetRequest(`/api/players?q=${encodeURIComponent(query)}`);
  return (await authFetch<Player[]>(request)) ?? [];
}

// usePlayerSearch searches players as the user types: search(query) sends a
// request once typing pauses; players are the results of the latest query
// (an older answer that comes later is dropped).
export function usePlayerSearch() {
  const [players, setPlayers] = useState<Player[]>([]);
  const timer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const latest = useRef('');

  const search = useCallback((query: string) => {
    latest.current = query;
    clearTimeout(timer.current);
    timer.current = setTimeout(async () => {
      const found = isSearchable(query)
        ? await searchPlayers(query.trim())
        : [];
      if (latest.current === query) setPlayers(found);
    }, SEARCH_DELAY_MS);
  }, []);

  useEffect(() => () => clearTimeout(timer.current), []);

  return { players, search };
}
