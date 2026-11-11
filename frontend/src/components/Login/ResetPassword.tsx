import { useTranslation } from 'react-i18next';
import { Button, Card } from '@heroui/react';
import { LinkButton } from '../LinkButton';
import { useCallback, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { globalFetch, PostRequest } from '../Request';
import { showToast } from '../Toast';
import PasswordInput from './PasswordInput';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';

export default function ResetPasswordPage() {
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') ?? '';
  const navigate = useNavigate();

  const [password, setPassword] = useState('');
  const { isPending: isLoading, run } = useSubmitOnce();

  const handleSubmit = useCallback(async () => {
    if (!token) {
      showToast({ title: t('auth.resetTokenMissing'), variant: 'danger' });
      return;
    }

    const data = await globalFetch<object>(
      PostRequest('/auth/reset-password', { token, password }),
      { redirectOn401: false }
    );

    if (!data) {
      return;
    }

    showToast({
      title: t('auth.passwordUpdated'),
      variant: 'success',
    });
    navigate('/login', { replace: true });
  }, [token, password, navigate, t]);

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      {isLoading ? (
        <LoadingSpinner label={t('common.loading')} />
      ) : (
        <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
          <Card className="w-full mx-5">
            <Card.Content className="gap-3 py-6">
              <h1 className="text-xl font-semibold">
                {t('auth.resetPasswordTitle')}
              </h1>
              {!token ? (
                <>
                  <p className="text-base text-muted">
                    {t('auth.resetLinkNoToken')}
                  </p>
                  <LinkButton variant="primary" to="/forgot-password">
                    {t('auth.requestNewLink')}
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
                      onPress={() => run(handleSubmit)}
                    >
                      {t('auth.updatePassword')}
                    </Button>
                    <LinkButton variant="outline" to="/login">
                      {t('common.cancel')}
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
