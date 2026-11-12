import { useTranslation } from 'react-i18next';
import { Button, useOverlayState } from '@heroui/react';
import { DeleteRequest, authFetch, eventGamePath } from '../Request';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { WarningIcon } from '../Icons/WarningIcon';
import Confirm from '../Confirm';
import { showToast } from '../Toast';

export default function AbortMenu({
  gameId,
  eventId,
  canAbort = false,
  onAborted,
}: {
  gameId: number | null;
  // eventId is the event of the shown game: the URL's on the live view, the
  // selected one at the table.
  eventId: number | null;
  canAbort?: boolean;
  // onAborted leaves the screen, as an aborted game has no result: home at the
  // table, the event's live games on the live view.
  onAborted: () => void;
}) {
  const { t } = useTranslation();
  const { isPending: isAborting, run } = useSubmitOnce();
  const { isOpen, open: onOpen, close: onClose } = useOverlayState();

  if (!canAbort) {
    return null;
  }

  const handleAbort = async () => {
    if (!gameId || !eventId) return;

    const abortRequest = DeleteRequest(
      eventGamePath(eventId, gameId, '/abort')
    );
    const data = await authFetch<null>(abortRequest);

    if (!data) return;

    showToast({
      title: t('game.aborted'),
      variant: 'warning',
    });

    onAborted();
  };

  return (
    <>
      <Button
        variant="danger"
        onPress={onOpen}
        isDisabled={!gameId || !eventId || isAborting}
        className="w-1/2"
      >
        <WarningIcon className="text-white" />
        {isAborting ? t('game.aborting') : t('game.abort')}
      </Button>
      <Confirm
        onClose={onClose}
        onConfirm={() => {
          onClose();
          void run(handleAbort);
        }}
        isOpen={isOpen}
        headerText={t('game.abortConfirmTitle')}
        cancelText={t('common.cancel')}
        confirmText={t('game.abort')}
        prompt={t('game.abortConfirmPrompt')}
        confirmVariant="danger"
      />
    </>
  );
}
