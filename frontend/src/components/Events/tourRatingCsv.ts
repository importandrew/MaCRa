import type { TFunction } from 'i18next';
import type { TourRating } from '../../models/TourRating';

// csvField quotes a value when it needs it (a comma, a quote, a line break).
function csvField(value: string | number): string {
  const text = String(value);
  return /[",\r\n]/.test(text) ? `"${text.replace(/"/g, '""')}"` : text;
}

// tourRatingCsv writes the rating by tour as CSV: place, participant
// number, name, TP and GP of each finished tour, the totals, and whether the
// player withdrew. A BOM first, so spreadsheets read the names as UTF-8.
export function tourRatingCsv(rating: TourRating, t: TFunction): string {
  const header: (string | number)[] = [
    t('tourRating.csv.place'),
    t('tourRating.csv.number'),
    t('tourRating.csv.name'),
    ...rating.tours.flatMap((tour) => [
      t('tourRating.csv.tourTP', { n: tour }),
      t('tourRating.csv.tourGP', { n: tour }),
    ]),
    t('rating.tp'),
    t('rating.gp'),
    t('tourRating.csv.withdrew'),
  ];
  const lines = rating.rows.map((row): (string | number)[] => [
    row.rank ?? '',
    row.number ?? '',
    row.name,
    ...row.tours.flatMap((result): (string | number)[] =>
      result ? [result.tp, result.gp] : ['', '']
    ),
    row.tp,
    row.gp,
    row.withdrawn ? t('tourRating.csv.yes') : '',
  ]);
  return (
    '﻿' +
    [header, ...lines]
      .map((line) => line.map(csvField).join(','))
      .join('\r\n') +
    '\r\n'
  );
}

// downloadText gives the viewer a text file to save.
export function downloadText(text: string, fileName: string, type: string) {
  const url = URL.createObjectURL(new Blob([text], { type }));
  const link = document.createElement('a');
  link.href = url;
  link.download = fileName;
  link.click();
  URL.revokeObjectURL(url);
}
