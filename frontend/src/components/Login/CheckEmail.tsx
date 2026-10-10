import { Button, Card, CardBody, Spacer, Spinner } from '@heroui/react';
import { useCallback, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { globalFetch, PostRequest } from '../Request';
import { showToast } from '../Toast';
import EmailInput from './EmailInput';

type CheckEmailState = {
  email?: string;
};

export default function CheckEmailPage() {
  const location = useLocation();
  const stateEmail = (location.state as CheckEmailState | null)?.email ?? '';

  const [email, setEmail] = useState(stateEmail);
  const [isLoading, setIsLoading] = useState(false);

  const handleResend = useCallback(async () => {
    if (!email) {
      return;
    }

    setIsLoading(true);
    const data = await globalFetch<object>(
      PostRequest('/auth/resend-verification', {
        email,
        origin: window.location.origin,
      }),
      { redirectOn401: false },
    );
    setIsLoading(false);

    if (!data) {
      return;
    }

    showToast({
      title: 'Verification email sent if the account is unverified',
      color: 'success',
    });
  }, [email]);

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      {isLoading ? (
        <Spinner label="Loading..." size="lg" />
      ) : (
        <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
          <Card fullWidth className="mx-5">
            <CardBody className="gap-3 py-6">
              <h1 className="text-xl font-semibold">Check your email</h1>
              <p className="text-base text-default-600">
                We sent a verification link
                {email ? (
                  <>
                    {' '}
                    to <span className="font-medium text-foreground">{email}</span>
                  </>
                ) : null}
                . Open it to activate your account before signing in.
              </p>
              {!stateEmail && (
                <>
                  <EmailInput onChange={setEmail} />
                  <Spacer />
                </>
              )}
              <div className="flex justify-center gap-3">
                <Button
                  color="primary"
                  isDisabled={!email}
                  onPress={handleResend}
                >
                  Resend verification email
                </Button>
                <Button as={Link} to="/login" variant="bordered">
                  Back to login
                </Button>
              </div>
            </CardBody>
          </Card>
        </div>
      )}
    </div>
  );
}
