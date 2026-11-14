import { currentLanguage } from '../i18n';

// getDateTimeFormatted formats a date and time in the app's language, so the
// month name matches the rest of the screen.
export const getDateTimeFormatted = (dateISO: string): string => {
  if (!dateISO) return '—';

  const locale = currentLanguage();
  const dateObj = new Date(dateISO);
  const date = dateObj.toLocaleDateString(locale, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
  const time = dateObj.toLocaleTimeString(locale, {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  });

  return `${date}, ${time}`;
};

// getTimeFormatted formats the time of day (hours and minutes) in the app's
// language, e.g. for entries of one game.
export const getTimeFormatted = (dateISO: string): string => {
  if (!dateISO) return '—';

  return new Date(dateISO).toLocaleTimeString(currentLanguage(), {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  });
};
