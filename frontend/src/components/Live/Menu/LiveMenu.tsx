import {
  Accordion,
  AccordionItem,
  DrawerBody,
  DrawerHeader,
} from '@heroui/react';
import { useEventCapabilities } from '../../../providers/event/useEventCapabilities';
import { BookIcon } from '../../Icons/BookIcon';
import { ListIcon } from '../../Icons/ListIcon';
import LogsMenu from '../../CurrentGame/Menu/LogsMenu';
import { FansMenu } from '../../CurrentGame/Menu/FansMenu';

export default function LiveMenu({
  gameId = null,
  eventId,
  onUndoSuccess,
}: {
  gameId: number | null;
  eventId?: number;
  onUndoSuccess?: () => void;
}) {
  const { hasCapability } = useEventCapabilities(eventId ?? null);
  const canUndoLastRound = hasCapability('undo_round');

  return (
    <>
      <DrawerHeader className="flex flex-col gap-3">Table menu</DrawerHeader>
      <DrawerBody className="flex flex-col">
        <Accordion>
          <AccordionItem
            startContent={<ListIcon />}
            key="logs"
            aria-label="Logs"
            title="Logs"
          >
            <LogsMenu
              gameId={gameId}
              eventId={eventId}
              canUndoLastRound={canUndoLastRound}
              onUndoSuccess={onUndoSuccess}
            />
          </AccordionItem>
          <AccordionItem
            startContent={<BookIcon />}
            key="fans"
            aria-label="Fans"
            title="Fans"
          >
            <FansMenu />
          </AccordionItem>
        </Accordion>
      </DrawerBody>
    </>
  );
}
