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
import { FansMenu } from './FansMenu';

export default function TableMenu({
  players = [],
  finedPlayerId = null,
  gameId = null,
  setFinedPlayerId,
}: {
  players: PlayerOrder[] | null;
  finedPlayerId: number | null;
  gameId: number | null;
  setFinedPlayerId: (value: number | null) => void;
}) {
  return (
    <>
      <DrawerHeader className="flex flex-col gap-3">Table menu</DrawerHeader>
      <DrawerBody>
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
              gameId={gameId}
              setFinedPlayerId={setFinedPlayerId}
            />
          </AccordionItem>
          <AccordionItem
            startContent={<ListIcon />}
            key="logs"
            aria-label="Logs"
            title="Logs"
          >
            <LogsMenu gameId={gameId}></LogsMenu>
          </AccordionItem>
          <AccordionItem
            startContent={<BookIcon />}
            key="fans"
            aria-label="Fans"
            title="Fans"
          >
            <FansMenu></FansMenu>
          </AccordionItem>
        </Accordion>
      </DrawerBody>
    </>
  );
}
