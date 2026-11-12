import { useTranslation } from 'react-i18next';
import { useOverlayState } from '@heroui/react';
import PendingButton from '../PendingButton';
import { PostRequest, authFetch, eventGamePath } from '../Request';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { WarningIcon } from '../Icons/WarningIcon';
import Confirm from '../Confirm';
import { showToast } from '../Toast';

export default function FinishMenu({
  gameId,
  eventId,
  canFinish = false,
  onFinished,
}: {
  gameId: number | null;
  // eventId is the event of the shown game: the URL's on the live view, the
  // selected one at the table.
  eventId: number | null;
  canFinish?: boolean;
  // onFinished shows the result: /last at the table, where the finished game
  // is the user's last one; the game's own page on the live view, where the
  // user may be a referee who didn't play it.
  onFinished: () => void;
}) {
  const { t } = useTranslation();
  const { isPending: isFinishing, run } = useSubmitOnce();
  const { isOpen, open: onOpen, close: onClose } = useOverlayState();

  const handleFinish = async () => {
    if (!gameId || !eventId) return;

    const finishRequest = PostRequest(
      eventGamePath(eventId, gameId, '/finish')
    );
    const data = await authFetch<null>(finishRequest);

    if (!data) return;

    showToast({
      title: t('game.finished'),
      variant: 'success',
    });

    onFinished();
  };

  if (!canFinish) {
    return null;
  }

  return (
    <>
      <PendingButton
        variant="primary"
        onPress={onOpen}
        isPending={isFinishing}
        isDisabled={!gameId || !eventId}
        className="w-1/2"
        icon={<WarningIcon className="text-white" />}
      >
        {t('game.finish')}
      </PendingButton>
      <Confirm
        onClose={onClose}
        onConfirm={() => {
          onClose();
          void run(handleFinish);
        }}
        isOpen={isOpen}
        headerText={t('game.finishConfirmTitle')}
        cancelText={t('common.cancel')}
        confirmText={t('game.finish')}
        prompt={t('game.finishConfirmPrompt')}
        confirmVariant="primary"
      />
    </>
  );
}
