import { useTranslation } from 'react-i18next';
import {
  Button,
  Card,
  Input,
  Label,
  Separator,
  TextField,
} from '@heroui/react';
import { useCallback, useEffect, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../providers/auth/authConstants';
import {
  fetchProfileApi,
  updateProfileApi,
} from '../../providers/auth/authApi';
import { UserData } from '../../models/UserData';
import { showToast } from '../Toast';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';
import PendingButton from '../PendingButton';

export default function ProfilePage() {
  const { user, updateUser } = useAuth();
  const [userData, setUserData] = useState<UserData | null>(user);
  const [name, setName] = useState(user?.name ?? '');
  const [isLoading, setIsLoading] = useState(true);
  const { isPending: isSaving, run } = useSubmitOnce();
  const navigate = useNavigate();
  const { t } = useTranslation();

  // The latest profile from the server; it also refreshes the signed-in user
  // everywhere (the name may have changed on another device).
  useEffect(() => {
    fetchProfileApi().then((profileData) => {
      if (profileData) {
        setUserData(profileData);
        setName(profileData.name);
        updateUser(profileData);
      }
      setIsLoading(false);
    });
  }, [updateUser]);

  const handleSave = useCallback(async () => {
    if (!name.trim()) {
      return;
    }

    const updatedData = await updateProfileApi(name.trim());

    if (!updatedData) return;

    // The header shows the new name at once.
    updateUser(updatedData);
    setUserData(updatedData);

    showToast({
      title: t('profile.updated'),
      variant: 'success',
    });
  }, [name, t, updateUser]);

  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-screen md:w-[600px] my-auto justify-center items-center">
        <LoadingSpinner label={t('common.loading')} />
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
        <Card className="w-full mx-5">
          <Card.Header className="items-center rounded-xl bg-accent py-2 text-accent-foreground">
            <Card.Title>{t('profile.title')}</Card.Title>
          </Card.Header>
          <Separator />
          <Card.Content className="gap-4 text-start">
            <TextField
              isDisabled
              fullWidth
              value={userData?.id.toString() || ''}
            >
              <Label>{t('profile.userId')}</Label>
              <Input />
            </TextField>
            <TextField fullWidth value={name} onChange={setName}>
              <Label>{t('auth.name')}</Label>
              <Input placeholder={t('profile.namePlaceholder')} />
            </TextField>
            <div className="flex justify-center mt-3 gap-3">
              <Button variant="tertiary" onPress={() => navigate('/')}>
                {t('common.cancel')}
              </Button>
              <PendingButton
                variant="primary"
                onPress={() => run(handleSave)}
                isPending={isSaving}
                isDisabled={!name.trim() || name.trim() === userData?.name}
              >
                {t('common.save')}
              </PendingButton>
            </div>
          </Card.Content>
        </Card>
      </div>
    </div>
  );
}
