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
import { passwordProblem } from './validation';

// autoComplete is 'current-password' to sign in and 'new-password' where the
// user sets a password (sign up, reset), so password managers suggest one.
export default function PasswordInput({
  value,
  onChange,
  autoComplete,
}: InputArgs & { autoComplete: 'current-password' | 'new-password' }) {
  const { t } = useTranslation();
  const [isVisible, setIsVisible] = useState(false);

  return (
    <TextField
      name="password"
      autoComplete={autoComplete}
      value={value}
      validate={(v: string) => {
        switch (passwordProblem(v, autoComplete === 'new-password')) {
          case 'required':
            return t('auth.passwordRequired');
          case 'tooLong':
            return t('auth.passwordTooLong');
          default:
            return null;
        }
      }}
      onChange={onChange}
    >
      <Label>{t('auth.password')}</Label>
      <InputGroup>
        <InputGroup.Input
          type={isVisible ? 'text' : 'password'}
          // min-w-0: an input keeps its intrinsic width in the flex row, which
          // pushes the eye button out of the field on a 320px phone.
          className="min-w-0"
        />
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
