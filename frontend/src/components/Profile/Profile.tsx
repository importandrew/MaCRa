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
import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';
import {
  fetchProfileApi,
  updateProfileApi,
} from '../../providers/auth/authApi';
import { UserData } from '../../models/UserData';
import { showToast } from '../Toast';

export default function ProfilePage() {
  const [userData, setUserData] = useState<UserData | null>(null);
  const [name, setName] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const loadProfile = async () => {
      setIsLoading(true);

      // Load from localStorage initially
      const localData = JSON.parse(
        localStorage.getItem(AUTH_USER_DATA) || '{}'
      ) as UserData;

      if (localData.id && localData.name) {
        setUserData(localData);
        setName(localData.name);
      }

      // Fetch latest from server
      const profileData = await fetchProfileApi();
      if (profileData) {
        setUserData(profileData);
        setName(profileData.name);
      }

      setIsLoading(false);
    };

    loadProfile();
  }, []);

  const handleSave = useCallback(async () => {
    if (!name.trim()) {
      showToast({
        title: 'Error',
        description: 'Name cannot be empty',
        color: 'danger',
      });
      return;
    }

    setIsSaving(true);
    const updatedData = await updateProfileApi(name.trim());

    if (!updatedData) {
      setIsSaving(false);
      return;
    }

    // Update localStorage
    localStorage.setItem(AUTH_USER_DATA, JSON.stringify(updatedData));
    setUserData(updatedData);

    showToast({
      title: 'Profile updated successfully',
      color: 'success',
    });

    setIsSaving(false);
  }, [name]);

  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-screen md:w-[600px] my-auto justify-center items-center">
        <Spinner label="Loading..." size="lg" />
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
        <Card fullWidth className="mx-5">
          <CardHeader className="flex justify-center bg-blue-400 text-gray-100">
            Profile
          </CardHeader>
          <Divider />
          <CardBody>
            <Input
              label="User ID"
              value={userData?.id || ''}
              isDisabled
              variant="bordered"
              fullWidth
            />
            <Spacer y={4} />
            <Input
              label="Name"
              value={name}
              onValueChange={setName}
              variant="bordered"
              fullWidth
              placeholder="Enter your name"
            />
            <Spacer y={4} />
            <div className="flex justify-center mt-3 gap-3">
              <Button
                color="default"
                variant="light"
                onPress={() => navigate('/')}
              >
                Cancel
              </Button>
              <Button
                color="primary"
                variant="solid"
                onPress={handleSave}
                isLoading={isSaving}
                isDisabled={!name.trim() || name.trim() === userData?.name}
              >
                Save
              </Button>
            </div>
          </CardBody>
        </Card>
      </div>
    </div>
  );
}
