import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { FieldError, Input, Label, TextArea, TextField } from '@heroui/react';
import type { EventVisibility } from '../../models/Event';
import SegmentedPicker from '../SegmentedPicker';
import {
  type EventSettings,
  eventNameProblem,
  isDescriptionTooLong,
  MAX_EVENT_DESCRIPTION_LENGTH,
  MAX_EVENT_NAME_LENGTH,
} from './eventSettings';

const VISIBILITIES: EventVisibility[] = ['private', 'public'];

// EventSettingsFields are the event's name, description and who sees it,
// the same when creating an event and managing it. In a finished event
// (isFinished) only who sees it can change.
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
        {/* One setting decides who sees the event (area 1): everything
            but managing it, to everyone or to its members. */}
        <SegmentedPicker
          label={t('events.whoSees')}
          options={VISIBILITIES.map((visibility) => ({
            value: visibility,
            label: t(`whoSees.${visibility}`),
          }))}
          selectedKey={value.visibility}
          onSelectionChange={(visibility) => onChange({ ...value, visibility })}
        />
        <p className="text-sm text-muted">
          {t(`whoSees.hint.${value.visibility}`)}
        </p>
      </div>
    </>
  );
}
