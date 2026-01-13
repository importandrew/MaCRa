import {
  Button,
  Card,
  CardBody,
  CardHeader,
  Divider,
  Input,
  Spacer,
  Spinner,
} from '@heroui/react';
import { useCallback, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import EmailInput from '../Login/EmailInput';
import PasswordInput from '../Login/PasswordInput';
import { globalFetch, PostRequest } from '../Request';
import { addToast } from '@heroui/react';

export default function SignUpPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const navigate = useNavigate();

  const handleSignup = useCallback(async () => {
    setIsLoading(true);

    const signupRequest = PostRequest('/auth/signup', {
      email,
      password,
      name,
    });

    const data = await globalFetch<object>(signupRequest);

    if (!data) {
      setIsLoading(false);
      return;
    }

    addToast({
      title: 'Registration successful! Now login, please',
      color: 'success',
    });
    navigate('/login');

    setIsLoading(false);
  }, [email, password, name, navigate]);

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      {isLoading ? (
        <Spinner label="Loading..." size="lg" />
      ) : (
        <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
          <Card fullWidth className="mx-5">
            <CardHeader className="flex justify-center bg-gray-400 text-gray-100">
              Sign Up
            </CardHeader>
            <Divider />
            <CardBody>
              <EmailInput onChange={setEmail} />
              <Spacer />
              <Input
                type="text"
                label="Name"
                placeholder="John Mah"
                variant="bordered"
                validate={(value: string) => {
                  return !value.length ? 'Name is required' : null;
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
            </CardBody>
          </Card>
        </div>
      )}
    </div>
  );
}
