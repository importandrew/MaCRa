import { useTranslation } from 'react-i18next';
import { Button, Card, Form } from '@heroui/react';
import { FormEvent, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { useNavigate } from 'react-router-dom';
import { createEvent } from '../../providers/event/eventApi';
import { showToast } from '../Toast';
import SegmentedPicker from '../SegmentedPicker';
import PendingButton from '../PendingButton';
import EventSettingsFields from './EventSettingsFields';
import { type EventSettings, isValidSettings } from './eventSettings';

export default function EventCreate() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [type, setType] = useState<'regular' | 'tournament'>('regular');
  // A new event is private by default, its rating and games for members.
  const [settings, setSettings] = useState<EventSettings>({
    name: '',
    description: '',
    visibility: 'private',
    ratingsReach: 'members',
    gamesReach: 'members',
  });
  const { isPending: isSubmitting, run } = useSubmitOnce();

  const onSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    void run(submit);
  };

  const submit = async () => {
    const id = await createEvent({
      type,
      name: settings.name.trim(),
      description: settings.description.trim(),
      visibility: settings.visibility,
      ratings_reach: settings.ratingsReach,
      games_reach: settings.gamesReach,
    });

    if (!id) return;

    showToast({
      title: t('events.created'),
      variant: 'success',
    });
    navigate(`/event/${id}`);
  };

  return (
    <div className="flex flex-col w-full max-w-[600px] mx-auto px-3 py-6 gap-4 min-w-0">
      <h1 className="text-2xl font-semibold">{t('events.createTitle')}</h1>

      <Card>
        <Card.Content className="gap-2">
          <p className="text-muted">{t('events.general')}</p>
          <Form
            id="create-event-form"
            onSubmit={onSubmit}
            className="flex flex-col gap-3"
          >
            <EventSettingsFields
              value={settings}
              onChange={setSettings}
              afterDescription={
                <>
                  <p className="text-muted">{t('events.type')}</p>
                  <SegmentedPicker
                    ariaLabel={t('events.type')}
                    options={[
                      { value: 'regular', label: t('eventType.regular') },
                      {
                        value: 'tournament',
                        label: t('eventType.tournament'),
                        isDisabled: true,
                      },
                    ]}
                    selectedKey={type}
                    onSelectionChange={setType}
                  />
                </>
              }
            />
          </Form>
        </Card.Content>
      </Card>

      <div className="flex gap-2 w-3/4 self-center">
        <Button
          type="button"
          variant="tertiary"
          className="flex-1"
          onPress={() => navigate('/events')}
        >
          {t('common.cancel')}
        </Button>
        <PendingButton
          type="submit"
          form="create-event-form"
          variant="primary"
          isPending={isSubmitting}
          isDisabled={!isValidSettings(settings)}
          className="flex-1"
        >
          {t('events.create')}
        </PendingButton>
      </div>
    </div>
  );
}
