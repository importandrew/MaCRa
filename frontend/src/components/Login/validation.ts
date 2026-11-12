// The checks of the auth fields, shared by each field (its error text) and
// its form (the submit button stays disabled while a field is invalid). The
// limits mirror the backend's.

export type FieldProblem = 'required' | 'tooLong' | null;

// maxEmailLength is the users.email column's limit.
const maxEmailLength = 50;
// maxNameLength is the users.name column's limit.
const maxNameLength = 50;
// maxPasswordBytes is bcrypt's limit.
const maxPasswordBytes = 72;

export function emailProblem(value: string): FieldProblem {
  if (!value.length) return 'required';
  if ([...value].length > maxEmailLength) return 'tooLong';
  return null;
}

export function nameProblem(value: string): FieldProblem {
  if (!value.length) return 'required';
  if ([...value].length > maxNameLength) return 'tooLong';
  return null;
}

// passwordProblem checks the length only for a new password, so an existing
// one is never refused here.
export function passwordProblem(value: string, isNew: boolean): FieldProblem {
  if (!value.length) return 'required';
  if (isNew && new TextEncoder().encode(value).length > maxPasswordBytes) {
    return 'tooLong';
  }
  return null;
}
