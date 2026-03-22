import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

export const CrossInCircleIcon = (
  props: JSX.IntrinsicAttributes & SVGProps<SVGSVGElement>
) => {
  return (
    <svg
      width="32px"
      height="32px"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <g
        fill="none"
        fillRule="evenodd"
        strokeLinecap="round"
        strokeLinejoin="round"
        stroke="#000000"
        transform="translate(2 2)"
        strokeWidth="1.3"
      >
        <circle cx="8.5" cy="8.5" r="8" />

        <g transform="matrix(0 1 -1 0 17 0)">
          <path strokeWidth="1.3" d="m5.5 11.5 6-6" />

          <path strokeWidth="1.3" d="m5.5 5.5 6 6" />
        </g>
      </g>
    </svg>
  );
};
