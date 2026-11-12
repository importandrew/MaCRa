import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';
import { useNavigate } from 'react-router-dom';

// CrashState is what a page shows instead of itself after a render error:
// what happened, the error text (to quote when reporting it) and two ways out.
export default function CrashState({ error }: { error: Error }) {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
    <div className="flex flex-col grow w-full my-auto justify-center items-center gap-4 px-6 text-center">
      <p className="text-2xl">{t('crash.title')}</p>
      <p className="text-muted">{t('crash.hint')}</p>
      <p className="text-sm text-muted break-words max-w-full">
        {error.message}
      </p>
      <div className="flex gap-3">
        <Button
          type="button"
          variant="primary"
          onPress={() => window.location.reload()}
        >
          {t('crash.reload')}
        </Button>
        <Button type="button" variant="secondary" onPress={() => navigate('/')}>
          {t('common.goHome')}
        </Button>
      </div>
    </div>
  );
}
