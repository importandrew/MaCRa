// GameRules mirrors models.GameRules: the rules that affect a game's score.
// An event has them (frozen after its first game) and each game keeps the
// copy it started with.
export type GameRules = {
  // hands: 16 (East, South, West, North rounds) or 8 (East and South).
  hands: number;
  min_hand: number;
  // false_mahjong: the two amounts a false mahjong may cost, smaller first.
  false_mahjong: [number, number];
};

// Limits of the rules, as the backend checks them.
export const HANDS_OPTIONS = [16, 8];
export const MIN_MIN_HAND = 8;
export const MAX_MIN_HAND = 88;
export const MAX_FALSE_MAHJONG = 300;

export const STANDARD_RULES: GameRules = {
  hands: 16,
  min_hand: 8,
  false_mahjong: [30, 60],
};

export const SHORT_RULES: GameRules = {
  hands: 8,
  min_hand: 24,
  false_mahjong: [30, 60],
};

export const RULE_PRESETS = { standard: STANDARD_RULES, short: SHORT_RULES };
export type RulePreset = keyof typeof RULE_PRESETS;

export function sameRules(a: GameRules, b: GameRules): boolean {
  return (
    a.hands === b.hands &&
    a.min_hand === b.min_hand &&
    a.false_mahjong[0] === b.false_mahjong[0] &&
    a.false_mahjong[1] === b.false_mahjong[1]
  );
}

// isValidFalseMahjong: a positive multiple of 3 up to MAX_FALSE_MAHJONG (the
// other three players get a third each).
export function isValidFalseMahjong(value: number): boolean {
  return (
    Number.isInteger(value) &&
    value > 0 &&
    value <= MAX_FALSE_MAHJONG &&
    value % 3 === 0
  );
}

export function isValidRules(rules: GameRules): boolean {
  const [small, big] = rules.false_mahjong;
  return (
    HANDS_OPTIONS.includes(rules.hands) &&
    Number.isInteger(rules.min_hand) &&
    rules.min_hand >= MIN_MIN_HAND &&
    rules.min_hand <= MAX_MIN_HAND &&
    isValidFalseMahjong(small) &&
    isValidFalseMahjong(big) &&
    small < big
  );
}
