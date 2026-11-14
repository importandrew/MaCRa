import { useTranslation } from 'react-i18next';
import { Card, Chip } from '@heroui/react';
import { useContext } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { AuthContext } from '../../providers/auth/authConstants';
import { useEvent } from '../../providers/event/eventConstants';
import { joinByInvite, useFetchInvite } from '../../providers/event/inviteApi';
import { LinkButton } from '../LinkButton';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';

// JoinEvent is /join/:token, the page an invite link opens: the event's name
// and a button to join it as a player. A visitor signs in first and comes
// back here; a member is told so and goes to the event.
export default function JoinEvent() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const { token = null } = useParams();
  const { isAuthenticated } = useContext(AuthContext);
  const { selectEvent } = useEvent();
  const { invite, isLoading, isNotFound, isError } = useFetchInvite(token);
  const { isPending, run } = useSubmitOnce();

  const join = () =>
    run(async () => {
      if (!token) return;
      const eventId = await joinByInvite(token);
      if (eventId == null) return;
      // Joining to play: the event becomes the selected one.
      await selectEvent(eventId);
      showToast({ title: t('invite.joined'), variant: 'success' });
      navigate(`/event/${eventId}`);
    });

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      <div className="flex flex-col max-w-screen h-full text-center justify-center items-center mx-5">
        <Card className="w-full mx-5">
          <Card.Content className="items-center gap-4 py-8">
            {isLoading && <LoadingSpinner label={t('common.loading')} />}
            {!isLoading && (isNotFound || isError || !invite) && (
              <>
                <p className="text-xl">{t('invite.invalid')}</p>
                <p className="text-sm text-muted">{t('invite.invalidHint')}</p>
                <LinkButton variant="tertiary" to="/">
                  {t('common.goHome')}
                </LinkButton>
              </>
            )}
            {!isLoading && invite && (
              <>
                <p className="text-muted">{t('invite.title')}</p>
                <h1 className="w-full break-words text-2xl font-semibold">
                  {invite.name}
                </h1>
                <Chip size="sm" variant="soft">
                  {t(`eventType.${invite.type}`)}
                </Chip>
                {invite.is_member ? (
                  <>
                    <p>{t('invite.alreadyMember')}</p>
                    <LinkButton
                      variant="primary"
                      to={`/event/${invite.event_id}`}
                    >
                      {t('invite.openEvent')}
                    </LinkButton>
                  </>
                ) : invite.is_full ? (
                  <p className="text-muted">{t('joinRequests.full')}</p>
                ) : isAuthenticated ? (
                  <PendingButton
                    variant="primary"
                    isPending={isPending}
                    onPress={join}
                  >
                    {t('invite.join')}
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
                    {t('invite.signInToJoin')}
                  </PendingButton>
                )}
              </>
            )}
          </Card.Content>
        </Card>
      </div>
    </div>
  );
}
