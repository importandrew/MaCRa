import { useTranslation } from 'react-i18next';
import { Button, useOverlayState } from '@heroui/react';
import { QRCodeSVG } from 'qrcode.react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import {
  inviteUrl,
  newInviteToken,
  turnOffInvite,
  useFetchInviteToken,
} from '../../providers/event/inviteApi';
import Confirm from '../Confirm';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';
import { copyText } from '../../utils/copyText';

// InviteLink is the invite link in Manage → Members: whoever opens it while
// signed in becomes a player. The owner turns it on, shares it (copy, the
// phone's share sheet, or the QR code to scan), replaces it if it leaked
// (the old one stops working) or turns it off.
export default function InviteLink({ eventId }: { eventId: number }) {
  const { t } = useTranslation();
  const { token, isLoading, isError, mutate } = useFetchInviteToken(eventId);
  const { isPending, run } = useSubmitOnce();
  const replaceConfirm = useOverlayState();
  const offConfirm = useOverlayState();

  if (isLoading)
    return <LoadingSpinner size="sm" label={t('common.loading')} />;
  if (isError) {
    return <p className="text-sm text-danger">{t('invite.loadFailed')}</p>;
  }

  const create = () =>
    run(async () => {
      replaceConfirm.close();
      const next = await newInviteToken(eventId);
      if (next) {
        mutate({ token: next });
        showToast({ title: t('invite.created'), variant: 'success' });
      }
    });

  const turnOff = () =>
    run(async () => {
      offConfirm.close();
      if (await turnOffInvite(eventId)) {
        mutate({ token: null });
        showToast({ title: t('invite.turnedOff'), variant: 'success' });
      }
    });

  if (!token) {
    return (
      <div className="flex flex-col gap-3">
        <p className="text-sm text-muted">{t('invite.offHint')}</p>
        <PendingButton
          variant="primary"
          isPending={isPending}
          onPress={create}
          className="w-3/4 self-center"
        >
          {t('invite.create')}
        </PendingButton>
      </div>
    );
  }

  const url = inviteUrl(token);
  const copy = async () => {
    if (await copyText(url)) {
      showToast({ title: t('invite.copied'), variant: 'success' });
    } else {
      showToast({ title: t('invite.copyFailed'), variant: 'warning' });
    }
  };
  const canShare = typeof navigator.share === 'function';
  const share = () =>
    navigator.share({ title: t('invite.shareTitle'), url }).catch(() => {
      // Closing the share sheet is not an error.
    });

  return (
    <div className="flex flex-col items-center gap-3">
      <p className="text-sm text-muted self-start">{t('invite.onHint')}</p>
      {/* White around the code: scanners need a quiet zone, also in dark
          mode. */}
      <div className="rounded-xl bg-white p-3">
        <QRCodeSVG value={url} size={192} aria-label={t('invite.qr')} />
      </div>
      <p className="w-full break-all text-center text-sm">{url}</p>
      <div className="flex flex-wrap justify-center gap-2">
        <Button variant="secondary" onPress={copy}>
          {t('invite.copy')}
        </Button>
        {canShare && (
          <Button variant="secondary" onPress={share}>
            {t('invite.share')}
          </Button>
        )}
      </div>
      <div className="flex flex-wrap justify-center gap-2">
        <PendingButton
          variant="tertiary"
          isPending={isPending}
          onPress={replaceConfirm.open}
        >
          {t('invite.replace')}
        </PendingButton>
        <Button variant="danger-soft" onPress={offConfirm.open}>
          {t('invite.turnOff')}
        </Button>
      </div>
      <Confirm
        isOpen={replaceConfirm.isOpen}
        onClose={replaceConfirm.close}
        onConfirm={create}
        headerText={t('invite.replaceTitle')}
        prompt={t('invite.replacePrompt')}
        cancelText={t('common.cancel')}
        confirmText={t('invite.replace')}
        confirmVariant="primary"
      />
      <Confirm
        isOpen={offConfirm.isOpen}
        onClose={offConfirm.close}
        onConfirm={turnOff}
        headerText={t('invite.turnOffTitle')}
        prompt={t('invite.turnOffPrompt')}
        cancelText={t('common.cancel')}
        confirmText={t('invite.turnOff')}
        confirmVariant="danger"
      />
    </div>
  );
}
