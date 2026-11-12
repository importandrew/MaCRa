import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';
import ForbiddenState from './ForbiddenState';
import NotFoundState, { type NotFound } from './NotFoundState';
import LoadingSpinner from './LoadingSpinner';

interface LoadingTemplateProps {
  isLoading: boolean;
  isError: boolean;
  isForbidden?: boolean;
  // notFound shows NotFoundState: the page's event or game doesn't exist.
  notFound?: NotFound | null;
  onRefresh: () => void;
  onGoHome?: () => void;
  children: React.ReactNode;
}

const LoadingTemplate = (props: LoadingTemplateProps) => {
  const { t } = useTranslation();
  const {
    isLoading,
    isError,
    isForbidden,
    notFound,
    onRefresh,
    onGoHome,
    children,
  } = props;

  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-auto my-auto justify-center items-center">
        <LoadingSpinner label={t('common.loading')} />
      </div>
    );
  }

  if (notFound) {
    return <NotFoundState what={notFound} onGoHome={onGoHome} />;
  }

  if (isForbidden) {
    return <ForbiddenState onGoHome={onGoHome} />;
  }

  if (isError) {
    return (
      <div className="flex flex-col h-full w-auto my-auto justify-center items-center">
        <p className="mb-4 text-center text-3xl">{t('common.error')}</p>
        <Button type="button" onPress={onRefresh}>
          {t('common.tryAgain')}
        </Button>
      </div>
    );
  }

  return <>{children}</>;
};

export default LoadingTemplate;
