import { Spinner } from '@heroui/react';

// LoadingSpinner is a spinner with a text label under it (HeroUI v3 Spinner
// has no label).
export default function LoadingSpinner({
  label,
  size = 'lg',
}: {
  label: string;
  size?: 'sm' | 'md' | 'lg';
}) {
  return (
    <div
      role="status"
      aria-label={label}
      className="flex flex-col items-center gap-2"
    >
      <Spinner size={size} />
      <span className="text-sm text-muted">{label}</span>
    </div>
  );
}
