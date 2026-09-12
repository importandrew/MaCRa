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
import AbortMenu from '../../CurrentGame/Menu/AbortMenu';
import FinishMenu from '../../CurrentGame/Menu/FinishMenu';
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
  const canAbort = hasCapability('abort_game');
  const canFinish = hasCapability('finish_game');

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

        <div className="mt-auto pt-4 flex justify-center gap-3 w-7/8 mx-auto">
          <FinishMenu gameId={gameId} canFinish={canFinish} />
          <AbortMenu gameId={gameId} canAbort={canAbort} />
        </div>
      </DrawerBody>
    </>
  );
}
