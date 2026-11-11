import { Button, Card } from '@heroui/react';
import { LinkButton } from '../LinkButton';
import { useCallback, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { globalFetch, PostRequest } from '../Request';
import { showToast } from '../Toast';
import EmailInput from './EmailInput';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState('');
  const { isPending: isLoading, run } = useSubmitOnce();
  const [sent, setSent] = useState(false);

  const handleSubmit = useCallback(async () => {
    const data = await globalFetch<object>(
      PostRequest('/auth/forgot-password', {
        email,
      }),
      { redirectOn401: false }
    );

    if (!data) {
      return;
    }

    setSent(true);
    showToast({
      title: 'If that email exists, a reset link was sent',
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
              <h1 className="text-xl font-semibold">Forgot password</h1>
              {sent ? (
                <>
                  <p className="text-base text-muted">
                    Check your email for a reset link.
                  </p>
                  <LinkButton variant="primary" to="/login">
                    Back to login
                  </LinkButton>
                </>
              ) : (
                <>
                  <EmailInput onChange={setEmail} />
                  <div className="flex justify-center gap-3">
                    <Button
                      variant="primary"
                      isDisabled={!email}
                      onPress={() => run(handleSubmit)}
                    >
                      Send reset link
                    </Button>
                    <LinkButton variant="outline" to="/login">
                      Cancel
                    </LinkButton>
                  </div>
                </>
              )}
            </Card.Content>
          </Card>
        </div>
      )}
    </div>
  );
}
