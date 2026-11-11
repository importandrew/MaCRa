import { Card, Separator } from '@heroui/react';
import { useTranslation } from 'react-i18next';
import {
  currentLanguage,
  LANGUAGES,
  setLanguage,
  type Language,
} from '../../i18n';
import SegmentedPicker from '../SegmentedPicker';

// SettingsPage holds the settings of the app on this device (Profile is about
// the account). For now only the language; it applies at once.
export default function SettingsPage() {
  // useTranslation re-renders the page when the language changes.
  const { t } = useTranslation();

  return (
    <div className="flex flex-col h-full w-screen md:w-[600px] my-auto">
      <div className="flex flex-col max-w-screen h-full text-2xl text-center justify-center items-center mx-5">
        <Card className="w-full mx-5">
          <Card.Header className="items-center rounded-xl bg-accent py-2 text-accent-foreground">
            <Card.Title>{t('settings.title')}</Card.Title>
          </Card.Header>
          <Separator />
          <Card.Content className="gap-4 text-start">
            <SegmentedPicker<Language>
              label={t('settings.language')}
              options={LANGUAGES.map(({ code, label }) => ({
                value: code,
                label,
              }))}
              selectedKey={currentLanguage()}
              onSelectionChange={setLanguage}
            />
          </Card.Content>
        </Card>
      </div>
    </div>
  );
}
