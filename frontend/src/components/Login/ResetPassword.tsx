import { Button, Card } from '@heroui/react';
import { LinkButton } from '../LinkButton';
import { useCallback, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { globalFetch, PostRequest } from '../Request';
import { showToast } from '../Toast';
import PasswordInput from './PasswordInput';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';

export default function ResetPasswordPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') ?? '';
  const navigate = useNavigate();

  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = useCallback(async () => {
    if (!token) {
      showToast({ title: 'Missing reset token', variant: 'danger' });
      return;
    }

    setIsLoading(true);
    const data = await globalFetch<object>(
      PostRequest('/auth/reset-password', { token, password }),
      { redirectOn401: false }
    );
    setIsLoading(false);

    if (!data) {
      return;
    }

    showToast({
      title: 'Password updated. Please sign in.',
      variant: 'success',
    });
    navigate('/login', { replace: true });
  }, [token, password, navigate]);

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      {isLoading ? (
        <LoadingSpinner label="Loading..." />
      ) : (
        <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
          <Card className="w-full mx-5">
            <Card.Content className="gap-3 py-6">
              <h1 className="text-xl font-semibold">Reset password</h1>
              {!token ? (
                <>
                  <p className="text-base text-muted">
                    This reset link is missing a token.
                  </p>
                  <LinkButton variant="primary" to="/forgot-password">
                    Request a new link
                  </LinkButton>
                </>
              ) : (
                <>
                  <PasswordInput
                    onChange={setPassword}
                    autoComplete="new-password"
                  />
                  <div className="flex justify-center gap-3">
                    <Button
                      variant="primary"
                      isDisabled={!password}
                      onPress={handleSubmit}
                    >
                      Update password
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
