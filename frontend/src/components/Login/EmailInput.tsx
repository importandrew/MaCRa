import { FieldError, Input, Label, TextField } from '@heroui/react';
import { InputArgs } from '../../models/Input';

// maxEmailLength mirrors the backend limit (the users.email column).
const maxEmailLength = 50;

export default function EmailInput({ onChange }: InputArgs) {
  return (
    <TextField
      type="email"
      name="email"
      autoComplete="email"
      validate={(value: string) => {
        if (!value.length) {
          return 'Email is required';
        }
        if ([...value].length > maxEmailLength) {
          return 'Email is too long';
        }
        return null;
      }}
      onChange={onChange}
    >
      <Label>Email</Label>
      <Input />
      <FieldError />
    </TextField>
  );
}
