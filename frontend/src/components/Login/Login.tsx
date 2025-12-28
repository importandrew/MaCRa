import {
  Button,
  Card,
  CardBody,
  CardHeader,
  Divider,
  Spacer,
  Spinner,
} from '@heroui/react';
import { useCallback, useContext, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { AuthContext } from '../../providers/auth/authConstants';
import { globalFetch, PostRequest } from '../Request';
import EmailInput from './EmailInput';
import PasswordInput from './PasswordInput';
import { UserData } from '../../models/UserData';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const { signIn } = useContext(AuthContext);
  const navigate = useNavigate();

  const handleLogin = useCallback(async () => {
    setIsLoading(true);

    const loginRequest = PostRequest('/auth/login', {
      email,
      password,
    });

    const data = await globalFetch<UserData>(loginRequest);

    if (!data || 'error' in data) {
      setIsLoading(false);
      return;
    }

    signIn(data);
    navigate('/');
    setIsLoading(false);
  }, [email, password, signIn, navigate]);

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      {isLoading ? (
        <Spinner label="Loading..." size="lg" />
      ) : (
        <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
          <Card fullWidth className="mx-5">
            <CardHeader className="flex justify-center bg-gray-400 text-gray-100">
              Login
            </CardHeader>
            <Divider />
            <CardBody>
              <EmailInput onChange={setEmail} />
              <Spacer />
              <PasswordInput onChange={setPassword} />
              <Spacer />
              <div className="flex justify-center mt-3">
                <Button
                  color="primary"
                  isDisabled={!email || !password}
                  variant="solid"
                  onPress={handleLogin}
                >
                  Sign In
                </Button>
              </div>
            </CardBody>
          </Card>
        </div>
      )}
    </div>
  );
}
