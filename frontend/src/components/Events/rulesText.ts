import type { TFunction } from 'i18next';
import { type GameRules, STANDARD_RULES } from '../../models/GameRules';

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
