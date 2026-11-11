import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';

type ForbiddenStateProps = {
  onGoHome?: () => void;
};

export default function ForbiddenState({ onGoHome }: ForbiddenStateProps) {
  const { t } = useTranslation();
  return (
    <div className="flex flex-col h-full w-auto my-auto justify-center items-center gap-4 px-6 text-center">
      <p className="text-3xl">{t('common.forbidden')}</p>
      <p className="text-muted text-lg">{t('common.noAccessToEvent')}</p>
      {onGoHome && (
        <Button type="button" variant="secondary" onPress={onGoHome}>
          {t('common.goHome')}
        </Button>
      )}
    </div>
  );
}
