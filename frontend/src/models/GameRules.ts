// GameRules mirrors models.GameRules: the rules that affect a game's score.
// An event has them (frozen after its first game) and each game keeps the
// copy it started with.
export type GameRules = {
  // hands: 16 (East, South, West, North rounds), 8, 4 or 1.
  hands: number;
  min_hand: number;
};

// Limits of the rules, as the backend checks them.
export const HANDS_OPTIONS = [16, 8, 4, 1];
export const MIN_MIN_HAND = 8;
export const MAX_MIN_HAND = 88;

// FALSE_MAHJONG_AMOUNTS: what a false mahjong may cost, in every game.
export const FALSE_MAHJONG_AMOUNTS = [30, 60];

// PENALTY_AMOUNTS: what an in-game penalty may cost, the official MCR list
// (models.PenaltyAmounts on the backend).
export const PENALTY_AMOUNTS = [5, 10, 20, 30, 40, 50, 60];

export const STANDARD_RULES: GameRules = { hands: 16, min_hand: 8 };

export const SHORT_RULES: GameRules = { hands: 8, min_hand: 8 };

export const RULE_PRESETS = { standard: STANDARD_RULES, short: SHORT_RULES };
export type RulePreset = keyof typeof RULE_PRESETS;

export function sameRules(a: GameRules, b: GameRules): boolean {
  return a.hands === b.hands && a.min_hand === b.min_hand;
}

export function isValidRules(rules: GameRules): boolean {
  return (
    HANDS_OPTIONS.includes(rules.hands) &&
    Number.isInteger(rules.min_hand) &&
    rules.min_hand >= MIN_MIN_HAND &&
    rules.min_hand <= MAX_MIN_HAND
  );
}
