import { ListBox } from '@heroui/react';
import { ReactNode } from 'react';
import { ChevronRightIcon } from './Icons/ChevronRightIcon';

export type NavListItem = {
  id: string;
  href: string;
  textValue: string;
  content: ReactNode;
};

// NavList is a bordered list of links (played games, event games). Items are
// real links, so long-press and "open in new tab" work on phones.
export default function NavList({
  ariaLabel,
  items,
  emptyContent,
}: {
  ariaLabel: string;
  items: NavListItem[];
  emptyContent: ReactNode;
}) {
  return (
    <div className="w-full max-w-[600px] md:border rounded-sm border-border">
      {items.length === 0 ? (
        <div className="py-6">{emptyContent}</div>
      ) : (
        <ListBox aria-label={ariaLabel} selectionMode="none">
          {items.map((item) => (
            <ListBox.Item
              key={item.id}
              id={item.id}
              href={item.href}
              textValue={item.textValue}
              className="border-b border-separator last:border-b-0"
            >
              <div className="flex w-full min-w-0 items-center gap-2">
                <div className="min-w-0 flex-1">{item.content}</div>
                <ChevronRightIcon className="text-xl text-muted shrink-0" />
              </div>
            </ListBox.Item>
          ))}
        </ListBox>
      )}
    </div>
  );
}
