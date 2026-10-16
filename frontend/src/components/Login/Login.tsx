import {
  Button,
  Card,
  CardBody,
  Input,
  Spacer,
  Spinner,
  Tab,
  Tabs,
} from '@heroui/react';
import { Key, useCallback, useContext, useState } from 'react';
import { Link, useLocation, useNavigate, useSearchParams } from 'react-router-dom';
import { AuthContext } from '../../providers/auth/authConstants';
import { UserData } from '../../models/UserData';
import { fetchWithStatus, globalFetch, PostRequest } from '../Request';
import EmailInput from './EmailInput';
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
  const [isLoading, setIsLoading] = useState(false);

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
    [setSearchParams],
  );

  const handleTabChange = useCallback(
    (key: Key) => {
      setTab(key === 'signup' ? 'signup' : 'login');
    },
    [setTab],
  );

  const handleLogin = useCallback(async () => {
    setIsLoading(true);

    const result = await fetchWithStatus<UserData>(
      PostRequest('/auth/login', { email, password }),
      { redirectOn401: false },
    );

    if (result.status === 403) {
      navigate('/check-email', { state: { email }, replace: true });
      return;
    }

    if (!result.data) {
      setIsLoading(false);
      return;
    }

    signIn(result.data);
    navigate(from, { replace: true });
    setIsLoading(false);
  }, [email, password, signIn, navigate, from]);

  const handleSignup = useCallback(async () => {
    setIsLoading(true);

    const data = await globalFetch<object>(
      PostRequest('/auth/signup', {
        email,
        password,
        name,
        origin: window.location.origin,
      }),
    );

    if (!data) {
      setIsLoading(false);
      return;
    }

    navigate('/check-email', { state: { email }, replace: true });
  }, [email, password, name, navigate]);

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      {isLoading ? (
        <Spinner label="Loading..." size="lg" />
      ) : (
        <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
          <Card fullWidth className="mx-5">
            <CardBody>
              <Tabs
                aria-label="Authentication"
                fullWidth
                selectedKey={tab}
                onSelectionChange={handleTabChange}
              >
                <Tab key="login" title="Login">
                  <div className="pt-2">
                    <EmailInput onChange={setEmail} />
                    <Spacer />
                    <PasswordInput onChange={setPassword} />
                    <Spacer />
                    <div className="flex flex-col items-center gap-2 mt-3">
                      <Button
                        color="primary"
                        isDisabled={!email || !password}
                        variant="solid"
                        onPress={handleLogin}
                      >
                        Sign In
                      </Button>
                      <Link
                        to="/forgot-password"
                        className="text-sm text-primary underline"
                      >
                        Forgot password?
                      </Link>
                    </div>
                  </div>
                </Tab>
                <Tab key="signup" title="Sign Up">
                  <div className="pt-2">
                    <EmailInput onChange={setEmail} />
                    <Spacer />
                    <Input
                      type="text"
                      label="Name"
                      placeholder="John Mah"
                      variant="bordered"
                      validate={(value: string) => {
                        if (!value.length) {
                          return 'Name is required';
                        }
                        if (value.length > 50) {
                          return 'Name is too long';
                        }
                        return null;
                      }}
                      onChange={(event) => setName(event.target.value)}
                    />
                    <Spacer />
                    <PasswordInput onChange={setPassword} />
                    <Spacer />
                    <div className="flex justify-center">
                      <Button
                        color="primary"
                        isDisabled={!email || !password || !name}
                        variant="solid"
                        className="mt-3"
                        onPress={handleSignup}
                      >
                        Sign Up
                      </Button>
                    </div>
                  </div>
                </Tab>
              </Tabs>
            </CardBody>
          </Card>
        </div>
      )}
    </div>
  );
}
