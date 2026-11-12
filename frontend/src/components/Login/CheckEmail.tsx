import { Trans, useTranslation } from 'react-i18next';
import { Card } from '@heroui/react';
import { LinkButton } from '../LinkButton';
import { useCallback, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { useLocation } from 'react-router-dom';
import { globalFetch, PostRequest } from '../Request';
import { showToast } from '../Toast';
import PendingButton from '../PendingButton';
import EmailInput from './EmailInput';
import { emailProblem } from './validation';

type CheckEmailState = {
  email?: string;
};

export default function CheckEmailPage() {
  const { t } = useTranslation();
  const location = useLocation();
  const stateEmail = (location.state as CheckEmailState | null)?.email ?? '';

  const [email, setEmail] = useState(stateEmail);
  const { isPending: isLoading, run } = useSubmitOnce();

  const handleResend = useCallback(async () => {
    if (!email) {
      return;
    }

    const data = await globalFetch<object>(
      PostRequest('/auth/resend-verification', {
        email,
      }),
      { redirectOn401: false }
    );

    if (!data) {
      return;
    }

    showToast({
      title: t('auth.verificationResent'),
      variant: 'success',
    });
  }, [email, t]);

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
        <Card className="w-full mx-5">
          <Card.Content className="gap-3 py-6">
            <h1 className="text-xl font-semibold">
              {t('auth.checkEmailTitle')}
            </h1>
            <p className="text-base text-muted">
              {/* Neutral on purpose: after a sign-up with a registered email
                the mail offers a password reset, not a verification link, and
                the page must not tell the two cases apart. */}
              {email ? (
                <Trans
                  i18nKey="auth.checkEmailSentTo"
                  values={{ email }}
                  components={{
                    email: (
                      <span className="font-medium text-foreground break-all" />
                    ),
                  }}
                />
              ) : (
                t('auth.checkEmailSent')
              )}
            </p>
            {!stateEmail && (
              <>
                <EmailInput value={email} onChange={setEmail} />
              </>
            )}
            <div className="flex justify-center gap-3">
              <PendingButton
                variant="primary"
                isPending={isLoading}
                isDisabled={emailProblem(email) != null}
                onPress={() => run(handleResend)}
              >
                {t('auth.resendVerification')}
              </PendingButton>
              <LinkButton variant="tertiary" to="/login">
                {t('auth.backToLogin')}
              </LinkButton>
            </div>
          </Card.Content>
        </Card>
      </div>
    </div>
  );
}
