export const getDateTimeFormatted = (dateISO: string): string => {
  if (!dateISO) return '—';

  const dateObj = new Date(dateISO);
  const date = dateObj.toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
  const time = dateObj.toLocaleTimeString(undefined, {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  });

  return `${date}, ${time}`;
};
