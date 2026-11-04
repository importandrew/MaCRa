import { Drawer } from '@heroui/react';
import { PlayerOrder } from '../../../models/GameData';
import { useFetchEventMe } from '../../../providers/event/eventApi';
import { useEvent } from '../../../providers/event/eventConstants';
import { BookIcon } from '../../Icons/BookIcon';
import { ListIcon } from '../../Icons/ListIcon';
import { WarningIcon } from '../../Icons/WarningIcon';
import PenaltyInput from './PenaltyInput';
import LogsMenu from './LogsMenu';
import AbortMenu from './AbortMenu';
import FinishMenu from './FinishMenu';
import { FansMenu } from './FansMenu';
import LazyAccordion from '../../LazyAccordion';

export default function TableMenu({
  players = [],
  finedPlayerId = null,
  gameId = null,
  setFinedPlayerId,
  onRefetchNow,
  onPenaltyAdded,
}: {
  players: PlayerOrder[] | null;
  finedPlayerId: number | null;
  gameId: number | null;
  setFinedPlayerId: (value: number | null) => void;
  onRefetchNow?: () => void;
  onPenaltyAdded?: () => void;
}) {
  const { selectedEventId } = useEvent();
  const { data: eventMe } = useFetchEventMe(selectedEventId);

  const canAbort = eventMe?.capabilities.includes('abort_game') ?? false;
  const canUndo = eventMe?.capabilities.includes('undo_round') ?? false;

  return (
    <>
      <Drawer.Header>
        <Drawer.Heading>Table menu</Drawer.Heading>
      </Drawer.Header>
      <Drawer.Body className="flex flex-col">
        <LazyAccordion
          defaultExpandedKeys={finedPlayerId ? ['penalty'] : []}
          sections={[
            {
              id: 'penalty',
              title: 'Penalty',
              icon: <WarningIcon className="text-warning" />,
              content: () => (
                <PenaltyInput
                  players={players}
                  finedPlayerId={finedPlayerId}
                  setFinedPlayerId={setFinedPlayerId}
                  onPenaltyAdded={onPenaltyAdded}
                />
              ),
            },
            {
              id: 'logs',
              title: 'Logs',
              icon: <ListIcon />,
              content: () => (
                <div className="flex flex-col gap-2">
                  <LogsMenu
                    gameId={gameId}
                    canUndoLastRound={canUndo}
                    onUndoSuccess={onRefetchNow}
                  />
                </div>
              ),
            },
            {
              id: 'fans',
              title: 'Fans',
              icon: <BookIcon />,
              content: () => <FansMenu />,
            },
          ]}
        />

        <div className="mt-auto pt-4 flex justify-center gap-3 w-7/8 mx-auto">
          <FinishMenu gameId={gameId} canFinish />
          <AbortMenu gameId={gameId} canAbort={canAbort} />
        </div>
      </Drawer.Body>
    </>
  );
}
