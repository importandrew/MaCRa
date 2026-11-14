import { useTranslation } from 'react-i18next';
import {
  type GameRules,
  HANDS_OPTIONS,
  MAX_MIN_HAND,
  MIN_MIN_HAND,
  RULE_PRESETS,
  type RulePreset,
  sameRules,
} from '../../models/GameRules';
import SegmentedPicker from '../SegmentedPicker';
import ValueStepper from '../ValueStepper';

const PRESETS = Object.keys(RULE_PRESETS) as RulePreset[];

// GameRulesFields are an event's game rules, the same when creating an event
// and in Manage → Rules: a preset to start from, then each value. With
// isReadOnly they only show the rules.
export default function GameRulesFields({
  value,
  onChange,
  isReadOnly = false,
}: {
  value: GameRules;
  onChange: (value: GameRules) => void;
  isReadOnly?: boolean;
}) {
  const { t } = useTranslation();
  const preset =
    PRESETS.find((name) => sameRules(RULE_PRESETS[name], value)) ?? null;

  return (
    <>
      {!isReadOnly && (
        <SegmentedPicker
          label={t('rules.preset')}
          options={PRESETS.map((name) => ({
            value: name,
            label: t(`rules.presets.${name}`),
          }))}
          selectedKey={preset}
          onSelectionChange={(name) => onChange(RULE_PRESETS[name])}
        />
      )}
      <SegmentedPicker
        label={t('rules.hands')}
        options={HANDS_OPTIONS.map((hands) => ({
          value: String(hands),
          label: String(hands),
        }))}
        selectedKey={String(value.hands)}
        onSelectionChange={(hands) =>
          onChange({ ...value, hands: Number(hands) })
        }
        isOptionDisabled={() => isReadOnly}
      />
      <div className="flex flex-col gap-1.5">
        <p className="text-sm">{t('rules.minHand')}</p>
        {isReadOnly ? (
          <p>{value.min_hand}</p>
        ) : (
          <>
            <ValueStepper
              label={t('rules.minHand')}
              value={value.min_hand}
              min={MIN_MIN_HAND}
              onChange={(minHand) => onChange({ ...value, min_hand: minHand })}
            />
            <p
              className={`text-xs text-center ${
                value.min_hand > MAX_MIN_HAND ? 'text-danger' : 'text-muted'
              }`}
            >
              {t('rules.minHandRange', {
                min: MIN_MIN_HAND,
                max: MAX_MIN_HAND,
              })}
            </p>
          </>
        )}
      </div>
    </>
  );
}
