import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

/** Judge's gavel (Referee A). */
export const GavelIcon = (
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
        d="M14.5 5.5 18 9l-2.2 2.2-3.5-3.5L14.5 5.5Zm-2.8 2.8L5.5 14.5M4 20h7M8 17.5h4"
      />
      <path
        stroke="currentColor"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.5"
        d="M12.5 6.5 16 10"
      />
    </svg>
  );
};
