# User scenarios

The agreed, desired behavior of MaCRa, area by area. When the code differs from this file, the code is wrong (or this file is out of date): fix one or the other, don't leave both.

Areas: 1. Roles and access · 2. Auth and account · 3. Events · 4. Playing a game · 5. Watching · 6. Ratings and history · 7. Landing and navigation. Tournament mode is a stub and will get its own area. Items marked "planned" are agreed but not built yet; the local known-issues and TODO lists track the gaps.

## 1. Roles and access

### Who is who in an event

- **Owner**: the user who created the event. Not a role, but the owner always has staff rights in their own event (see below), whatever roles they hold.
- **Member roles**: player, referee, observer. One user can hold several.
- **Member**: the owner or anyone with at least one role.
- **Staff**: referees, observers and the owner.
- When the owner creates a **regular** event, they also get the player and referee roles, so they can play and fix rounds right away. Tournament mode is a stub; its rules will be designed separately.

### Event settings that affect access

- **Visibility**: public (anyone, incl. anonymous visitors, can open the event) or private (members only). **Default: private.**
- **Ratings reach** and **games reach**: who sees the rating and the games — public, members or staff. "Public" means everyone only in a public event; in a private event it falls back to members. Defaults: public in a public event, members in a private one.

### What each user can see

- **Open the event** (overview): anyone in a public event; members in a private event. Without access to the event a user sees nothing in it, not even their own past games.
- **Member list**: members see the full list (names and roles). Everyone else sees only what the overview shows: the staff by name and the number of players.
- **Rating**: staff always; otherwise as set by ratings reach.
- **Other people's games** (live and finished lists, game pages, logs, players at the table): staff always; otherwise as set by games reach. In a tournament, "members" currently means staff only (stub).
- **Own games**: a player who sat at a game always sees that game (live, result, log, players) and their own game history, regardless of games reach, as long as they can open the event.
- **Game status** (live / ended / aborted, used to pick the right game page): anyone who can open the event.

### What each user can do

- **Create a game**: a player of a regular, unfinished event, who is one of the four players; all four must be players of the event.
- **Enter round results and penalties**: any of the four players seated in the game.
- **Undo the last round**: referees (and so the owner) only. Players can't, because a dishonest player could quietly undo a round and enter a better one. Revisit once undos leave a visible audit trail.
- **Abort a game** (regular events): referees at any time; the seated players only while no round has been entered (a game started by mistake).
- **Finish a game early** (regular events): a seated player or a referee.
- **Manage the event and its members, finish the event**: the owner. An event with live games can't be finished.
- **After the event is finished**: only its visibility can change. The name stays (players search by it), members stay (the rating is final), no new games.
- **Removing a member** who is in a live game is not allowed. A removed member of a private event loses all access to it.

## 2. Auth and account

### Sign up and email verification

- A visitor signs up with email, name and password. The backend validates all three, not only the form:
  - **email**: required, a plausible address, at most 254 characters;
  - **name**: required, at most 50 characters (same rule as renaming in the profile); names are not unique;
  - **password**: required, any length up to bcrypt's 72 bytes. Short passwords are allowed on purpose.
- Emails are **case-insensitive**: lookups and uniqueness use the lowercased address; mail goes to the address as typed.
- After sign-up a verification link is mailed (valid 24 hours). The user can ask to resend it. Neither response reveals whether the email is registered.
- Opening the link verifies the email and shows "Go to login". No automatic sign-in: the link may be opened on someone else's device.

### Sign in and sessions

- Sign in with email and password. A wrong email and a wrong password get the **same** answer ("email or password is incorrect"), so nobody can probe which emails are registered. An unverified email (with the right password) leads to the "check your email" screen.
- The session lives in an HttpOnly cookie only. The token is never returned in a response body nor stored in `localStorage`; the frontend keeps only the user's ID and name there.
- Sessions are **server-side**: each sign-in creates a session that lasts 400 days and slides while used, so an active user never gets signed out.
- **Sign out** ends this session on the server, not only in this browser.
- **Sign out everywhere** (profile) ends all of the user's sessions. A device list may come later.

### Forgotten password

- "Forgot password?" mails a reset link (valid 1 hour); the response doesn't reveal whether the account exists.
- Setting a new password through the link: ends all other sessions, marks the email as verified (the user proved they own the mailbox) and signs the user in right away.

### Profile

- Shows the user's name and their **user ID** (others use it to find them).
- Change the name; change the password (current + new; ends the other sessions).
- Later: change email; delete account. Deleting never deletes games or ratings: the account is anonymized to "Deleted player".

### Finding players

- Used when adding members and picking players. A numeric query matches a user ID exactly; otherwise names are searched (exact match first, then prefix, then substring), at least 2 characters, at most 20 results.
- Every player selector shows the user ID next to the name, so namesakes can be told apart. Email is never shown to other users.

### Abuse protection

- At most one auth email (verification or reset) per address per minute and about five per hour; beyond that the response still says "sent", without sending.
- nginx limits requests to `/auth/*` per IP (server config, kept outside the repo). No account lockout after failed sign-ins: it would let anyone lock someone else out.

## 3. Events

### Finding and choosing an event

- **Events list** (`/events`): two sections, "My events" (where I'm a member) and "Public events". In each, finished events are collapsed under "Show finished". Search and filters come later.
- **Select / switch event** (user menu): only events where I'm a member; active ones on top, finished ones in a collapsed section (still selectable, to see the last game and history). The selected event is where `/new`, `/now`, `/last` and `/history` work.

### Creating an event

- Any signed-in user can create one: name (required, trimmed, at most 100 characters), description (optional, at most 1000), type (tournament not available yet), visibility (private by default) and ratings/games reach with presets. The backend checks the same limits on create and edit.

### Event page

- **Overview**: name, type, finished mark, description, owner, number of players (players are not listed by name: there can be hundreds, and they may not want to be shown), names of referees and observers (staff are few and responsible for the event; a section with more than 5 names is collapsed).
- **≡ menu**: Manage, Rating, Live games, Finished games, as the user's rights allow; "Finish event" for the owner.

### Managing an event (owner)

- **Settings**: name, description, visibility, ratings/games reach. Changing visibility changes a reach only when it becomes invalid (public → members when the event turns private); a custom reach is kept otherwise.
- **Members**: add a user (server-side search by name or ID) with roles; tap a member to switch their roles; remove a member after a confirmation. The owner can't be removed and may drop their own player role. A member in a live game can't be removed.
- **Joining by invite**: the owner shares an invite link (or its QR code); a signed-in user who opens it becomes a player of the event. The owner can regenerate the link, which makes the old one stop working. (Planned.)
- **Leaving**: a member who is not the owner and not in a live game can leave a regular event themselves. Leaving a tournament is not possible. Past games and rating rows of a member who left or was removed stay as they are.
- **Deleting**: the owner can delete an event only while it has no games (created by mistake); otherwise it can only be finished.
- **Finishing**: the owner, after a confirmation, when no game is live. It can't be undone: the rating is final. After that only the visibility can change; members are read-only.

### Event settings (planned)

Two kinds, set by the owner in Manage → Rules, starting from a preset ("MCR standard": 16 hands, minimum 8; "Short": 8 hands, minimum 24) and adjusted value by value. The overview lists the settings that differ from the standard, so players know what they play.

- **Game rules** affect the score: number of hands (16 = East, South, West, North rounds; 8 = East and South), minimum hand value (8 by default), false mahjong amounts (30/60 by default). They can be changed only while the event has no games; after the first game they are frozen, so all games of one rating follow the same rules. Each game also stores a copy of the rules it started with and is always replayed by them.
- **Permissions** don't affect the score and can be changed at any time: who may finish a game early (seated players and referees / referees only), abort it (referees / also seated players before the first round / also seated players at any time), undo a round (referees / also seated players, once undos are audited), create games (players / referees only). Defaults are the ones in area 1.
- **Display and input** (can be changed at any time): minimum number of games to appear in the rating (the rating itself is computed the same way); hand input mode — total value (now) or picking fans, once fan-based input exists.
- Time limits per game are tournament-only and will be designed with the tournament mode. Not planned: a choice of who enters rounds, round confirmation by a second player, penalty switches, other rating methods or placement points, automatic seating.

## 4. Playing a game

Applies to regular events; tournament games will get their own rules.

### Starting a game (`/new`)

- A player of the selected event fills the four seats (East, South, West, North) from the event's players, IDs shown next to names; they must include themselves. Clear empties the seats, Shuffle mixes them.
- Players still in a live game **can** be picked: the next table is often prepared while two of its players finish their game. Play then fails with a clear message ("Player X is still in a game — wait until it ends") and the form keeps its seats.
- "Play again" on the result page opens the draft with the same players **in the same seats**; they shuffle themselves if they want.

### The table (`/now`)

- Shows the four seats in seating order, the scores and the round (東1…北4).
- Usually one player at the table keeps the score; any of the four may enter.
- **Entering a round**: tap the round kanji; pick the outcome (draw, self-draw, deal-in, penalty), 👍/👎 next to the seats, the hand value (stepper from 8, or type the number) or the fine (from 1; false mahjong 30 or 60; optionally move to the next round). While filling it in, each seat shows the score change the entry will make, so mistakes are seen before OK. After OK the scores update and the entry resets.
- **What a round holds**: exactly one main outcome (draw, self-draw or deal-in), which ends the round, and any number of penalties and false mahjongs, before it or instead of it (a penalty with "move to the next round" ends the round without a main outcome).
- **No silent duplicates**: an entry carries the number of entries in the game's round log as the player's table shows it (results and penalties alike, so a penalty tapped twice is caught too, though it doesn't change the round). If the log has changed since, the entry is refused, the table refreshes, and the message says what changed: "Round 東3 has already been entered", "A penalty has just been entered", or, after a referee's undo, "The last round has been undone" / "The last entry has been undone". The player decides whether to enter it again. Entering two forgotten rounds in a row still works.
- **Refreshing**: scores never change on their own, so a player's picture of the game isn't disturbed. The app checks every 10–15 seconds whether the table changed and then highlights the Refresh button; the player refreshes when ready.
- **Table menu**: round log, fan table, Finish (count the game now, with a confirmation; seated players and referees; not in tournaments), Abort (don't count it; referees, or the seated players before the first round), undo of the last round (referees).

### The end of a game

- After the 16th round the game ends by itself and the result page (`/last`) opens: final scores, chart, log, "Play again".
- If the game has ended while a player's table was still open (another player entered the last round, or it was finished), their next entry or Refresh shows "The game has ended" and opens the result page; if it was aborted, "The game was aborted" and the home page.
- A referee can undo the last round within an hour after the end; the game becomes live again. After that the result is final. The undo is refused if one of the players has already started another game ("Cannot undo the final round: <Name> has already started another game"): a player can't have two live games.
- **Undo is never applied twice**: like an entry, an undo carries the number of log entries the referee sees. If someone has entered or undone an entry since (two referees pressing Undo at once, a double tap), it is refused, the log reloads, and the message says what changed, in the same words as for an entry ("A round has just been entered", "A penalty has just been entered", "The last entry has already been undone").
- **Abandoned games**: a game left live (players walked away) blocks its players and the event's finish. Staff see a mark on live games with no new round for more than ~3 hours and decide themselves (abort or finish); nothing happens automatically.

## 5. Watching

### Game lists

- **Live games** and **Finished games** (event menu, by rights): each game shows the round progress (ring with the round kanji) and the four players with wind and score; finished games load in pages ("Load more").
- Live games has an **Auto-refresh** switch (off by default, remembered on the device). A game that ended stays in Live games for about 10 minutes, marked ✓ with the final scores, then moves to Finished games, so people around the table (or a club screen) can still see the result.

### Watching a game (`/event/N/game/M`)

- The page shows the live table or, for an ended game, the result (same as `/last`). Opening one's own game from the list shows it the same way; one's own table to play is reached from the home page.
- The live table can be rotated (who is at the bottom) and refreshed; Refresh keeps the rotation. Like at the table, scores change only on Refresh by default, and the Refresh button lights up when the game moved on; an **Auto-refresh** switch (remembered on the device) updates the scores by itself, e.g. for a club screen.
- The table menu (log, fans, Finish, Abort, Undo) follows the user's rights for that game: this is where referees manage other people's games.

### Stream overlay (HUD)

- `GET /api/event/:id/game/:gameId/hud` serves a compact game state for a stream overlay (OBS): round, winds, scores, differences, a message about the last round, the final result. It follows the same access rules as watching the game. It hasn't been updated for a while and must keep up with the rest (penalties, event settings such as 8 hands). The stream is run by a signed-in user, usually an observer (observers see all games of the event). The overlay is a page of the app for OBS's built-in Browser Source, signed in once inside OBS; if that is impractical, a revocable view-only link per event.

## 6. Ratings and history

### Event rating

- A table of players: place, name, games, TP, GP; sortable by any column. Only ended games count (not live or aborted ones). Members who left or were removed stay in it.
- Order: total TP, then total GP, then the best single games by GP. Players equal in all of that share the place (5, 5, then 7).
- The user's own row is highlighted; when it is far down, it stays pinned at the bottom of the screen.
- Players with fewer games than the event's minimum (event settings) are listed in a separate block "Fewer than N games" below the table.
- Tapping a row opens that player's history in the event (for users who can see the event's results; one's own always).

### History

- **History** (`/history`): the user's ended games in the selected event, newest first, 30 at a time with "Load more". On top, a summary for the event: number of games, total and average TP, how many times each place. Each game shows the place, TP, GP, date and the three opponents (exact look to be agreed on a mock-up first). One request returns a page of games with all of that.
- **Across events** (planned): the user picks which of their events to include (events with different rules shouldn't be mixed blindly). Only one's own games: other players' games in events the viewer isn't part of aren't available.

### Statistics (planned)

- **Player**: average / max / min TP and GP, hand values, number of wins (self-draw, deal-in), deal-ins, draws, points won and lost, and similar; for one event or a chosen set of the user's events.
- **Event**: the same figures for the whole event.

## 7. Landing and navigation

### Header

- Logo (home), "Events" (always visible), and the user menu: for signed-in users name, Profile, current event with Switch event, Settings, Log out; for visitors a generic person icon with Log in and Settings.

### Home page

- Visitors: "Login to play" (the game buttons are about playing; events to watch are in the header).
- Signed in without a selected event: "Select event to play".
- With a selected event: a small "Event" caption over the event name (a link to the event), then icon shortcuts with small captions (Rating, Live, Finished) as rights allow, then the game buttons: "Current game" when there is one, otherwise "New game" when the user may create games; "Last game" and "Played games" when there are past games.
- "Select event to play" opens the event picker (events where the user is a member). With no such events, the picker links to the events list.
- When there is nothing to do in the selected event, the place of the game buttons shows the reason: "You are not a player in this event", "The event is finished". Switching events stays in the user menu (no duplicate button). If the event is no longer accessible (403), the selection is cleared with a one-time toast ("No access to <event>") and "Select event to play" appears.
- No automatic event selection.

### Menus and overlays

- The events list has a round (+) button in the bottom-right corner to create an event (visitors go through sign-in first); no ≡ menu there.
- Event pages, the table (`/now`) and the live view have a ≡ menu button in the bottom-right corner; it must not cover page content (tentative fix: bottom padding on those pages, to be checked on screenshots).
- On Android, the system Back button or gesture closes an open menu, dialog or sheet instead of leaving the page.
- A crash in a page shows a short message with "Reload" and "Home" under a working header, never a blank white screen.

### Settings

- "Settings" in the user menu (for visitors too). Profile is about the account; Settings are about the app on this device.
- Language (English and Russian for a start; from the browser language on the first visit), theme (light / dark / system), keep the screen on at the table (on by default, only on the game screen), default position of the Auto-refresh switches, larger text.
- All settings are stored on the device. The language is also saved in the account of a signed-in user, so emails come in that language: changing it while signed in updates both; on a new device with no choice yet, the account's language is used.

### Screen sizes

- Phones come first; the desktop layout must keep working too.
- The narrowest supported width is 320 CSS px (the WCAG 2.1 reflow minimum; e.g. iPhone SE 1st gen). From 320 px up, every page works without horizontal scrolling and nothing is cut off.
- Narrower screens (e.g. 280 px on the cover screen of a folding phone) are best effort: sideways scrolling is acceptable there, and cheap fixes are welcome.
- UI changes are checked at 320 px, at a common phone width (360–390 px) and on desktop.
