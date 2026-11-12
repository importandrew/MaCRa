import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';

// PageSubject is what a page is about: an event or a game. For NotFoundState
// it is what is missing (a bad id in the URL counts as missing).
export type PageSubject = 'event' | 'game';
export type NotFound = PageSubject;

// NotFoundState is a page whose event or game doesn't exist. No retry: asking
// again won't make it appear.
export default function NotFoundState({
  what,
  onGoHome,
}: {
  what: NotFound;
  onGoHome?: () => void;
}) {
  const { t } = useTranslation();
  return (
    <div className="flex flex-col h-full w-auto my-auto justify-center items-center gap-4 px-6 text-center">
      <p className="text-3xl">
        {what === 'event'
          ? t('apiErrors.event_not_found')
          : t('apiErrors.game_not_found')}
      </p>
      {onGoHome && (
        <Button type="button" variant="secondary" onPress={onGoHome}>
          {t('common.goHome')}
        </Button>
      )}
    </div>
  );
}
