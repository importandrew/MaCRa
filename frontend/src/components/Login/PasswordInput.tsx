import { Input } from '@heroui/react';
import { useCallback, useState } from 'react';
import { InputArgs } from '../../models/Input';
import { EyeClosedIcon } from './EyeClosed';
import { EyeIcon } from './EyeOpen';

export default function PasswordInput({ onChange }: InputArgs) {
  const [isVisible, setIsVisible] = useState(false);
  const toggleVisibility = useCallback(() => setIsVisible((prev) => !prev), []);

  return (
    <Input
      label="Password"
      variant="bordered"
      placeholder=" "
      validate={(value: string) => {
        return !value.length ? 'Password is required' : null;
      }}
      endContent={
        <button
          className="focus:outline-none"
          type="button"
          onClick={toggleVisibility}
        >
          {isVisible ? (
            <EyeClosedIcon className="text-2xl text-default-400 pointer-events-none" />
          ) : (
            <EyeIcon className="text-2xl text-default-400 pointer-events-none" />
          )}
        </button>
      }
      type={isVisible ? 'text' : 'password'}
      onChange={(event) => onChange(event.target.value)}
    />
  );
}
