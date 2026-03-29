import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

export const CheckInCircleIcon = (
  props: JSX.IntrinsicAttributes & SVGProps<SVGSVGElement>
) => {
  return (
    <svg
      fill="#000000"
      width="32px"
      height="32px"
      viewBox="0 0 36 36"
      version="1.1"
      preserveAspectRatio="xMidYMid meet"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <path
        // stroke="currentColor"
        className="clr-i-outline clr-i-outline-path-1"
        d="M18,6A12,12,0,1,0,30,18,12,12,0,0,0,18,6Zm0,22A10,10,0,1,1,28,18,10,10,0,0,1,18,28Z"
      ></path>
      <path
        // stroke="currentColor"
        className="clr-i-outline clr-i-outline-path-2"
        d="M16.34,23.74l-5-5a1,1,0,0,1,1.41-1.41l3.59,3.59,6.78-6.78a1,1,0,0,1,1.41,1.41Z"
      ></path>
      <rect x="0" y="0" width="32" height="32" fillOpacity="0" />
    </svg>
  );
};
