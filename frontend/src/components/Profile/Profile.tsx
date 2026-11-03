import {
  Button,
  Card,
  Input,
  Label,
  Separator,
  Spinner,
  TextField,
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
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';

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
        variant: 'danger',
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
      variant: 'success',
    });

    setIsSaving(false);
  }, [name]);

  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-screen md:w-[600px] my-auto justify-center items-center">
        <LoadingSpinner label="Loading..." />
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
        <Card className="w-full mx-5">
          <Card.Header className="items-center rounded-xl bg-blue-400 py-2 text-gray-100">
            <Card.Title>Profile</Card.Title>
          </Card.Header>
          <Separator />
          <Card.Content className="gap-4 text-start">
            <TextField
              isDisabled
              fullWidth
              value={userData?.id.toString() || ''}
            >
              <Label>User ID</Label>
              <Input />
            </TextField>
            <TextField fullWidth value={name} onChange={setName}>
              <Label>Name</Label>
              <Input placeholder="Enter your name" />
            </TextField>
            <div className="flex justify-center mt-3 gap-3">
              <Button variant="ghost" onPress={() => navigate('/')}>
                Cancel
              </Button>
              <Button
                variant="primary"
                onPress={handleSave}
                isPending={isSaving}
                isDisabled={!name.trim() || name.trim() === userData?.name}
              >
                {({ isPending }) => (
                  <>
                    {isPending && <Spinner color="current" size="sm" />}
                    Save
                  </>
                )}
              </Button>
            </div>
          </Card.Content>
        </Card>
      </div>
    </div>
  );
}
