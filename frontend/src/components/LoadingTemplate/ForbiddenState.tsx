import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';
import type { PageSubject } from './NotFoundState';

type ForbiddenStateProps = {
  // subject is what the page shows, for the text: "this event" or "this game".
  subject?: PageSubject;
  onGoHome?: () => void;
};

export default function ForbiddenState({
  subject = 'event',
  onGoHome,
}: ForbiddenStateProps) {
  const { t } = useTranslation();
  return (
    <div className="flex flex-col h-full w-auto my-auto justify-center items-center gap-4 px-6 text-center">
      <p className="text-3xl">{t('common.forbidden')}</p>
      <p className="text-muted text-lg">
        {subject === 'game'
          ? t('common.noAccessToGame')
          : t('common.noAccessToEvent')}
      </p>
      {onGoHome && (
        <Button type="button" variant="secondary" onPress={onGoHome}>
          {t('common.goHome')}
        </Button>
      )}
    </div>
  );
}
