import { Avatar } from '@heroui/react';

// initials returns the first letter of each of the first two words, like the
// HeroUI v2 Avatar did from its `name` prop ("Dev One" → "DO").
function initials(name: string): string {
  return name
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((word) => word.charAt(0))
    .join('')
    .toUpperCase();
}

// UserAvatar shows a player's initials in a round avatar.
export default function UserAvatar({
  name,
  size = 'md',
  className,
}: {
  name: string;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}) {
  return (
    <Avatar size={size} className={className}>
      <Avatar.Fallback>{initials(name)}</Avatar.Fallback>
    </Avatar>
  );
}
