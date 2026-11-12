import { useTranslation } from 'react-i18next';
import { Card, FieldError, Input, Label, Tabs, TextField } from '@heroui/react';
import { Key, useCallback, useContext, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import {
  Link,
  useLocation,
  useNavigate,
  useSearchParams,
} from 'react-router-dom';
import { AuthContext } from '../../providers/auth/authConstants';
import { UserData } from '../../models/UserData';
import { fetchWithStatus, globalFetch, PostRequest } from '../Request';
import PendingButton from '../PendingButton';
import EmailInput from './EmailInput';
import { emailProblem, nameProblem, passwordProblem } from './validation';
import PasswordInput from './PasswordInput';

type AuthTab = 'login' | 'signup';

function tabFromSearchParams(searchParams: URLSearchParams): AuthTab {
  return searchParams.get('tab') === 'signup' ? 'signup' : 'login';
}

export default function LoginPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const tab = tabFromSearchParams(searchParams);

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const { isPending: isLoading, run } = useSubmitOnce();
  const { t } = useTranslation();

  const { signIn } = useContext(AuthContext);
  const navigate = useNavigate();
  const location = useLocation();
  const from = location.state?.from?.pathname || '/';

  const setTab = useCallback(
    (next: AuthTab) => {
      if (next === 'signup') {
        setSearchParams({ tab: 'signup' }, { replace: true });
      } else {
        setSearchParams({}, { replace: true });
      }
    },
    [setSearchParams]
  );

  const handleTabChange = useCallback(
    (key: Key) => {
      setTab(key === 'signup' ? 'signup' : 'login');
    },
    [setTab]
  );

  const handleLogin = useCallback(async () => {
    const result = await fetchWithStatus<UserData>(
      PostRequest('/auth/login', { email, password }),
      { redirectOn401: false }
    );

    if (result.status === 403) {
      navigate('/check-email', { state: { email }, replace: true });
      return;
    }

    if (!result.data) {
      return;
    }

    signIn(result.data);
    navigate(from, { replace: true });
  }, [email, password, signIn, navigate, from]);

  const handleSignup = useCallback(async () => {
    const data = await globalFetch<object>(
      PostRequest('/auth/signup', {
        email,
        password,
        name,
      })
    );

    if (!data) {
      return;
    }

    navigate('/check-email', { state: { email }, replace: true });
  }, [email, password, name, navigate]);

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
        <Card className="w-full mx-5">
          <Card.Content>
            <Tabs selectedKey={tab} onSelectionChange={handleTabChange}>
              <Tabs.ListContainer>
                <Tabs.List aria-label={t('auth.tabs')} className="w-full">
                  <Tabs.Tab id="login" className="flex-1">
                    {t('auth.logIn')}
                    <Tabs.Indicator />
                  </Tabs.Tab>
                  <Tabs.Tab id="signup" className="flex-1">
                    {t('auth.signUp')}
                    <Tabs.Indicator />
                  </Tabs.Tab>
                </Tabs.List>
              </Tabs.ListContainer>
              <Tabs.Panel id="login">
                <div className="flex flex-col gap-3 pt-2 text-start">
                  <EmailInput value={email} onChange={setEmail} />
                  <PasswordInput
                    value={password}
                    onChange={setPassword}
                    autoComplete="current-password"
                  />
                  <div className="flex flex-col items-center gap-2 mt-3">
                    <PendingButton
                      variant="primary"
                      isPending={isLoading}
                      isDisabled={
                        emailProblem(email) != null ||
                        passwordProblem(password, false) != null
                      }
                      onPress={() => run(handleLogin)}
                    >
                      {t('auth.signIn')}
                    </PendingButton>
                    <Link
                      to="/forgot-password"
                      className="text-sm text-accent underline"
                    >
                      {t('auth.forgotPasswordLink')}
                    </Link>
                  </div>
                </div>
              </Tabs.Panel>
              <Tabs.Panel id="signup">
                <div className="flex flex-col gap-3 pt-2 text-start">
                  <EmailInput value={email} onChange={setEmail} />
                  <TextField
                    name="name"
                    autoComplete="name"
                    value={name}
                    validate={(value: string) => {
                      switch (nameProblem(value)) {
                        case 'required':
                          return t('auth.nameRequired');
                        case 'tooLong':
                          return t('auth.nameTooLong');
                        default:
                          return null;
                      }
                    }}
                    onChange={setName}
                  >
                    <Label>{t('auth.name')}</Label>
                    <Input placeholder={t('auth.namePlaceholder')} />
                    <FieldError />
                  </TextField>
                  <PasswordInput
                    value={password}
                    onChange={setPassword}
                    autoComplete="new-password"
                  />
                  <div className="flex justify-center">
                    <PendingButton
                      variant="primary"
                      isPending={isLoading}
                      isDisabled={
                        emailProblem(email) != null ||
                        nameProblem(name) != null ||
                        passwordProblem(password, true) != null
                      }
                      className="mt-3"
                      onPress={() => run(handleSignup)}
                    >
                      {t('auth.signUp')}
                    </PendingButton>
                  </div>
                </div>
              </Tabs.Panel>
            </Tabs>
          </Card.Content>
        </Card>
      </div>
    </div>
  );
}
