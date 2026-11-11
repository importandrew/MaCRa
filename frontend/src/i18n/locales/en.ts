// en holds the English texts, the reference for every other language: its
// shape (keys and nesting) is what each translation must match, see Messages.
//
// Keys are grouped by screen or component. Plurals use i18next suffixes
// (`_one`, `_other` in English; Russian also has `_few` and `_many`) and are
// called as t('key', { count }). Interpolation: {{name}}.
const en = {
  header: {
    events: 'Events',
    logIn: 'Login',
    userMenu: 'User menu',
    profile: 'Profile',
    settings: 'Settings',
    logOut: 'Log Out',
    noEventSelected: 'No event selected',
    selectEvent: 'Select event',
    switchEvent: 'Switch event',
  },
  landing: {
    rating: 'Rating',
    liveGames: 'Live games',
    finishedGames: 'Finished games',
    loginToPlay: 'Login to play',
    selectEventToPlay: 'Select event to play',
    loading: 'Loading...',
    currentGame: 'Current game',
    newGame: 'New game',
    lastGame: 'Last game',
    playedGames: 'Played games',
  },
  settings: {
    title: 'Settings',
    language: 'Language',
  },
};

export default en;
