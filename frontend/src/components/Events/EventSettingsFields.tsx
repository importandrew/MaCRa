import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { FieldError, Input, Label, TextArea, TextField } from '@heroui/react';
import {
  CONTENT_REACHES,
  type ContentReach,
  type EventVisibility,
} from '../../models/Event';
import SegmentedPicker from '../SegmentedPicker';
import {
  type EventSettings,
  eventNameProblem,
  isDescriptionTooLong,
  MAX_EVENT_DESCRIPTION_LENGTH,
  MAX_EVENT_NAME_LENGTH,
  withVisibility,
} from './eventSettings';

const VISIBILITIES: EventVisibility[] = ['private', 'public'];

// EventSettingsFields are the event's name, description, visibility and
// reaches, the same when creating an event and managing it. In a finished
// event (isFinished) only the visibility can change.
export default function EventSettingsFields({
  value,
  onChange,
  isFinished = false,
  afterDescription,
}: {
  value: EventSettings;
  onChange: (value: EventSettings) => void;
  isFinished?: boolean;
  // afterDescription is shown between the description and the visibility
  // (the event type when creating one).
  afterDescription?: ReactNode;
}) {
  const { t } = useTranslation();
  const reachOptions = CONTENT_REACHES.map((reach) => ({
    value: reach,
    label: t(`reach.${reach}`),
  }));
  // A public reach means "everyone" only in a public event.
  const isReachOptionDisabled = (reach: ContentReach) =>
    isFinished || (reach === 'public' && value.visibility === 'private');

  return (
    <>
      <TextField
        fullWidth
        isRequired
        isDisabled={isFinished}
        value={value.name}
        onChange={(name) => onChange({ ...value, name })}
        validate={(name: string) => {
          switch (eventNameProblem(name)) {
            case 'required':
              return t('events.nameRequired');
            case 'tooLong':
              return t('events.nameTooLong', { max: MAX_EVENT_NAME_LENGTH });
            default:
              return null;
          }
        }}
      >
        <Label>{t('events.name')}</Label>
        <Input />
        <FieldError />
      </TextField>
      <TextField
        fullWidth
        isDisabled={isFinished}
        value={value.description}
        onChange={(description) => onChange({ ...value, description })}
        validate={(description: string) =>
          isDescriptionTooLong(description)
            ? t('events.descriptionTooLong', {
                max: MAX_EVENT_DESCRIPTION_LENGTH,
              })
            : null
        }
      >
        <Label>{t('events.description')}</Label>
        <TextArea />
        <FieldError />
      </TextField>
      <div className="w-full flex flex-col gap-2">
        {afterDescription}
        <p className="text-muted">{t('events.visibility')}</p>
        <SegmentedPicker
          label={t('events.general')}
          options={VISIBILITIES.map((visibility) => ({
            value: visibility,
            label: t(`visibility.${visibility}`),
          }))}
          selectedKey={value.visibility}
          onSelectionChange={(visibility) =>
            // A finished event keeps its reaches (only the visibility may
            // change); a public reach in a private event means members anyway.
            onChange(
              isFinished
                ? { ...value, visibility }
                : withVisibility(value, visibility)
            )
          }
        />
        <div className="flex flex-col gap-5">
          <SegmentedPicker
            label={t('events.ratingsReach')}
            options={reachOptions}
            selectedKey={value.ratingsReach}
            onSelectionChange={(ratingsReach) =>
              onChange({ ...value, ratingsReach })
            }
            isOptionDisabled={isReachOptionDisabled}
          />
          <SegmentedPicker
            label={t('events.gamesReach')}
            options={reachOptions}
            selectedKey={value.gamesReach}
            onSelectionChange={(gamesReach) =>
              onChange({ ...value, gamesReach })
            }
            isOptionDisabled={isReachOptionDisabled}
          />
        </div>
      </div>
    </>
  );
}
