import { useTranslation } from 'react-i18next';
import { SearchField } from '@heroui/react';

// EventSearchField filters a list of events by name (eventListFilters.ts).
export default function EventSearchField({
  value,
  onChange,
}: {
  value: string;
  onChange: (value: string) => void;
}) {
  const { t } = useTranslation();
  return (
    <SearchField
      fullWidth
      aria-label={t('events.search')}
      value={value}
      onChange={onChange}
    >
      <SearchField.Group>
        <SearchField.SearchIcon />
        <SearchField.Input placeholder={t('events.search')} />
        <SearchField.ClearButton aria-label={t('events.clearSearch')} />
      </SearchField.Group>
    </SearchField>
  );
}
