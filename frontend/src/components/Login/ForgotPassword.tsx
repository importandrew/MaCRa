import { Button, Card, CardBody, Spacer, Spinner } from '@heroui/react';
import { useCallback, useState } from 'react';
import { Link } from 'react-router-dom';
import { globalFetch, PostRequest } from '../Request';
import { showToast } from '../Toast';
import EmailInput from './EmailInput';

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [sent, setSent] = useState(false);

  const handleSubmit = useCallback(async () => {
    setIsLoading(true);

    const data = await globalFetch<object>(
      PostRequest('/auth/forgot-password', {
        email,
        origin: window.location.origin,
      }),
      { redirectOn401: false }
    );

    setIsLoading(false);
    if (!data) {
      return;
    }

    setSent(true);
    showToast({
      title: 'If that email exists, a reset link was sent',
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
              <h1 className="text-xl font-semibold">Forgot password</h1>
              {sent ? (
                <>
                  <p className="text-base text-default-600">
                    Check your email for a reset link.
                  </p>
                  <Button as={Link} to="/login" color="primary">
                    Back to login
                  </Button>
                </>
              ) : (
                <>
                  <EmailInput onChange={setEmail} />
                  <Spacer />
                  <div className="flex justify-center gap-3">
                    <Button
                      color="primary"
                      isDisabled={!email}
                      onPress={handleSubmit}
                    >
                      Send reset link
                    </Button>
                    <Button as={Link} to="/login" variant="bordered">
                      Cancel
                    </Button>
                  </div>
                </>
              )}
            </CardBody>
          </Card>
        </div>
      )}
    </div>
  );
}
