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

// autoComplete is 'current-password' to sign in and 'new-password' where the
// user sets a password (sign up, reset), so password managers suggest one.
export default function PasswordInput({
  onChange,
  autoComplete,
}: InputArgs & { autoComplete: 'current-password' | 'new-password' }) {
  const [isVisible, setIsVisible] = useState(false);

  return (
    <TextField
      name="password"
      autoComplete={autoComplete}
      validate={(value: string) =>
        !value.length ? 'Password is required' : null
      }
      onChange={onChange}
    >
      <Label>Password</Label>
      <InputGroup>
        <InputGroup.Input type={isVisible ? 'text' : 'password'} />
        <InputGroup.Suffix className="pe-0">
          <Button
            isIconOnly
            size="sm"
            variant="ghost"
            aria-label={isVisible ? 'Hide password' : 'Show password'}
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
