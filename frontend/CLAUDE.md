# Frontend

- Install with `npm ci`, which installs from the lockfile without modifying it.
- API access lives in `src/providers/*/…Api.ts` as hooks and functions built on `components/Request.ts` (`authFetch`, `publicFetch`, `eventPath`, …). That module handles errors globally: it shows a toast on non-OK responses (except 403), and on 401 it logs out and redirects to `/login`. Don't duplicate this handling in callers.
- The app is mobile first: about 90% of users are on phones, but the desktop layout must keep working. Check UI changes at phone width first. On phones, side menus and modals take the full screen width, so there's no backdrop to tap to close them.
- UI is HeroUI v3 (React Aria Components, compound components like `Card.Content`, styles from `@heroui/styles`; no provider, no Tailwind plugin, no framer-motion). Shared building blocks in `src/components/`:
  - `ContextMenu/ContextMenu` — the ≡ button plus its side menu; use it for every ≡ menu. An open v3 Drawer makes everything outside it inert, so `SideMenu` draws its own ≡ close button at the same screen spot.
  - `LinkButton` — a react-router link styled as a button (v3 Button has no `as`). HeroUI `href` props navigate through react-router via `RouterProvider` in `App.tsx`.
  - `buttonTint()` — soft colored fill for a button (`variant="secondary"` + `className={buttonTint('accent')}`), like v2's `variant="flat"` with a color; v3 has no such variants.
  - `Confirm` (AlertDialog), `SegmentedPicker` (ToggleButtonGroup for one-of-few choices), `LazyAccordion` (v3 Accordion keeps every panel mounted; this one renders a panel only while open), `NavList` (list of links), `UserAvatar`, `LoadingTemplate/LoadingSpinner`.
- Colors use v3 tokens (`text-muted`, `bg-default`, `bg-accent-soft`, `border-border`, …); v2 numbered scales like `text-default-500` no longer exist. Dark mode is class-only (`index.css`), and the app never sets it.
- Number inputs that gate a submit button use a `TextField` with `inputMode="numeric"`, not `NumberField`: NumberField commits its value only on blur.
