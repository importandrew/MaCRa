import { useTranslation } from 'react-i18next';
import {
  Button,
  FieldError,
  InputGroup,
  Label,
  TextField,
} from '@heroui/react';
import { useState } from 'react';
import { InputArgs } from '../../models/Input';
import { EyeClosedIcon } from './EyeClosed';
import { EyeIcon } from './EyeOpen';

// maxPasswordBytes mirrors the backend limit (bcrypt's).
const maxPasswordBytes = 72;

// autoComplete is 'current-password' to sign in and 'new-password' where the
// user sets a password (sign up, reset), so password managers suggest one.
export default function PasswordInput({
  onChange,
  autoComplete,
}: InputArgs & { autoComplete: 'current-password' | 'new-password' }) {
  const { t } = useTranslation();
  const [isVisible, setIsVisible] = useState(false);

  return (
    <TextField
      name="password"
      autoComplete={autoComplete}
      validate={(value: string) => {
        if (!value.length) {
          return t('auth.passwordRequired');
        }
        // bcrypt on the backend takes at most 72 bytes; only a new password
        // is checked, so an existing one is never refused here.
        if (
          autoComplete === 'new-password' &&
          new TextEncoder().encode(value).length > maxPasswordBytes
        ) {
          return t('auth.passwordTooLong');
        }
        return null;
      }}
      onChange={onChange}
    >
      <Label>{t('auth.password')}</Label>
      <InputGroup>
        <InputGroup.Input type={isVisible ? 'text' : 'password'} />
        <InputGroup.Suffix className="pe-0">
          <Button
            isIconOnly
            size="sm"
            variant="ghost"
            aria-label={
              isVisible ? t('auth.hidePassword') : t('auth.showPassword')
            }
            onPress={() => setIsVisible((prev) => !prev)}
          >
            {isVisible ? (
              <EyeClosedIcon className="text-xl text-muted" />
            ) : (
              <EyeIcon className="text-xl text-muted" />
            )}
          </Button>
        </InputGroup.Suffix>
      </InputGroup>
      <FieldError />
    </TextField>
  );
}
