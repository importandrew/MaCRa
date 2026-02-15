import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

export const UndoIcon = (
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
        d="M5 10h7a4 4 0 0 1 4 4v2M5 10l3-3m-3 3l3 3"
      />
    </svg>
  );
};
