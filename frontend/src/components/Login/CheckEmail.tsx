import { Button, Card } from '@heroui/react';
import { LinkButton } from '../LinkButton';
import { useCallback, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { useLocation } from 'react-router-dom';
import { globalFetch, PostRequest } from '../Request';
import { showToast } from '../Toast';
import EmailInput from './EmailInput';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';

type CheckEmailState = {
  email?: string;
};

export default function CheckEmailPage() {
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
        origin: window.location.origin,
      }),
      { redirectOn401: false }
    );

    if (!data) {
      return;
    }

    showToast({
      title: 'Verification email sent if the account is unverified',
      variant: 'success',
    });
  }, [email]);

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      {isLoading ? (
        <LoadingSpinner label="Loading..." />
      ) : (
        <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
          <Card className="w-full mx-5">
            <Card.Content className="gap-3 py-6">
              <h1 className="text-xl font-semibold">Check your email</h1>
              <p className="text-base text-muted">
                We sent a verification link
                {email ? (
                  <>
                    {' '}
                    to{' '}
                    <span className="font-medium text-foreground">{email}</span>
                  </>
                ) : null}
                . Open it to activate your account before signing in.
              </p>
              {!stateEmail && (
                <>
                  <EmailInput onChange={setEmail} />
                </>
              )}
              <div className="flex justify-center gap-3">
                <Button
                  variant="primary"
                  isDisabled={!email}
                  onPress={() => run(handleResend)}
                >
                  Resend verification email
                </Button>
                <LinkButton variant="outline" to="/login">
                  Back to login
                </LinkButton>
              </div>
            </Card.Content>
          </Card>
        </div>
      )}
    </div>
  );
}
