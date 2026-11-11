import { useTranslation } from 'react-i18next';
import { useState } from 'react';
import { RefreshIcon } from './Icons/RefreshIcon';
import PendingButton from './PendingButton';

// RefreshButton is the round refresh button on game screens. It shows a
// spinner instead of the icon until onRefresh finishes.
export default function RefreshButton({
  onRefresh,
  className = '',
}: {
  onRefresh: () => Promise<unknown> | void;
  className?: string;
}) {
  const { t } = useTranslation();
  const [isRefreshing, setIsRefreshing] = useState(false);

  const refresh = async () => {
    setIsRefreshing(true);
    try {
      await onRefresh();
    } finally {
      setIsRefreshing(false);
    }
  };

  return (
    <PendingButton
      variant="ghost"
      size="lg"
      isIconOnly
      aria-label={t('common.refresh')}
      isPending={isRefreshing}
      onPress={refresh}
      className={`rounded-full ${className}`}
      icon={<RefreshIcon className="text-muted" />}
    />
  );
}
