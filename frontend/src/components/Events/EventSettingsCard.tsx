import { useTranslation } from 'react-i18next';
import { Card, Form, Label, Radio, RadioGroup } from '@heroui/react';
import { type FormEvent, useState } from 'react';
import type { Event } from '../../models/Event';
import {
  ABORT_GAME_OPTIONS,
  type EventSettings,
  FINISH_GAME_OPTIONS,
  isValidSettings,
  MAX_MIN_RATING_GAMES,
  MIN_MIN_RATING_GAMES,
  sameSettings,
} from '../../models/EventSettings';
import ValueStepper from '../ValueStepper';
import { updateEventSettings } from '../../providers/event/eventApi';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';

// SectionHeading starts a section of the card: its heading and a hint.
function SectionHeading({ title, hint }: { title: string; hint: string }) {
  return (
    <div className="flex flex-col gap-1">
      <h2 className="font-medium">{title}</h2>
      <p className="text-sm text-muted">{hint}</p>
    </div>
  );
}

// SettingRadios is one setting: a vertical list of its options (the labels
// are too long for a SegmentedPicker on a phone).
function SettingRadios<T extends string>({
  label,
  options,
  value,
  onChange,
  isDisabled,
  optionLabel,
}: {
  label: string;
  options: T[];
  value: T;
  onChange: (value: T) => void;
  isDisabled: boolean;
  optionLabel: (option: T) => string;
}) {
  return (
    // secondary: the grey controls HeroUI uses on a card's surface.
    <RadioGroup
      variant="secondary"
      value={value}
      onChange={(v) => onChange(v as T)}
      isDisabled={isDisabled}
    >
      <Label>{label}</Label>
      {options.map((option) => (
        <Radio key={option} value={option}>
          <Radio.Content>
            <Radio.Control>
              <Radio.Indicator />
            </Radio.Control>
            {optionLabel(option)}
          </Radio.Content>
        </Radio>
      ))}
    </RadioGroup>
  );
}

// EventSettingsCard is the card under the game rules in Manage → Settings:
// the event's settings that don't affect the score, in sections with their
// own headings (Permissions: who may finish or abort a game; Rating: the
// minimum games for a place). They can change at any time, except in a
// finished event, and are saved together.
export default function EventSettingsCard({
  event,
  onSaved,
}: {
  event: Event;
  onSaved: () => Promise<void>;
}) {
  const { t } = useTranslation();
  // draft holds the owner's edits; until there are any (and again after a
  // save), the form shows the event's settings.
  const [draft, setDraft] = useState<EventSettings | null>(null);
  const settings = draft ?? event.settings;
  const { isPending, run } = useSubmitOnce();
  const isFinished = Boolean(event.finished_at);

  const handleSave = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    void run(async () => {
      if (await updateEventSettings(event.id, settings)) {
        showToast({ title: t('eventSettings.saved'), variant: 'success' });
        // The edits stay on screen until the event has reloaded with them:
        // clearing them first would show the old values for a moment.
        await onSaved();
        setDraft(null);
      }
    });
  };

  return (
    <>
      <Card>
        <Card.Content className="gap-4">
          <p className="text-sm text-muted">
            {isFinished
              ? t('eventSettings.finishedReadOnly')
              : t('eventSettings.hint')}
          </p>
          <Form
            id="event-settings-form"
            onSubmit={handleSave}
            className="flex flex-col gap-4"
          >
            <section className="flex flex-col gap-4">
              <SectionHeading
                title={t('eventSettings.permissions')}
                hint={t('eventSettings.permissionsHint')}
              />
              <SettingRadios
                label={t('eventSettings.finishGame')}
                options={FINISH_GAME_OPTIONS}
                value={settings.finish_game}
                onChange={(finish_game) =>
                  setDraft({ ...settings, finish_game })
                }
                isDisabled={isFinished}
                optionLabel={(option) => t(`eventSettings.finish.${option}`)}
              />
              <SettingRadios
                label={t('eventSettings.abortGame')}
                options={ABORT_GAME_OPTIONS}
                value={settings.abort_game}
                onChange={(abort_game) => setDraft({ ...settings, abort_game })}
                isDisabled={isFinished}
                optionLabel={(option) => t(`eventSettings.abort.${option}`)}
              />
            </section>
            <section className="flex flex-col gap-4">
              <SectionHeading
                title={t('eventSettings.rating')}
                hint={t('eventSettings.ratingHint')}
              />
              <div className="flex flex-col gap-1.5">
                <p className="text-sm">{t('eventSettings.minRatingGames')}</p>
                {isFinished ? (
                  <p>{settings.min_rating_games}</p>
                ) : (
                  <>
                    <ValueStepper
                      label={t('eventSettings.minRatingGames')}
                      value={settings.min_rating_games}
                      min={MIN_MIN_RATING_GAMES}
                      onChange={(min_rating_games) =>
                        setDraft({ ...settings, min_rating_games })
                      }
                    />
                    <p
                      className={`text-xs text-center ${
                        settings.min_rating_games > MAX_MIN_RATING_GAMES
                          ? 'text-danger'
                          : 'text-muted'
                      }`}
                    >
                      {t('eventSettings.minRatingGamesRange', {
                        min: MIN_MIN_RATING_GAMES,
                        max: MAX_MIN_RATING_GAMES,
                      })}
                    </p>
                  </>
                )}
              </div>
            </section>
          </Form>
        </Card.Content>
      </Card>
      {!isFinished && (
        <PendingButton
          type="submit"
          form="event-settings-form"
          variant="primary"
          isPending={isPending}
          isDisabled={
            !isValidSettings(settings) || sameSettings(settings, event.settings)
          }
          className="w-3/4 self-center"
        >
          {t('eventSettings.save')}
        </PendingButton>
      )}
    </>
  );
}
