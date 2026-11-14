import type { TFunction } from 'i18next';
import { type GameRules, STANDARD_RULES } from '../../models/GameRules';
import type { EventSettings } from '../../models/EventSettings';

// rulesDifferences lists, as texts, the rules that differ from the MCR
// standard ones, so players know what they play; empty for the standard.
export function rulesDifferences(rules: GameRules, t: TFunction): string[] {
  const differences: string[] = [];
  if (rules.hands !== STANDARD_RULES.hands) {
    differences.push(t('rules.handsCount', { count: rules.hands }));
  }
  if (rules.min_hand !== STANDARD_RULES.min_hand) {
    differences.push(t('rules.minHandValue', { value: rules.min_hand }));
  }
  return differences;
}

// permissionDifferences lists, as texts, the permissions that differ from
// the defaults; empty for the defaults.
export function permissionDifferences(
  settings: EventSettings,
  t: TFunction
): string[] {
  const differences: string[] = [];
  if (settings.finish_game === 'referees') {
    differences.push(t('eventSettings.finishByReferees'));
  }
  if (settings.abort_game === 'referees') {
    differences.push(t('eventSettings.abortByReferees'));
  } else if (settings.abort_game === 'seated') {
    differences.push(t('eventSettings.abortBySeated'));
  }
  return differences;
}
