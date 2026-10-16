import { Button, Card, CardBody, Spacer, Spinner } from '@heroui/react';
import { useEffect, useRef, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { globalFetch, PostRequest } from '../Request';

export default function VerifyEmailPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') ?? '';
  const [status, setStatus] = useState<'loading' | 'ok' | 'error'>('loading');
  const started = useRef(false);

  useEffect(() => {
    if (!token) {
      setStatus('error');
      return;
    }

    // Strict Mode remounts effects in dev; only send one verify request.
    if (started.current) {
      return;
    }
    started.current = true;

    (async () => {
      const data = await globalFetch<object>(
        PostRequest('/auth/verify-email', { token }),
        { redirectOn401: false }
      );
      setStatus(data ? 'ok' : 'error');
    })();
  }, [token]);

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
        <Card fullWidth className="mx-5">
          <CardBody className="gap-3 py-8">
            {status === 'loading' && <Spinner label="Verifying email..." size="lg" />}
            {status === 'ok' && (
              <>
                <p>Email verified. You can sign in now.</p>
                <Spacer />
                <Button as={Link} to="/login" color="primary">
                  Go to login
                </Button>
              </>
            )}
            {status === 'error' && (
              <>
                <p>This verification link is invalid or expired.</p>
                <Spacer />
                <Button as={Link} to="/login" color="primary" variant="bordered">
                  Back to login
                </Button>
              </>
            )}
          </CardBody>
        </Card>
      </div>
    </div>
  );
}
