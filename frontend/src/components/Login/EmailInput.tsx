import { FieldError, Input, Label, TextField } from '@heroui/react';
import { InputArgs } from '../../models/Input';

export default function EmailInput({ onChange }: InputArgs) {
  return (
    <TextField
      type="email"
      name="email"
      autoComplete="email"
      validate={(value: string) => (!value.length ? 'Email is required' : null)}
      onChange={onChange}
    >
      <Label>Email</Label>
      <Input />
      <FieldError />
    </TextField>
  );
}
