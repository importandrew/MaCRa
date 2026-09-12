import {
  Button,
  Card,
  CardBody,
  Form,
  Input,
  Tab,
  Tabs,
  Textarea,
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

const VISIBILITY_OPTIONS: { value: EventVisibility; label: string }[] = [
  { value: 'private', label: 'Private' },
  { value: 'public', label: 'Public' },
];

const TYPE_OPTIONS: { value: 'regular' | 'tournament'; label: string, isDisabled: boolean }[] = [
  { value: 'regular', label: 'Regular',  isDisabled: false},
  { value: 'tournament', label: 'Tournament', isDisabled: true},
];

type VisibilityTabsProps<T extends string> = {
  label: string;
  ariaLabel: string;
  options: { value: T; label: string }[];
  selectedKey: T;
  onSelectionChange: (key: T) => void;
  isOptionDisabled?: (value: T) => boolean;
};

function VisibilityTabs<T extends string>({
  label,
  ariaLabel,
  options,
  selectedKey,
  onSelectionChange,
  isOptionDisabled,
}: VisibilityTabsProps<T>) {
  return (
    <div className="flex flex-col gap-1.5">
      <p className="text-sm">{label}</p>
      <Tabs
        aria-label={ariaLabel}
        fullWidth
        selectedKey={selectedKey}
        onSelectionChange={(key) => onSelectionChange(key as T)}
      >
        {options.map((opt) => (
          <Tab
            key={opt.value}
            title={opt.label}
            isDisabled={isOptionDisabled?.(opt.value)}
          />
        ))}
      </Tabs>
    </div>
  );
}

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
      color: 'success',
    });
    navigate(`/event/${id}`);
  };

  return (
    <div className="flex flex-col w-full max-w-[600px] mx-auto px-3 py-6 gap-4 min-w-0">
      <h1 className="text-2xl font-semibold">Create event</h1>

      <Card>
        <CardBody className="gap-2">
          <p className="text-default-500">General</p>
          <Form id="create-event-form" onSubmit={onSubmit}>
            <Input
              label="Name"
              value={name}
              onValueChange={setName}
              isRequired
              variant="bordered"
            />
            <Textarea
              label="Description"
              value={description}
              onValueChange={setDescription}
              variant="bordered"
            />
            <div className="w-full flex flex-col gap-2">
              <p className="text-default-500">Type</p>
              <Tabs
                aria-label="Event type"
                fullWidth
                selectedKey={type}
                onSelectionChange={(key) =>
                  setType(key as 'regular' | 'tournament')
                }
              >
                {TYPE_OPTIONS.map((opt) => (
                  <Tab key={opt.value} title={opt.label} isDisabled={opt.isDisabled} />
                ))}
              </Tabs>
              <p className="text-default-500">Visibility</p>
              <VisibilityTabs
                label="General"
                ariaLabel="Event visibility"
                options={VISIBILITY_OPTIONS}
                selectedKey={visibility}
                onSelectionChange={applyVisibilityPreset}
              />
              <div className="flex flex-col gap-5">
                <VisibilityTabs
                  label="Ratings"
                  ariaLabel="Ratings visibility"
                  options={CONTENT_REACH_OPTIONS}
                  selectedKey={ratingsReach}
                  onSelectionChange={setRatingsReach}
                  isOptionDisabled={isReachOptionDisabled}
                />
                <VisibilityTabs
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
        </CardBody>
      </Card>

      <div className="flex gap-2 w-3/4 self-center">
        <Button
          type="button"
          variant="flat"
          className="flex-1"
          onPress={() => navigate('/events')}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          form="create-event-form"
          color="primary"
          isLoading={isSubmitting}
          isDisabled={!name.trim()}
          className="flex-1"
        >
          Create
        </Button>
      </div>
    </div>
  );
}
