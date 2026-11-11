import { useTranslation } from 'react-i18next';
import { Card } from '@heroui/react';
import { LinkButton } from '../LinkButton';
import { useEffect, useRef, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { globalFetch, PostRequest } from '../Request';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';

export default function VerifyEmailPage() {
  const { t } = useTranslation();
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
              <LoadingSpinner label={t('auth.verifyingEmail')} />
            )}
            {status === 'ok' && (
              <>
                <p>{t('auth.emailVerified')}</p>
                <LinkButton variant="primary" to="/login">
                  {t('auth.goToLogin')}
                </LinkButton>
              </>
            )}
            {status === 'error' && (
              <>
                <p>{t('auth.verificationLinkInvalid')}</p>
                <LinkButton variant="outline" to="/login">
                  {t('auth.backToLogin')}
                </LinkButton>
              </>
            )}
          </Card.Content>
        </Card>
      </div>
    </div>
  );
}
