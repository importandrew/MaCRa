import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

/** Shield badge with whistle (Referee D). */
export const RefereeBadgeIcon = (
  props: JSX.IntrinsicAttributes & SVGProps<SVGSVGElement>
) => {
  return (
    <svg
      className="w-6 h-6 text-gray-800 dark:text-white"
      aria-hidden="true"
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      fill="none"
      viewBox="0 0 24 24"
      {...props}
    >
      <path
        stroke="currentColor"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.5"
        d="M12 3.5 18.5 6v5.2c0 4.1-2.7 7.4-6.5 8.8-3.8-1.4-6.5-4.7-6.5-8.8V6L12 3.5Z"
      />
      <path
        stroke="currentColor"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.5"
        d="M10.5 12.5a2 2 0 1 0 0-4 2 2 0 0 0 0 4Zm2-2.8L15 8.5v1.5l-2 1"
      />
    </svg>
  );
};
