export const WIND_LABELS = ['East', 'South', 'West', 'North'];
export const WIND_SYMBOLS = ['東', '南', '西', '北'];
export const OUTCOME_LABELS = ['Draw', 'Self', 'Deal-in'];

/**
 * Calculates the wind label (kanji + position) from round number
 * @param round Round number (1-16)
 * @returns Wind label string (e.g., "東1", "南2", "西3", "北4")
 */
export function getWindLabel(round: number): string {
  const windIndex = Math.floor((round - 1) / 4);
  const position = ((round - 1) % 4) + 1;
  const wind = WIND_SYMBOLS[windIndex];
  return `${wind}${position}`;
}
