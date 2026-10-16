import { Button, Card, CardBody, Spacer, Spinner } from '@heroui/react';
import { useCallback, useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { globalFetch, PostRequest } from '../Request';
import { showToast } from '../Toast';
import PasswordInput from './PasswordInput';

export default function ResetPasswordPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') ?? '';
  const navigate = useNavigate();

  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = useCallback(async () => {
    if (!token) {
      showToast({ title: 'Missing reset token', color: 'danger' });
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
      color: 'success',
    });
    navigate('/login', { replace: true });
  }, [token, password, navigate]);

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      {isLoading ? (
        <Spinner label="Loading..." size="lg" />
      ) : (
        <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
          <Card fullWidth className="mx-5">
            <CardBody className="gap-3 py-6">
              <h1 className="text-xl font-semibold">Reset password</h1>
              {!token ? (
                <>
                  <p className="text-base text-default-600">
                    This reset link is missing a token.
                  </p>
                  <Button as={Link} to="/forgot-password" color="primary">
                    Request a new link
                  </Button>
                </>
              ) : (
                <>
                  <PasswordInput onChange={setPassword} />
                  <Spacer />
                  <div className="flex justify-center gap-3">
                    <Button
                      color="primary"
                      isDisabled={!password}
                      onPress={handleSubmit}
                    >
                      Update password
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
