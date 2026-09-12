import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

/** Scales of justice (Referee B). */
export const ScalesIcon = (
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
        d="M12 3v15m-5 3h10M12 6l-6.5 3.5M12 6l6.5 3.5M5.5 9.5 4 14h3l-1.5-4.5Zm13 0L17 14h3l-1.5-4.5Z"
      />
    </svg>
  );
};
