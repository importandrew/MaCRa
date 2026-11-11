import type { Messages } from '../types';

// ru holds the Russian texts. Messages makes the build fail on a missing or
// misspelled key; Russian plurals add the `_few` and `_many` forms.
const ru: Messages = {
  header: {
    events: 'События',
    logIn: 'Войти',
    userMenu: 'Меню пользователя',
    profile: 'Профиль',
    settings: 'Настройки',
    logOut: 'Выйти',
    noEventSelected: 'Событие не выбрано',
    selectEvent: 'Выбрать событие',
    switchEvent: 'Сменить событие',
  },
  landing: {
    rating: 'Рейтинг',
    liveGames: 'Идущие игры',
    finishedGames: 'Завершённые игры',
    loginToPlay: 'Войти, чтобы играть',
    selectEventToPlay: 'Выбрать событие для игры',
    loading: 'Загрузка...',
    currentGame: 'Текущая игра',
    newGame: 'Новая игра',
    lastGame: 'Последняя игра',
    playedGames: 'Сыгранные игры',
  },
  settings: {
    title: 'Настройки',
    language: 'Язык',
  },
};

export default ru;
