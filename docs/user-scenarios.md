# User scenarios

The agreed, desired behavior of MaCRa, area by area. When the code differs from this file, the code is wrong (or this file is out of date): fix one or the other, don't leave both.

Areas: 1. Roles and access · 2. Auth and account · 3. Events · 4. Playing a game · 5. Watching · 6. Ratings and history · 7. Landing and navigation · 8. Tournaments (planned). Items marked "planned" are agreed but not built yet; the local known-issues and TODO lists track the gaps.

## 1. Roles and access

### Who is who in an event

- **Owner**: the user who created the event. Not a role. The owner always sees everything in their own event, as staff do, and manages it (settings, members, finishing), whatever roles they hold. Acting as a referee (undo, abort) takes the referee role, also for the owner: an organizer who isn't a referee can leave it out and can't press those by accident. In a regular event the owner gets the referee role at creation.
- **Member roles**: player, referee, observer; in a tournament also substitute, never together with player (area 8). One user can hold several.
- **Member**: the owner or anyone with at least one role.
- **Staff**: referees, observers and the owner.
- When the owner creates a **regular** event, they also get the player and referee roles, so they can play and fix rounds right away. In a tournament the owner gets no roles by default (area 8).

### Event settings that affect access

- **Visibility**: public (anyone, incl. anonymous visitors, can open the event) or private (members only). **Default: private.**
- **Ratings reach** and **games reach**: who sees the rating and the games — public, members or staff. "Public" means everyone only in a public event; in a private event it falls back to members. Defaults: public in a public event, members in a private one.

### What each user can see

- **Open the event** (overview): anyone in a public event; members in a private event. Without access to the event a user sees nothing in it, not even their own past games.
- **Member list**: members see the full list (names and roles). Everyone else sees only what the overview shows: the staff by name and the number of players.
- **Rating**: staff always; otherwise as set by ratings reach.
- **Other people's games** (live and finished lists, game pages, logs, players at the table): staff always; otherwise as set by games reach. In a tournament, live tables are also hidden during a tour (area 8).
- **Own games**: a player who sat at a game always sees that game (live, result, log, players) and their own game history, regardless of games reach, as long as they can open the event.
- **Game status** (live / ended / aborted, used to pick the right game page): anyone who can open the event.

### What each user can do

- **Create a game**: a player of a regular, unfinished event, who is one of the four players; all four must be players of the event.
- **Enter round results**: any of the four players seated in the game.
- **Enter penalties and false mahjongs**: referees always, in any live game of the event, from the game's view (also a game they play in); at the table, the seated players, unless the event's permissions leave it to referees (area 3), and a seated referee.
- **Undo the last round**: by default referees only (the owner too when they have the referee role): undos are visible in the game's log ("The table", area 4), but a referee at hand still keeps undos honest. The owner may let the seated players undo too, while the game is live (area 3, "Event settings"), for an event whose players they trust. Undoing the final round of an ended game stays with referees: it changes the rating.
- **Abort a game** (regular events): referees at any time; by default also the seated players while no round has been entered (a game started by mistake). The event's permissions can leave it to referees or let the seated players abort at any time (area 3, "Event settings").
- **Finish a game early** (regular events): a referee, and by default a seated player; the event's permissions can leave it to referees.
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

- Used when adding members (a server search as the user types, with a hint "A name (2+ letters) or a user ID") and picking players (the event's players, filtered the same way on the device). A numeric query matches a user ID exactly; otherwise names are searched (exact match first, then prefix, then substring), at least 2 characters, at most 20 results.
- Every player selector shows the user ID next to the name, so namesakes can be told apart. Email is never shown to other users.

### Abuse protection

- At most one auth email (verification or reset) per address per minute and about five per hour; beyond that the response still says "sent", without sending.
- nginx limits requests to `/auth/*` per IP (server config, kept outside the repo). No account lockout after failed sign-ins: it would let anyone lock someone else out.

## 3. Events

### Finding and choosing an event

- **Events list** (`/events`): two tabs, "My (N)" (events where I'm a member) and "Public (N)", so the public events are one tap away however many of my own there are. Without a choice it opens "My", or "Public" when I have no events; the last choice is remembered on the device. A visitor sees only the public events, without tabs. In each tab, active events come first, the ones played most recently on top (the start of their last game, or creation without games), and finished events are collapsed under "Finished (N)", the latest finished on top. A search field above filters both tabs by part of the name, ignoring case (the tab counts show the matches, so a match in the other tab is visible); matching finished events show at once.
- **Select / switch event** (user menu): only events where I'm a member; active ones on top, in the same order as the list, finished ones in a collapsed section (still selectable, to see the last game and history); a search by name from 6 events on. The selected event is where `/new`, `/now`, `/last` and `/history` work.

### Creating an event

- Any signed-in user can create one: name (required, trimmed, at most 100 characters), description (optional, at most 1000), type (regular or tournament; a tournament starts with 5 tours of 120 minutes, the hand in progress finished when time is up, and its owner gets no roles), visibility (private by default) and ratings/games reach with presets. The backend checks the same limits on create and edit.

### Event page

- **Overview**: name, type, finished mark, description, owner, number of players (players are not listed by name: there can be hundreds, and they may not want to be shown; with a player limit, "12 / 16"), names of referees and observers (staff are few and responsible for the event; a section with more than 5 names is collapsed). In a public event that accepts requests, a non-member also sees "Request to join" (see "Requests to join" below).
- **≡ menu**: Manage, Rating, Live games, Finished games, as the user's rights allow; "Finish event" for the owner. The owner sees the number of pending requests to join next to Manage.

### Managing an event (owner)

- **Settings**: name, description, visibility, ratings/games reach. Changing visibility changes a reach only when it becomes invalid (public → members when the event turns private); a custom reach is kept otherwise.
- **Members**: add a user (server-side search by name or ID) with roles; tap a member to switch their roles (a member keeps at least one role: to take all of them away, remove the member); remove a member after a confirmation. The owner can't be removed and may drop any of their roles, even all of them; they stay in the member list and can take a role again. A member in a live game can't be removed, and their player role can't be taken away until the game ends.
- **Overview**: players are a count; referees and observers are listed by name, the first five of a longer list with "N more" for the rest.
- **Joining by invite**: in Manage → Members the owner turns on an invite link and shares it (copy, the phone's share sheet, or its QR code). Whoever opens it sees the event's name (also in a private event: the link is the access) and joins as a player with one tap; a visitor signs in first and comes back to the link; a member is told they already are one. The owner can replace the link (the old one and its QR code stop working) or turn it off. A finished event doesn't take new players by the link.
- **Requests to join**: in Manage → Members the owner can switch on "Accept requests" (off by default; it works only while the event is public). Then a signed-in non-member sees "Request to join" on the event's overview; a visitor signs in first. Only the owner sees the pending requests (Manage → Members, oldest first) and accepts one (the user becomes a player) or declines it. The user sees on the overview that the request is pending and can withdraw it; a declined user sees "Request declined" and can't ask again (the owner can still invite or add them). Pending requests stay when the owner switches requests off or the event turns private, and the owner can still answer them; a user who becomes a member another way loses their pending request. A finished event takes no requests, and its pending ones can't be accepted.
- **Player limit**: optional, set by the owner in Manage → Members. While the event has that many players, a request can't be sent or accepted and the invite link doesn't take new players ("No places left"). The owner may still add players by hand, and lowering the limit below the current count removes nobody.
- **Leaving**: a member who is not the owner and not in a live game can leave a regular event themselves. Leaving a tournament is not possible. Past games and rating rows of a member who left or was removed stay as they are.
- **Deleting**: the owner can delete an event only while it has no games (created by mistake); otherwise it can only be finished.
- **Finishing**: the owner, after a confirmation, when no game is live. It can't be undone: the rating is final. After that only the visibility can change; members are read-only.

### Event settings

Game rules, permissions and the minimum games for the rating are in place; the hand input mode is planned. Two kinds, set by the owner in Manage → Settings, each kind in its own card (game rules also when creating the event), starting from a preset ("MCR standard": 16 hands, minimum 8; "Short": 8 hands, minimum 8) and adjusted value by value. The overview lists the settings that differ from the standard, so players know what they play.

- **Game rules** affect the score: number of hands (16, 8, 4 or 1; 16 = East, South, West, North rounds), minimum hand value (8 by default). False mahjong is always 30 or 60, not a setting. They can be changed only while the event has no games (aborted games don't count); after the first game they are frozen, so all games of one rating follow the same rules. Each game also stores a copy of the rules it started with and is always replayed by them.
- **Permissions** don't affect the score and can be changed at any time (not after the event is finished): who may finish a game early (seated players and referees / referees only), abort it (referees / also seated players before the first round / also seated players at any time). Referees always may. Defaults are the ones in area 1., enter penalties and false mahjongs (referees only / also the seated players; referees always may, from the game's view), undo the last round (referees / also seated players while the game is live; the owner decides without waiting for an undo audit trail, owner's decision of 2026-11-14). Creating games by referees for other players' tables comes with tournament seating (area 8), not as a regular-event setting.
- **Display and input** (can be changed at any time, not after the event is finished): minimum number of ended games for a place in the rating, 1–100, 1 by default (the rating itself is computed the same way; players with fewer are listed below it without a place); hand input mode — total value (now) or picking fans, once fan-based input exists.
- Time limits per game are tournament-only (area 8). Not planned for regular events: a choice of who enters rounds, round confirmation by a second player, other rating methods or placement points, automatic seating.

## 4. Playing a game

Applies to regular events; tournament games follow area 8.

### Starting a game (`/new`)

- A player of the selected event fills the four seats (East, South, West, North) from the event's players, typing part of a name or a user ID, IDs shown next to names; they must include themselves. Clear empties the seats, Shuffle mixes them.
- Players still in a live game **can** be picked: the next table is often prepared while two of its players finish their game. Play then fails with a clear message ("Player X is still in a game — wait until it ends") and the form keeps its seats.
- "Play again" on the result page opens the draft with the same players **in the same seats**; they shuffle themselves if they want.

### The table (`/now`)

- Shows the four seats in seating order, the scores and the round (東1…北4).
- Usually one player at the table keeps the score; any of the four may enter.
- **Entering a round**: tap the round kanji; pick the outcome (draw, self-draw, deal-in, penalty), 👍/👎 next to the seats, the hand value (stepper from the game's minimum hand, 8 by default, or type the number) or the fine (a penalty from the official MCR list, 5, 10, 20, 30, 40, 50 or 60, taken from the offender only; false mahjong 30 or 60; optionally move to the next round). While filling it in, each seat shows the score change the entry will make, so mistakes are seen before OK. After OK the scores update and the entry resets.
- **What a round holds**: exactly one main outcome (draw, self-draw or deal-in), which ends the round, and any number of penalties and false mahjongs, before it or instead of it (a penalty with "move to the next round" ends the round without a main outcome).
- **No silent duplicates**: an entry carries the number of entries in the game's round log as the player's table shows it (results and penalties alike, so a penalty tapped twice is caught too, though it doesn't change the round). If the log has changed since, the entry is refused, the table refreshes, and the message says what changed: "Round 東3 has already been entered", "A penalty has just been entered", or, after a referee's undo, "The last round has been undone" / "The last entry has been undone". The player decides whether to enter it again. Entering two forgotten rounds in a row still works.
- **Refreshing**: scores never change on their own, so a player's picture of the game isn't disturbed. The app checks every 10–15 seconds whether the table changed and then highlights the Refresh button; the player refreshes when ready.
- **Table menu**: round log, fan table, Finish (count the game now, with a confirmation; seated players and referees; not in tournaments), Abort (don't count it; referees, or the seated players before the first round), undo of the last round (referees, or also the seated players if the event's settings allow it).

### The end of a game

- After the last round (the 16th by default, or as the event's rules set) the game ends by itself and the result page (`/last`) opens: final scores, chart, log, "Play again".
- If the game has ended while a player's table was still open (another player entered the last round, or it was finished), their next entry or Refresh shows "The game has ended" and opens the result page; if it was aborted, "The game was aborted" and the home page.
- A referee can undo the last round within an hour after the end; the game becomes live again. After that the result is final. The undo is refused if one of the players has already started another game ("Cannot undo the final round: <Name> has already started another game"): a player can't have two live games.
- **Who entered what, and undos**: every new log entry records who entered it (a seated player, or a referee from the game's view); tapping its row in the log shows "Entered by <name> at <time>" (entries from before this was recorded show only the time). An undone entry stays in the log for everyone who can see it, at the time it was undone: a row like any other (round, scores, Δ), struck through and dimmed, that no longer counts in the scores; tapping it shows who entered it and "Undone by <name> at <time>". A "Hide undone entries" checkbox above the log (only when the game has some) hides them for the viewer.
- **The log in the table menu matches the table**: it reloads together with the table (on Refresh, after an entry or an undo), not when the menu opens, so the referee never undoes an entry their table doesn't show yet.
- **Undo is never applied twice**: like an entry, an undo carries the number of log entries the referee's table shows. If someone has entered or undone an entry since (two referees pressing Undo at once, a double tap), it is refused, the log reloads, and the message says what changed, in the same words as for an entry ("A round has just been entered", "A penalty has just been entered", "The last entry has already been undone").
- **Abandoned games**: a game left live (players walked away) blocks its players and the event's finish. Staff see a mark on live games with no new round for more than ~3 hours and decide themselves (abort or finish); nothing happens automatically.

## 5. Watching

### Game lists

- **Live games** and **Finished games** (event menu, by rights): each game shows the round progress (ring with the round kanji) and the four players with wind and score; finished games load in pages ("Load more").
- Live games has an **Auto-refresh** switch (off by default, remembered on the device). A game that ended stays in Live games for about 10 minutes, marked ✓ with the final scores, then moves to Finished games, so people around the table (or a club screen) can still see the result.

### Watching a game (`/event/N/game/M`)

- The page shows the live table or, for an ended game, the result (same as `/last`). Opening one's own game from the list shows it the same way; one's own table to play is reached from the home page.
- The live table can be rotated (who is at the bottom) and refreshed; Refresh keeps the rotation. Like at the table, scores change only on Refresh by default, and the Refresh button lights up when the game moved on; an **Auto-refresh** switch (remembered on the device) updates the scores by itself, e.g. for a club screen.
- The table menu (log, fans, Finish, Abort, Undo) follows the user's rights for that game: this is where referees manage other people's games.
- **Penalties from the view**: for a referee of the event, the round kanji in the center is a button, in any live game, also one they play in: that is where a referee hands out penalties. It's the referee's own entry into a game, and other referee entries (e.g. round results) may join it later. It opens the same entry as at the table, with only Penalty and False mahjong active (the other outcomes are shown disabled, so the layout is the same); the penalized player is marked with − at their seat. Others, the seated players included, see the kanji as text.
- At the table, when the event leaves penalties to referees, Penalty is shown disabled in the round entry, in its usual place.

### Stream overlay (HUD)

- `GET /api/event/:id/game/:gameId/hud` serves a compact game state for a stream overlay (OBS): round, winds, scores, differences, a message about the last round, the final result. It follows the same access rules as watching the game. It hasn't been updated for a while and must keep up with the rest (penalties, event settings such as 8 hands). The stream is run by a signed-in user, usually an observer (observers see all games of the event). The overlay is a page of the app for OBS's built-in Browser Source, signed in once inside OBS; if that is impractical, a revocable view-only link per event.

## 6. Ratings and history

### Event rating

- A table of players: place, name, games, TP, GP; sortable by any column. A long name is cut with an ellipsis so rows keep one height; the full name shows on hover and on the player's statistics, which a tap opens. The table never scrolls sideways, also at 320 px. Only ended games count (not live or aborted ones). Members who left or were removed stay in it.
- Order: total TP, then total GP, then the best single games by GP. Players equal in all of that share the place (5, 5, then 7).
- The user's own row is highlighted; when it is far down, it stays pinned at the bottom of the screen.
- Players with fewer ended games than the event's minimum (event settings) get no place: in the same table, under a separator row "Fewer than N games", they follow the ranked players in grey, in the same order, with "—" for the place. Places count only the players above the minimum. Sorting by a column sorts each group on its own; the group without places stays under the separator. Members and staff also see there the event's players who haven't played yet (0 games, by name; the separator says "No games yet" when the minimum is 1), so players know who is in the event; others don't, as they don't get the member list.
- Tapping a row (the name included) opens that player's statistics in the event (area 6, "Statistics"); their game history there is planned (for users who can see the event's results; one's own always).

### History

- **History** (`/history`): the user's ended games in the selected event, newest first, 30 at a time with "Load more". On top, a summary for the event: number of games, total and average TP, how many times each place. Each game shows the place, TP, GP, date and the three opponents (exact look to be agreed on a mock-up first). One request returns a page of games with all of that.
- **Across events** (planned): the user picks which of their events to include (events with different rules shouldn't be mixed blindly). Only one's own games: other players' games in events the viewer isn't part of aren't available.

### Statistics

- **Player, in one event**: from the event's ended games. Per game: games, places (the average place; a tap expands 1st–4th with their share of the games; a tie in GP shares the better place), average TP and GP, best and worst game (GP). Per hand: hands played, wins and their share of the hands (a tap expands self-draw and deal-in), average and biggest winning hand, deal-ins (the player dealt in to someone's win) and draws with their shares, penalties (a tap expands point penalties and false mahjongs).
- **Where**: one's own on `/history`, tab "Statistics" next to "Games", for the selected event; anyone's from the event's rating: a tap on a player's row opens `/event/N/player/M` with their name and statistics. One's own statistics are always visible (while the user can open the event); other players' follow the event's ratings reach, as the rating does.
- **Later**: the event's statistics as a whole (the same figures for all players, the biggest hand and who won it); a player's statistics over a chosen set of their events.

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

## 8. Tournaments

Being built in stages. So far: creating a tournament, its status (preparation until the first tour starts, then in progress, then finished), its settings in Manage → Settings (number of tours, tour length, what happens when time is up; the regular-event permissions to finish or abort a game and the rating's minimum games don't apply and are hidden), the substitute role, and the tournament's desk (≡ → Tournament, for its staff) with the participants and their numbers: the referee's draw, redone after the players change, and withdrawals once the first tour has started, when the players are closed (they can't be added or removed, also not by the invite link or a request). No games can be played in a tournament yet. The rules follow the official *Mahjong Competition Rules* (WMO / WMCC, 2021 edition) where they cover a question; where the rules are silent, these are the organizer's practice. A few questions are still open with referees; they are listed at the end.

### Format

- 12–160 players, 3–15 rounds (tours), 1–3 days, in one hall or several (the common screen is shown in each).
- A tour is one game (hanchan) at every table: 16 hands, or until the tour's time is up.
- No team ranking, stages or cuts at first. A team may later be an optional label on a participant, so the data should allow it.

### Roles

- **Owner**: the organizer. Sees everything, manages the settings, participants and roles, finishes the tournament. Gets no roles by default: referee actions take the referee role, playing takes the player role (as in all events, area 1).
- **Referee**: runs the tournament — the draw of participant numbers; the seating (generate, edit, cancel, publish); starting, aborting and finishing tours; adding time to a table; seating substitutes; penalties and undoing the last entry, at any time.
- **Observer**: staff; sees everything, live tables included, at any time (e.g. the person streaming).
- **Player**: plays and enters the rounds at their own table.
- **Substitute**: a role of its own, exclusive with player. Sits only where a referee seats them. Referee and observer combine with either player or substitute.

### Participants

- Everyone registers in advance with an account (no guests). Ways in: the invite link, being added by the organizer, or a request to join that the organizer accepts or declines.
- A participant can't leave a tournament themselves.
- **Participant numbers** come from a draw (a shuffle), a referee's button on the tournament's desk that must be pressed before the first tour's seating is published. Random numbers, so nobody can pick opponents by the order of registration. The desk lists the players by number for the staff (referees, observers, the owner); a redraw asks for a confirmation.
- Before the first tour starts, the participants may change: after a change, the numbers no longer count and the desk asks for a new draw (cancelling the seating first if it was published).
- **Once the first tour starts, the list of players is closed.** After that only staff and substitutes can be added.
- **A player who withdraws** keeps their past games in the rating, is marked "withdrew" and isn't seated any more; a substitute takes their number in the seating.
- **Late players**: the app doesn't track lateness. The referee decides and seats or removes a substitute by hand. Under the official rules, being late by up to 15 minutes costs 10 or 20 points (an out-of-game penalty), and by more than 15 minutes the player misses that tour and a substitute plays; from the next tour the player plays again.

### Visibility

- Non-members see what they see in any event (the number of players, the staff by name).
- Participants see the other participants: names in the rating table, only the number on the overview.
- **Live tables during a tour** are visible only to staff and to the players of the tour who have finished their own game. Everyone else, anonymous visitors included, doesn't see them, so signing out doesn't get around it.
- A setting hides the rating after tour N until the end; staff always see it.

### Seating

- The method is chosen for each tour:
  - random, with as few repeat meetings as possible;
  - Swiss without repeats: neighbours in the current rating play together, repeat meetings are avoided;
  - interval: by the current rating with step k (by default the number of tables), e.g. table 1 gets places 1, 1+k, 1+2k, 1+3k;
  - fixed without repeats for all tours at once;
  - imported from a file (seatings prepared by the organizer).
- All methods come from **one seating generator of our own**. It searches for a limited time and shows the referee the number of repeat meetings. No repeats can't be guaranteed near the limit of (N−1)/3 tours for N players.
- Table and wind are assigned automatically, by participant number. Winds are not balanced on purpose: in MCR their effect is minimal.
- Before a tour starts, the referee can regenerate or edit its seating; after the start, only through substitutes.

### A tour, step by step

1. The referee prepares the seating and **publishes** it: the tour's games are created (status `created`). Each player sees their own table, seat and opponents as "Current game", not the other tables. Rounds can't be entered yet.
2. **Cancel seating** (before the start): the created games are deleted; the seating can be edited or made again.
3. **Start tour**: the tour's timer starts, the games become live, the tables enter rounds.
4. The next tour's seating can be published only after this tour is finished: a player must not have a created and a live game at once, and the Swiss seating depends on the tour's results anyway.
5. **Abort tour** (after the start, before finishing; with a confirmation): for a critical mistake noticed too late (e.g. a wrong seating). All the tour's games are aborted, the tour doesn't count and the tour number doesn't advance; the tour starts over from step 1. The aborted games, empty ones included, stay in the database with the aborted status (each attempt has the same number of games); they don't count in the rating and aren't listed anywhere for now.
6. **Finish tour**: available when all the tour's games have ended; the referee sees which tables haven't. Only then do the tour's results enter the rating, and after it the last entry can't be undone any more. Until then the referee can always undo, players as the settings allow. A finished tour can't be reopened: its results are final.
- The referee can add time to a single table (e.g. one that started late through the organizers' fault).
- No planned start time per tour at first: the timer starts with the button.

### Timer

- One timer per tour, started by Start tour; time is up automatically. The referee may add time to a table.
- **What happens when time is up** is a tournament setting:
  - the hand in progress is played to the end (**default**, the organizer's practice);
  - the hand in progress and N more hands (usually 1–2);
  - the game ends at once and the hand in progress doesn't count (the official rule).

  Hands that the setting lets the table finish are entered after the signal.
- The timer shows on the table screen and on the common screen, with no warnings as the end nears. It is the one exception to "the table screen changes only on the player's tap" (area 4): only the timer ticks.

### Playing

- Games are created by the referee publishing a tour's seating, never by players.
- No one can finish a game early or abort a single game in a tournament (the referee aborts whole tours, above).
- **Only the table enters rounds** (not the referee). No paper score sheets.
- Mistakes are fixed only by undoing the last entry, one at a time; there is no editing an entry in the middle of the game.
- **Penalties in a game**: the referee always; players only if the setting "players may enter penalties" is on. The amounts follow the official list (5, 10, 20, 30, 40, 50 or 60 points), taken from the offender only.
- **Penalties outside a hand** (lateness, conduct): by the referee, tied to the player's game in that tour: points off that game's result, with a reason, from the same list. Being part of the game's score, they can change the places at that table and so the tournament points.
- Undoing the last entry: the referee always; players only if the setting "players may undo" is on; only until the tour is finished.
- Players don't confirm a game's result. A confirmation screen before entering a hand may come later, for all events.

### Substitutes

- One mechanism: **another person takes a seat from hand N on**, for the whole tour or in the middle of a game.
- The referee seats and removes substitutes by hand from the table menu. From then on the substitute has the game as "Current game" on the home page, sees the table and enters rounds; the replaced player can no longer enter rounds but still sees their game. A substitute can't sit in two live games at once.
- The seat's score carries on. The log may mark from which hand someone else sat, if it doesn't clutter the log.
- **Points**: the opponents get tournament points for their places as usual; the place held by a substitute gives tournament points to nobody. A substitute doesn't take part in the rating: they play only so that the table is complete. The game's record shows who actually sat.
- **Substitutes in the rating** are hidden by default; a switch on the rating page shows them, always at the bottom whatever their points (their results mean little: substitutes swap with each other).
- **Data**: a seat keeps a history of who sat from which hand, from the start, even if mid-game substitution comes later.
- **An incomplete table**: no skipped tour, no three-player game. For now a substitute must be found; when starting a tour, the referee sees "N players missing". An "empty seat" (a technical player, the table effectively playing three) may come later, once it's decided how to score it.

### Rating

- The main table of a tournament is by tour: tournament points (TP) and game points (GP) for each tour, and the totals.
- It shows the standings after the last finished tour, with a line "Tour N in progress".
- Order: TP, then GP (official rules). TP for places are 4, 2, 1, 0, fixed.
- Players who withdrew are listed below the active players, marked "withdrew" (tentative).
- Export to CSV.

### Screens

- **A player's home page**: the current tour and table ("Current game"). Everything else is on the event's pages.
- **The common screen** for the hall (projector, TV): the current tour's seating (all tables: players find their table by it), the timer, the rating; it refreshes by itself.
- No push notifications.

### Settings and when they freeze

| Setting | Can change |
|---|---|
| Game rules (number of hands, minimum hand) | until the first game |
| What happens when time is up | until the first tour starts |
| Number of tours | until the last tour starts; prepared future tours are recomputed |
| Tour length | until that tour starts; during it, only adding time to a table |
| Seating method of a tour | until that tour's seating is published |
| "Players may enter penalties", "players may undo" | any time |
| "Hide the rating after tour N" | any time |

### Lifecycle

Preparation (registration, draw, seating) → in progress (tours) → finished (by the owner's button).

### First stage and later

- **First stage**:
  - the tournament and its participants (invite link, adding by hand, requests to join);
  - the roles, substitute included;
  - the draw of participant numbers;
  - tours: publish, start, abort, finish;
  - the seating, all methods and import, with our own generator;
  - the timer and what happens when time is up;
  - games from the seating, rounds entered by the table;
  - the referee: penalties in and outside games, undo, substitutes for a whole tour;
  - live tables hidden during a tour;
  - the rating by tour, CSV;
  - the common screen;
  - maybe mid-game substitutes.
- **Later**: mid-game substitutes (if not in the first stage), the empty seat, teams, the confirmation screen before entering a hand, planned tour start times, stages and cuts.

### Open questions (asking referees)

- Does a player who misses a tour (absent, or late by more than 15 minutes) get 0 TP and 0 GP for it?
- A player played part of a game and left, and a substitute finished it: does that player get anything for the game?
- In a game with a substitute, do the opponents get GP from the actual result and TP for their places, with nobody getting the substitute's place's TP?
- When a player and a substitute share a place (e.g. 2nd–3rd), does the player get their share of TP (1.5)?
- What breaks a tie after TP and GP? If everything is equal, do the players share the place?
