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
  const [small, big] = rules.false_mahjong;
  if (
    small !== STANDARD_RULES.false_mahjong[0] ||
    big !== STANDARD_RULES.false_mahjong[1]
  ) {
    differences.push(t('rules.falseMahjongValue', { small, big }));
  }
  return differences;
}
