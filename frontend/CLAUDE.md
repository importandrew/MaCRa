# Frontend

- Install with `npm ci`, which installs from the lockfile without modifying it.
- API access lives in `src/providers/*/…Api.ts` as hooks and functions built on `components/Request.ts` (`authFetch`, `publicFetch`, `eventPath`, …). That module handles errors globally: it shows a toast on non-OK responses (except 403), and on 401 it logs out and redirects to `/login`. Don't duplicate this handling in callers.
- The app is mobile first: about 90% of users are on phones, but the desktop layout must keep working. Check UI changes at phone width first. On phones, side menus and modals take the full screen width, so there's no backdrop to tap to close them.
- UI is HeroUI v3 (React Aria Components, compound components like `Card.Content`, styles from `@heroui/styles`; no provider, no Tailwind plugin, no framer-motion). Shared building blocks in `src/components/`:
  - `ContextMenu/ContextMenu` — the ≡ button plus its side menu; use it for every ≡ menu. The menu is full width on phones (no backdrop to tap) and closes with its ✕ button or by swiping the header to the right. While it's open, React Aria makes everything outside it inert, so the ≡ button can't close it.
  - `PendingButton` — a Button that shows a spinner before its label while `isPending` (v3 leaves the pending look to the caller).
  - `LinkButton` — a react-router link styled as a button (v3 Button has no `as`). HeroUI `href` props navigate through react-router via `RouterProvider` in `App.tsx`.
  - `Confirm` (AlertDialog), `SegmentedPicker` (ToggleButtonGroup for one-of-few choices), `ValueStepper` (a numeric field with −10/−1/+1/+10 buttons under it, so game screens rarely need the keyboard; shows an error below its minimum), `LazyAccordion` (v3 Accordion keeps every panel mounted; this one renders a panel only while open), `NavList` (list of links), `UserAvatar`, `LoadingTemplate/LoadingSpinner`.
- Buttons use only the built-in v3 variants, by importance on the screen: `primary` (filled) for the main action, `secondary` for the next one, `tertiary` for the rest, `danger`/`danger-soft` for destructive ones. No custom color fills.
- Colors use v3 tokens (`text-muted`, `bg-default`, `bg-accent-soft`, `border-border`, …); v2 numbered scales like `text-default-500` no longer exist. Dark mode is class-only (`index.css`), and the app never sets it.
- Number inputs that gate a submit button use a `TextField` with `inputMode="numeric"`, not `NumberField`: NumberField commits its value only on blur.
