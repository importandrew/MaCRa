# Frontend

- Install with `npm ci`, which installs from the lockfile without modifying it.
- API access lives in `src/providers/*/…Api.ts` as hooks and functions built on `components/Request.tsx` (`authFetch`, `publicFetch`, `eventPath`, …). That module handles errors globally: it shows a toast on non-OK responses (except 403), and on 401 it logs out and redirects to `/login`. Don't duplicate this handling in callers.
