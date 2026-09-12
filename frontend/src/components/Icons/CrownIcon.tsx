import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

export const CrownIcon = (
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
        d="M5 17h14l1.5-9-4.5 3L12 5l-4 6L3.5 8 5 17Zm0 0h14v2H5v-2Z"
      />
    </svg>
  );
};
