import type { ContentReach, EventVisibility } from '../../models/Event';

// EventSettings are the fields an event is created with and managed by.
export type EventSettings = {
  name: string;
  description: string;
  visibility: EventVisibility;
  ratingsReach: ContentReach;
  gamesReach: ContentReach;
};

// The limits mirror the backend's (models.MaxEventNameLength, ...).
export const MAX_EVENT_NAME_LENGTH = 100;
export const MAX_EVENT_DESCRIPTION_LENGTH = 1000;

export type NameProblem = 'required' | 'tooLong' | null;

export function eventNameProblem(name: string): NameProblem {
  const trimmed = name.trim();
  if (!trimmed) return 'required';
  if ([...trimmed].length > MAX_EVENT_NAME_LENGTH) return 'tooLong';
  return null;
}

export function isDescriptionTooLong(description: string): boolean {
  return [...description.trim()].length > MAX_EVENT_DESCRIPTION_LENGTH;
}

// isValidSettings says whether the form may be submitted.
export function isValidSettings(settings: EventSettings): boolean {
  return (
    eventNameProblem(settings.name) == null &&
    !isDescriptionTooLong(settings.description)
  );
}

// withVisibility changes the visibility and a reach only where it becomes
// invalid: a public reach in a private event turns into members. A reach the
// owner set otherwise is kept.
export function withVisibility(
  settings: EventSettings,
  visibility: EventVisibility
): EventSettings {
  const fit = (reach: ContentReach): ContentReach =>
    visibility === 'private' && reach === 'public' ? 'members' : reach;
  return {
    ...settings,
    visibility,
    ratingsReach: fit(settings.ratingsReach),
    gamesReach: fit(settings.gamesReach),
  };
}
