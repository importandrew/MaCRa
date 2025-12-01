import { Input } from '@heroui/react';
import { InputArgs } from '../../models/Input';

export default function EmailInput({ onChange }: InputArgs) {
  return (
    <Input
      type="email"
      label="Email"
      variant="bordered"
      placeholder=" "
      validate={(value: string) => {
        return !value.length ? 'Email is required' : null;
      }}
      onChange={(event) => onChange(event.target.value)}
    />
  );
}
