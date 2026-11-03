import { Card } from '@heroui/react';
import { LinkButton } from '../LinkButton';
import { useEffect, useRef, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { globalFetch, PostRequest } from '../Request';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';

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
        <Card className="w-full mx-5">
          <Card.Content className="gap-3 py-8">
            {status === 'loading' && (
              <LoadingSpinner label="Verifying email..." />
            )}
            {status === 'ok' && (
              <>
                <p>Email verified. You can sign in now.</p>
                <LinkButton variant="primary" to="/login">
                  Go to login
                </LinkButton>
              </>
            )}
            {status === 'error' && (
              <>
                <p>This verification link is invalid or expired.</p>
                <LinkButton variant="outline" to="/login">
                  Back to login
                </LinkButton>
              </>
            )}
          </Card.Content>
        </Card>
      </div>
    </div>
  );
}
