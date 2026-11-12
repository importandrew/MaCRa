import type { ReactNode } from 'react';
import './Table.css';

// TableLayout is the table of the game screens (/now and the live view): a
// 3×3 grid with the seats around the center, Refresh in the bottom left
// corner and the ≡ menu in the bottom right one. Each seat is turned to face
// its side of the table (Table.css). The caller fills the center cell itself,
// as the screens lay it out differently.
export default function TableLayout({
  top,
  left,
  right,
  bottom,
  center,
  refresh,
  menu,
  className = '',
}: {
  top: ReactNode;
  left: ReactNode;
  right: ReactNode;
  bottom: ReactNode;
  center: ReactNode;
  refresh: ReactNode;
  menu: ReactNode;
  className?: string;
}) {
  return (
    <div className={`flex justify-center h-full w-screen ${className}`}>
      <div className="grid w-screen grid-cols-[25%_50%_25%] grid-rows-[20%_60%_20%] text-center">
        <div></div>
        <div className="relative self-center h-full w-full">
          <div className="table-top-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
            {top}
          </div>
        </div>
        <div></div>
        <div className="relative self-center h-full w-full">
          <div className="table-right-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
            {left}
          </div>
        </div>
        {center}
        <div className="relative self-center h-full w-full">
          <div className="table-left-side absolute flex flex-col w-screen md:w-max whitespace-nowrap">
            {right}
          </div>
        </div>
        <div>{refresh}</div>
        <div className="relative self-center h-full w-full">
          <div className="table-bottom-side absolute px-6 flex flex-col w-screen whitespace-nowrap">
            {bottom}
          </div>
        </div>
        <div className="flex justify-end items-end">{menu}</div>
      </div>
    </div>
  );
}
