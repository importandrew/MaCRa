import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

export const WhistleIcon = (
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
        d="M9.5 14a3.5 3.5 0 1 0 0-7 3.5 3.5 0 0 0 0 7Zm0 0v2m3.2-6.5L20 7v3l-4.5 1.5"
      />
    </svg>
  );
};
