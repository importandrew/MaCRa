import { useTranslation } from 'react-i18next';
import { useContext } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { AuthContext } from '../../providers/auth/authConstants';
import { useEventLayout } from '../../providers/event/eventLayoutContext';
import {
  requestToJoin,
  withdrawJoinRequest,
} from '../../providers/event/joinRequestApi';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';

// JoinRequestBox is the overview's offer to a non-member of a public event
// that takes requests to join (EventMe.join_request): ask to join (a visitor
// signs in first), the pending request with a way to withdraw it, a declined
// request, or no places left. Nothing for anyone else.
export default function JoinRequestBox() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated } = useContext(AuthContext);
  const { eventId, me, refetch } = useEventLayout();
  const { isPending, run } = useSubmitOnce();
  const state = me?.join_request;

  if (!state) return null;

  const ask = () =>
    run(async () => {
      if (await requestToJoin(eventId)) {
        showToast({ title: t('joinRequests.sent'), variant: 'success' });
      }
      // Also after an error: the state may have changed meanwhile (the
      // event got full or stopped taking requests).
      await refetch();
    });

  const withdraw = () =>
    run(async () => {
      if (await withdrawJoinRequest(eventId)) {
        showToast({ title: t('joinRequests.withdrawn'), variant: 'success' });
      }
      await refetch();
    });

  return (
    <div className="flex flex-col items-center gap-2 text-center">
      {state === 'available' &&
        (isAuthenticated ? (
          <PendingButton variant="primary" isPending={isPending} onPress={ask}>
            {t('joinRequests.ask')}
          </PendingButton>
        ) : (
          <PendingButton
            variant="primary"
            onPress={() =>
              navigate('/login', {
                state: { from: { pathname: location.pathname } },
              })
            }
          >
            {t('joinRequests.signInToAsk')}
          </PendingButton>
        ))}
      {state === 'pending' && (
        <>
          <p>{t('joinRequests.pending')}</p>
          <PendingButton
            variant="tertiary"
            isPending={isPending}
            onPress={withdraw}
          >
            {t('joinRequests.withdraw')}
          </PendingButton>
        </>
      )}
      {state === 'declined' && (
        <p className="text-muted">{t('joinRequests.declined')}</p>
      )}
      {state === 'full' && (
        <p className="text-muted">{t('joinRequests.full')}</p>
      )}
    </div>
  );
}
