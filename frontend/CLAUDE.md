# Frontend

- Install with `npm ci`, which installs from the lockfile without modifying it.
- API access lives in `src/providers/*/…Api.ts` as hooks and functions built on `components/Request.ts` (`authFetch`, `publicFetch`, `eventPath`, …). That module handles errors globally: it shows a toast on non-OK responses (except 403), and on 401 it logs out and redirects to `/login`. Don't duplicate this handling in callers.
- The app is mobile first: about 90% of users are on phones, but the desktop layout must keep working. Check UI changes at phone width first. On phones, side menus and modals take the full screen width, so there's no backdrop to tap to close them.
