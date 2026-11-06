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
import { useNavigate } from 'react-router-dom';
import {
  CONTENT_REACH_OPTIONS,
  ContentReach,
  EventVisibility,
} from '../../models/Event';
import { createEvent } from '../../providers/event/eventApi';
import { showToast } from '../Toast';
import SegmentedPicker from '../SegmentedPicker';
import PendingButton from '../PendingButton';

const VISIBILITY_OPTIONS: { value: EventVisibility; label: string }[] = [
  { value: 'private', label: 'Private' },
  { value: 'public', label: 'Public' },
];

const TYPE_OPTIONS: {
  value: 'regular' | 'tournament';
  label: string;
  isDisabled: boolean;
}[] = [
  { value: 'regular', label: 'Regular', isDisabled: false },
  { value: 'tournament', label: 'Tournament', isDisabled: true },
];

export default function EventCreate() {
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [type, setType] = useState<'regular' | 'tournament'>('regular');
  const [visibility, setVisibility] = useState<EventVisibility>('private');
  const [ratingsReach, setRatingsReach] = useState<ContentReach>('members');
  const [gamesReach, setGamesReach] = useState<ContentReach>('members');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const publicReachDisabled = visibility === 'private';
  const isReachOptionDisabled = (value: ContentReach) =>
    value === 'public' && publicReachDisabled;

  const applyVisibilityPreset = (value: EventVisibility) => {
    setVisibility(value);
    const reach: ContentReach = value === 'public' ? 'public' : 'members';
    setRatingsReach(reach);
    setGamesReach(reach);
  };

  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsSubmitting(true);

    const id = await createEvent({
      type,
      name: name.trim(),
      description: description.trim() || undefined,
      visibility,
      ratings_reach: ratingsReach,
      games_reach: gamesReach,
    });

    setIsSubmitting(false);

    if (!id) return;

    showToast({
      title: 'Event created',
      variant: 'success',
    });
    navigate(`/event/${id}`);
  };

  return (
    <div className="flex flex-col w-full max-w-[600px] mx-auto px-3 py-6 gap-4 min-w-0">
      <h1 className="text-2xl font-semibold">Create event</h1>

      <Card>
        <Card.Content className="gap-2">
          <p className="text-muted">General</p>
          <Form
            id="create-event-form"
            onSubmit={onSubmit}
            className="flex flex-col gap-3"
          >
            <TextField fullWidth isRequired value={name} onChange={setName}>
              <Label>Name</Label>
              <Input />
            </TextField>
            <TextField fullWidth value={description} onChange={setDescription}>
              <Label>Description</Label>
              <TextArea />
            </TextField>
            <div className="w-full flex flex-col gap-2">
              <p className="text-muted">Type</p>
              <SegmentedPicker
                ariaLabel="Event type"
                options={TYPE_OPTIONS}
                selectedKey={type}
                onSelectionChange={setType}
              />
              <p className="text-muted">Visibility</p>
              <SegmentedPicker
                label="General"
                ariaLabel="Event visibility"
                options={VISIBILITY_OPTIONS}
                selectedKey={visibility}
                onSelectionChange={applyVisibilityPreset}
              />
              <div className="flex flex-col gap-5">
                <SegmentedPicker
                  label="Ratings"
                  ariaLabel="Ratings visibility"
                  options={CONTENT_REACH_OPTIONS}
                  selectedKey={ratingsReach}
                  onSelectionChange={setRatingsReach}
                  isOptionDisabled={isReachOptionDisabled}
                />
                <SegmentedPicker
                  label="Live games & results"
                  ariaLabel="Live games and results visibility"
                  options={CONTENT_REACH_OPTIONS}
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
          Cancel
        </Button>
        <PendingButton
          type="submit"
          form="create-event-form"
          variant="primary"
          isPending={isSubmitting}
          isDisabled={!name.trim()}
          className="flex-1"
        >
          Create
        </PendingButton>
      </div>
    </div>
  );
}
