import { useTranslation } from 'react-i18next';
import { Card, FieldError, Form, Input, Label, TextField } from '@heroui/react';
import { type FormEvent, useState } from 'react';
import {
  type GameRules,
  HANDS_OPTIONS,
  isValidFalseMahjong,
  isValidRules,
  MAX_FALSE_MAHJONG,
  MAX_MIN_HAND,
  MIN_MIN_HAND,
  RULE_PRESETS,
  type RulePreset,
  sameRules,
} from '../../models/GameRules';
import type { Event } from '../../models/Event';
import { updateEventRules } from '../../providers/event/eventApi';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import SegmentedPicker from '../SegmentedPicker';
import ValueStepper from '../ValueStepper';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';

const PRESETS = Object.keys(RULE_PRESETS) as RulePreset[];

// EventRulesForm is Manage → Rules: the event's game rules, starting from a
// preset and adjusted value by value. They can change only until the first
// game (the server refuses them after it) and not in a finished event; then
// the form only shows them.
export default function EventRulesForm({
  event,
  onSaved,
}: {
  event: Event;
  onSaved: () => void;
}) {
  const { t } = useTranslation();
  // draft holds the owner's edits; until there are any (and again after a
  // save), the form shows the event's rules.
  const [draft, setDraft] = useState<GameRules | null>(null);
  const rules = draft ?? event.rules;
  const { isPending, run } = useSubmitOnce();

  const isFinished = Boolean(event.finished_at);
  const isReadOnly = isFinished || Boolean(event.rules_frozen);
  const preset =
    PRESETS.find((name) => sameRules(RULE_PRESETS[name], rules)) ?? null;
  const [small, big] = rules.false_mahjong;

  const setFalseMahjong = (index: 0 | 1, text: string) => {
    const digits = text.replace(/\D/g, '');
    const amounts: [number, number] = [...rules.false_mahjong];
    amounts[index] = digits === '' ? 0 : Number(digits);
    setDraft({ ...rules, false_mahjong: amounts });
  };

  const handleSave = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    void run(async () => {
      if (await updateEventRules(event.id, rules)) {
        showToast({ title: t('rules.saved'), variant: 'success' });
        setDraft(null);
        onSaved();
      }
    });
  };

  return (
    <>
      <Card>
        <Card.Content className="gap-4">
          <p className="text-sm text-muted">
            {isFinished
              ? t('rules.finishedReadOnly')
              : event.rules_frozen
                ? t('rules.frozen')
                : t('rules.editableHint')}
          </p>
          <Form
            id="event-rules-form"
            onSubmit={handleSave}
            className="flex flex-col gap-4"
          >
            {!isReadOnly && (
              <SegmentedPicker
                label={t('rules.preset')}
                options={PRESETS.map((name) => ({
                  value: name,
                  label: t(`rules.${name}`),
                }))}
                selectedKey={preset}
                onSelectionChange={(name) => setDraft(RULE_PRESETS[name])}
              />
            )}
            <SegmentedPicker
              label={t('rules.hands')}
              options={HANDS_OPTIONS.map((hands) => ({
                value: String(hands),
                label: t(hands === 8 ? 'rules.hands8' : 'rules.hands16'),
              }))}
              selectedKey={String(rules.hands)}
              onSelectionChange={(hands) =>
                setDraft({ ...rules, hands: Number(hands) })
              }
              isOptionDisabled={() => isReadOnly}
            />
            <div className="flex flex-col gap-1.5">
              <p className="text-sm">{t('rules.minHand')}</p>
              {isReadOnly ? (
                <p>{rules.min_hand}</p>
              ) : (
                <>
                  <ValueStepper
                    label={t('rules.minHand')}
                    value={rules.min_hand}
                    min={MIN_MIN_HAND}
                    onChange={(minHand) =>
                      setDraft({ ...rules, min_hand: minHand })
                    }
                  />
                  <p
                    className={`text-xs text-center ${
                      rules.min_hand > MAX_MIN_HAND
                        ? 'text-danger'
                        : 'text-muted'
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
            <div className="flex flex-col gap-1.5">
              <p className="text-sm">{t('rules.falseMahjong')}</p>
              <p className="text-xs text-muted">
                {t('rules.falseMahjongHint')}
              </p>
              <div className="grid grid-cols-2 gap-3">
                <TextField
                  isDisabled={isReadOnly}
                  value={small ? String(small) : ''}
                  onChange={(text) => setFalseMahjong(0, text)}
                  isInvalid={!isReadOnly && !isValidFalseMahjong(small)}
                >
                  <Label>{t('rules.falseMahjongSmall')}</Label>
                  <Input inputMode="numeric" pattern="\d*" autoComplete="off" />
                  <FieldError>
                    {t('rules.falseMahjongInvalid', { max: MAX_FALSE_MAHJONG })}
                  </FieldError>
                </TextField>
                <TextField
                  isDisabled={isReadOnly}
                  value={big ? String(big) : ''}
                  onChange={(text) => setFalseMahjong(1, text)}
                  isInvalid={
                    !isReadOnly && (!isValidFalseMahjong(big) || big <= small)
                  }
                >
                  <Label>{t('rules.falseMahjongBig')}</Label>
                  <Input inputMode="numeric" pattern="\d*" autoComplete="off" />
                  <FieldError>
                    {isValidFalseMahjong(big)
                      ? t('rules.falseMahjongOrder')
                      : t('rules.falseMahjongInvalid', {
                          max: MAX_FALSE_MAHJONG,
                        })}
                  </FieldError>
                </TextField>
              </div>
            </div>
          </Form>
        </Card.Content>
      </Card>
      {!isReadOnly && (
        <PendingButton
          type="submit"
          form="event-rules-form"
          variant="primary"
          isPending={isPending}
          isDisabled={!isValidRules(rules) || sameRules(rules, event.rules)}
          className="w-3/4 self-center"
        >
          {t('rules.save')}
        </PendingButton>
      )}
    </>
  );
}
