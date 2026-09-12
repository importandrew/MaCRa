import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

export const TrophyIcon = (
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
        d="M8 4h8v4a4 4 0 0 1-8 0V4Zm0 0H6a2 2 0 0 0-2 2v1a3 3 0 0 0 3 3m10-6h2a2 2 0 0 1 2 2v1a3 3 0 0 1-3 3m-5 5v3m0 0H9m3 0h3m-6 3h6"
      />
    </svg>
  );
};
