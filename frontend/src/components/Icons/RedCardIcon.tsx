import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

/** Person holding a card (Referee C). */
export const RedCardIcon = (
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
        d="M9 10a2.5 2.5 0 1 0 0-5 2.5 2.5 0 0 0 0 5Zm-3.5 9a3.5 3.5 0 0 1 7 0"
      />
      <rect
        x="14"
        y="4"
        width="6"
        height="8"
        rx="0.75"
        stroke="currentColor"
        strokeWidth="1.5"
      />
      <path
        stroke="currentColor"
        strokeLinecap="round"
        strokeWidth="1.5"
        d="M12.5 11.5 14 9"
      />
    </svg>
  );
};
