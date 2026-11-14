import { useTranslation } from 'react-i18next';
import { Disclosure } from '@heroui/react';
import { type ReactNode, useState } from 'react';

// FinishedEvents is the collapsed "Finished (N)" part of a list of events:
// finished events stay reachable (the last game, history) without pushing
// the active ones down. While a search is on (isSearching) it is open, so
// matching finished events show at once.
export default function FinishedEvents({
  count,
  isSearching,
  children,
}: {
  count: number;
  isSearching: boolean;
  children: ReactNode;
}) {
  const { t } = useTranslation();
  const [isExpanded, setIsExpanded] = useState(false);
  if (count === 0) return null;

  return (
    <Disclosure
      isExpanded={isSearching || isExpanded}
      onExpandedChange={setIsExpanded}
    >
      <Disclosure.Heading>
        <Disclosure.Trigger className="flex items-center gap-1 py-2 text-muted">
          {t('events.finishedCount', { n: count })}
          <Disclosure.Indicator />
        </Disclosure.Trigger>
      </Disclosure.Heading>
      <Disclosure.Content>
        {/* Not Disclosure.Body: its p-2 wrapper made these cards narrower
            than the active ones above. */}
        <div className="flex flex-col gap-3 pt-1">{children}</div>
      </Disclosure.Content>
    </Disclosure>
  );
}
