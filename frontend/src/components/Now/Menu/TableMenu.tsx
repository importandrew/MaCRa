import {
  Accordion,
  AccordionItem,
  DrawerBody,
  DrawerHeader,
} from '@heroui/react';
import { PlayerOrder } from '../../../models/GameData';
import { BookIcon } from '../../Icons/BookIcon';
import { ListIcon } from '../../Icons/ListIcon';
import { WarningIcon } from '../../Icons/WarningIcon';
import PenaltyInput from './PenaltyInput';
import LogsMenu from './LogsMenu';
import AbortMenu from './AbortMenu';
import FinishMenu from './FinishMenu';
import { FansMenu } from './FansMenu';

export default function TableMenu({
  players = [],
  finedPlayerId = null,
  gameId = null,
  setFinedPlayerId,
  onRefetchNow,
}: {
  players: PlayerOrder[] | null;
  finedPlayerId: number | null;
  gameId: number | null;
  setFinedPlayerId: (value: number | null) => void;
  onRefetchNow?: () => void;
}) {
  return (
    <>
      <DrawerHeader className="flex flex-col gap-3">Table menu</DrawerHeader>
      <DrawerBody className="flex flex-col">
        <Accordion defaultExpandedKeys={[finedPlayerId ? 'penalty' : '']}>
          <AccordionItem
            key="penalty"
            aria-label="Penalty"
            title="Penalty"
            startContent={<WarningIcon className="text-warning" />}
          >
            <PenaltyInput
              players={players}
              finedPlayerId={finedPlayerId}
              setFinedPlayerId={setFinedPlayerId}
            />
          </AccordionItem>
          <AccordionItem
            startContent={<ListIcon />}
            key="logs"
            aria-label="Logs"
            title="Logs"
          >
            <div className="flex flex-col gap-2">
              <LogsMenu
                gameId={gameId}
                canUndoLastRound={true}
                onUndoSuccess={onRefetchNow}
              />
            </div>
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
          <FinishMenu gameId={gameId} />
          <AbortMenu gameId={gameId} />
        </div>
      </DrawerBody>
    </>
  );
}
