import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

/** Ranking / leaderboard podium (Phosphor ranking–inspired). */
export const RankingIcon = (
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
        d="M3.5 20.5h17M5 20.5V11h4.5v9.5M9.5 20.5V4h5v16.5M14.5 20.5V14H19v6.5"
      />
    </svg>
  );
};
