import { useTranslation } from 'react-i18next';
import { FieldError, Input, Label, TextField } from '@heroui/react';
import { InputArgs } from '../../models/Input';
import { emailProblem } from './validation';

export default function EmailInput({ value, onChange }: InputArgs) {
  const { t } = useTranslation();
  return (
    <TextField
      type="email"
      name="email"
      autoComplete="email"
      value={value}
      validate={(v: string) => {
        switch (emailProblem(v)) {
          case 'required':
            return t('auth.emailRequired');
          case 'tooLong':
            return t('auth.emailTooLong');
          default:
            return null;
        }
      }}
      onChange={onChange}
    >
      <Label>{t('auth.email')}</Label>
      <Input />
      <FieldError />
    </TextField>
  );
}
