import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

/** Single mahjong tile with 中 (Player A). */
export const MahjongTileIcon = (
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
      <rect
        x="6.5"
        y="3.5"
        width="11"
        height="17"
        rx="1.5"
        stroke="currentColor"
        strokeWidth="1.5"
      />
      <text
        x="12"
        y="14.5"
        textAnchor="middle"
        fill="currentColor"
        fontSize="9"
        fontWeight="600"
        fontFamily="serif"
      >
        中
      </text>
    </svg>
  );
};
