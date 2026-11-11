import { useTranslation } from 'react-i18next';
import {
  Button,
  Card,
  Form,
  Input,
  Label,
  TextArea,
  TextField,
} from '@heroui/react';
import { FormEvent, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { useNavigate } from 'react-router-dom';
import {
  CONTENT_REACHES,
  ContentReach,
  EventVisibility,
} from '../../models/Event';
import { createEvent } from '../../providers/event/eventApi';
import { showToast } from '../Toast';
import SegmentedPicker from '../SegmentedPicker';
import PendingButton from '../PendingButton';

const VISIBILITIES: EventVisibility[] = ['private', 'public'];

export default function EventCreate() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [type, setType] = useState<'regular' | 'tournament'>('regular');
  const [visibility, setVisibility] = useState<EventVisibility>('private');
  const [ratingsReach, setRatingsReach] = useState<ContentReach>('members');
  const [gamesReach, setGamesReach] = useState<ContentReach>('members');
  const { isPending: isSubmitting, run } = useSubmitOnce();

  const publicReachDisabled = visibility === 'private';
  const isReachOptionDisabled = (value: ContentReach) =>
    value === 'public' && publicReachDisabled;

  const applyVisibilityPreset = (value: EventVisibility) => {
    setVisibility(value);
    const reach: ContentReach = value === 'public' ? 'public' : 'members';
    setRatingsReach(reach);
    setGamesReach(reach);
  };

  const onSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    void run(submit);
  };

  const submit = async () => {
    const id = await createEvent({
      type,
      name: name.trim(),
      description: description.trim() || undefined,
      visibility,
      ratings_reach: ratingsReach,
      games_reach: gamesReach,
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
            <TextField fullWidth isRequired value={name} onChange={setName}>
              <Label>{t('events.name')}</Label>
              <Input />
            </TextField>
            <TextField fullWidth value={description} onChange={setDescription}>
              <Label>{t('events.description')}</Label>
              <TextArea />
            </TextField>
            <div className="w-full flex flex-col gap-2">
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
              <p className="text-muted">{t('events.visibility')}</p>
              <SegmentedPicker
                label={t('events.general')}
                options={VISIBILITIES.map((value) => ({
                  value,
                  label: t(`visibility.${value}`),
                }))}
                selectedKey={visibility}
                onSelectionChange={applyVisibilityPreset}
              />
              <div className="flex flex-col gap-5">
                <SegmentedPicker
                  label={t('events.ratingsReach')}
                  options={CONTENT_REACHES.map((value) => ({
                    value,
                    label: t(`reach.${value}`),
                  }))}
                  selectedKey={ratingsReach}
                  onSelectionChange={setRatingsReach}
                  isOptionDisabled={isReachOptionDisabled}
                />
                <SegmentedPicker
                  label={t('events.gamesReach')}
                  options={CONTENT_REACHES.map((value) => ({
                    value,
                    label: t(`reach.${value}`),
                  }))}
                  selectedKey={gamesReach}
                  onSelectionChange={setGamesReach}
                  isOptionDisabled={isReachOptionDisabled}
                />
              </div>
            </div>
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
          isDisabled={!name.trim()}
          className="flex-1"
        >
          {t('events.create')}
        </PendingButton>
      </div>
    </div>
  );
}
