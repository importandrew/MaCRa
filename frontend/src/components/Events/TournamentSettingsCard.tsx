import { useTranslation } from 'react-i18next';
import { Card, Form } from '@heroui/react';
import { type FormEvent, useState } from 'react';
import type { Event } from '../../models/Event';
import {
  isValidTournamentSettings,
  MAX_EXTRA_HANDS,
  MAX_TOUR_MINUTES,
  MAX_TOURS,
  MIN_EXTRA_HANDS,
  MIN_TOUR_MINUTES,
  MIN_TOURS,
  TIES_CHOICES,
  TIME_UP_RULES,
  tiesChoice,
  type TournamentSettings,
  withTies,
} from '../../models/TournamentSettings';
import { updateTournament } from '../../providers/event/eventApi';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import ListStepper from '../ListStepper';
import PendingButton from '../PendingButton';
import ValueStepper from '../ValueStepper';
import { showToast } from '../Toast';
import { SectionHeading, SettingRadios } from './EventSettingsCard';

// range lists the whole numbers from min to max.
const range = (min: number, max: number) =>
  Array.from({ length: max - min + 1 }, (_, i) => min + i);

// TournamentSettingsCard is a tournament's own settings in Manage →
// Settings: the number of tours, the tour length and what happens when time
// is up. Saved together; a finished tournament shows them read-only.
export default function TournamentSettingsCard({
  event,
  onSaved,
}: {
  event: Event;
  onSaved: () => Promise<void>;
}) {
  const { t } = useTranslation();
  const [draft, setDraft] = useState<TournamentSettings | null>(null);
  const saved = event.tournament;
  const settings = draft ?? saved;
  const { isPending, run } = useSubmitOnce();
  const isFinished = Boolean(event.finished_at);
  // What happens when time is up is frozen once the first tour started.
  const isStarted = Boolean(event.started_at);

  if (!saved || !settings) return null;

  // tie_breakers may come as null: compare as an empty list.
  const normalized = (s: TournamentSettings) =>
    JSON.stringify({ ...s, tie_breakers: s.tie_breakers ?? [] });
  const isChanged = normalized(settings) !== normalized(saved);

  const handleSave = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    void run(async () => {
      if (await updateTournament(event.id, settings)) {
        showToast({ title: t('tournament.saved'), variant: 'success' });
        await onSaved();
        setDraft(null);
      }
    });
  };

  return (
    <>
      <Card>
        <Card.Content className="gap-4">
          <SectionHeading
            title={t('tournament.title')}
            hint={t('tournament.hint')}
          />
          <Form
            id="tournament-settings-form"
            onSubmit={handleSave}
            className="flex flex-col gap-4"
          >
            {isFinished ? (
              <p>
                {t('tournament.summary', {
                  count: settings.tours,
                  minutes: settings.tour_minutes,
                })}
              </p>
            ) : (
              <>
                <p className="text-sm">{t('tournament.tours')}</p>
                <ListStepper
                  label={t('tournament.tours')}
                  values={range(MIN_TOURS, MAX_TOURS)}
                  value={settings.tours}
                  onChange={(tours) => setDraft({ ...settings, tours })}
                />
                <div className="flex flex-col gap-1.5">
                  <p className="text-sm">{t('tournament.tourMinutes')}</p>
                  <ValueStepper
                    label={t('tournament.tourMinutes')}
                    value={settings.tour_minutes}
                    min={MIN_TOUR_MINUTES}
                    onChange={(tour_minutes) =>
                      setDraft({ ...settings, tour_minutes })
                    }
                  />
                  <p
                    className={`text-xs text-center ${
                      settings.tour_minutes > MAX_TOUR_MINUTES
                        ? 'text-danger'
                        : 'text-muted'
                    }`}
                  >
                    {t('tournament.tourMinutesRange', {
                      min: MIN_TOUR_MINUTES,
                      max: MAX_TOUR_MINUTES,
                    })}
                  </p>
                </div>
              </>
            )}
            <SettingRadios
              label={t('tournament.timeUpLabel')}
              options={TIME_UP_RULES}
              value={settings.time_up}
              onChange={(time_up) =>
                setDraft({
                  ...settings,
                  time_up,
                  extra_hands:
                    time_up === 'extra_hands'
                      ? Math.max(settings.extra_hands, MIN_EXTRA_HANDS)
                      : 0,
                })
              }
              isDisabled={isFinished || isStarted}
              optionLabel={(option) => t(`tournament.timeUp.${option}`)}
            />
            {isStarted && !isFinished && (
              <p className="text-xs text-muted">
                {t('tournament.timeUpFrozen')}
              </p>
            )}
            {settings.time_up === 'extra_hands' &&
              !isFinished &&
              !isStarted && (
                <>
                  <p className="text-sm">{t('tournament.extraHands')}</p>
                  <ListStepper
                    label={t('tournament.extraHands')}
                    values={range(MIN_EXTRA_HANDS, MAX_EXTRA_HANDS)}
                    value={settings.extra_hands}
                    onChange={(extra_hands) =>
                      setDraft({ ...settings, extra_hands })
                    }
                  />
                </>
              )}
            {/* Any time: it changes only how the rating orders players
                equal in TP and GP. */}
            <SettingRadios
              label={t('tournament.tiesLabel')}
              options={TIES_CHOICES}
              value={tiesChoice(settings)}
              onChange={(choice) => setDraft(withTies(settings, choice))}
              isDisabled={isFinished}
              optionLabel={(option) => t(`tournament.ties.${option}`)}
            />
            {settings.ties === 'break' && (
              <p className="text-xs text-muted">{t('tournament.tiesHint')}</p>
            )}
          </Form>
        </Card.Content>
      </Card>
      {!isFinished && (
        <PendingButton
          type="submit"
          form="tournament-settings-form"
          variant="primary"
          isPending={isPending}
          isDisabled={!isValidTournamentSettings(settings) || !isChanged}
          className="w-3/4 self-center"
        >
          {t('tournament.save')}
        </PendingButton>
      )}
    </>
  );
}
