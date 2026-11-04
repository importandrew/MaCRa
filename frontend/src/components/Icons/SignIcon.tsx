import { SVGProps } from 'react';
import { JSX } from 'react/jsx-runtime';

const PATHS = {
  plus: 'M12 5v14M5 12h14',
  minus: 'M5 12h14',
};

// SignIcon draws a + or − centered in its box, unlike the text glyphs,
// which the font sets above the middle of the line.
export const SignIcon = ({
  sign,
  ...props
}: { sign: keyof typeof PATHS } & JSX.IntrinsicAttributes &
  SVGProps<SVGSVGElement>) => {
  return (
    <svg
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
        d={PATHS[sign]}
      />
    </svg>
  );
};
