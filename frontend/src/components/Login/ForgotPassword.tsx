import { useTranslation } from 'react-i18next';
import { Button, Card } from '@heroui/react';
import { LinkButton } from '../LinkButton';
import { useCallback, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { globalFetch, PostRequest } from '../Request';
import { showToast } from '../Toast';
import EmailInput from './EmailInput';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';

export default function ForgotPasswordPage() {
  const { t } = useTranslation();
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
      title: t('auth.resetLinkSent'),
      variant: 'success',
    });
  }, [email, t]);

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      {isLoading ? (
        <LoadingSpinner label={t('common.loading')} />
      ) : (
        <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
          <Card className="w-full mx-5">
            <Card.Content className="gap-3 py-6">
              <h1 className="text-xl font-semibold">
                {t('auth.forgotPasswordTitle')}
              </h1>
              {sent ? (
                <>
                  <p className="text-base text-muted">
                    {t('auth.checkEmailForReset')}
                  </p>
                  <LinkButton variant="primary" to="/login">
                    {t('auth.backToLogin')}
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
                      {t('auth.sendResetLink')}
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
