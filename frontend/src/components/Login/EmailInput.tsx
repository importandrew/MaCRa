import { useTranslation } from 'react-i18next';
import { FieldError, Input, Label, TextField } from '@heroui/react';
import { InputArgs } from '../../models/Input';

// maxEmailLength mirrors the backend limit (the users.email column).
const maxEmailLength = 50;

export default function EmailInput({ onChange }: InputArgs) {
  const { t } = useTranslation();
  return (
    <TextField
      type="email"
      name="email"
      autoComplete="email"
      validate={(value: string) => {
        if (!value.length) {
          return t('auth.emailRequired');
        }
        if ([...value].length > maxEmailLength) {
          return t('auth.emailTooLong');
        }
        return null;
      }}
      onChange={onChange}
    >
      <Label>{t('auth.email')}</Label>
      <Input />
      <FieldError />
    </TextField>
  );
}
