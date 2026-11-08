import { Drawer } from '@heroui/react';
import { useNavigate } from 'react-router-dom';
import { useFetchEventMe } from '../../../providers/event/eventApi';
import { useEvent } from '../../../providers/event/eventConstants';
import { BookIcon } from '../../Icons/BookIcon';
import { ListIcon } from '../../Icons/ListIcon';
import LogsMenu from './LogsMenu';
import AbortMenu from './AbortMenu';
import FinishMenu from './FinishMenu';
import { FansMenu } from './FansMenu';
import LazyAccordion from '../../LazyAccordion';
import { GameLogsState } from '../../../providers/game/gameDataApi';

// TableMenu is the ≡ menu of the game screen. Penalties are entered in the
// round card (tap the round kanji), not here.
export default function TableMenu({
  gameId = null,
  onRefetchNow,
  logs,
  logSize,
}: {
  gameId: number | null;
  // onRefetchNow reloads the table; it is also called after an undo refused
  // because the log has changed.
  onRefetchNow?: () => Promise<unknown>;
  // logs and logSize are the table's, see useTableLogs.
  logs: GameLogsState;
  logSize: number;
}) {
  const navigate = useNavigate();
  const { selectedEventId, selectedEvent } = useEvent();
  const { data: eventMe } = useFetchEventMe(selectedEventId);

  const canAbort = eventMe?.capabilities.includes('abort_game') ?? false;
  const canUndo = eventMe?.capabilities.includes('undo_round') ?? false;
  // The event-level finish_game is computed without a game, so it leaves out
  // seated players. At the table the user is always seated, and a seated
  // player may finish in a regular event (as on the backend). TODO 25 replaces
  // this with game-level capabilities.
  const canFinish = selectedEvent?.type === 'regular';

  return (
    <>
      <Drawer.Header>
        <Drawer.Heading>Table menu</Drawer.Heading>
      </Drawer.Header>
      <Drawer.Body className="flex flex-col">
        <LazyAccordion
          sections={[
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
                    onStale={onRefetchNow}
                    tableLogs={logs}
                    tableLogSize={logSize}
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
          <FinishMenu
            gameId={gameId}
            eventId={selectedEventId}
            canFinish={canFinish}
            onFinished={() => navigate('/last')}
          />
          <AbortMenu
            gameId={gameId}
            eventId={selectedEventId}
            canAbort={canAbort}
          />
        </div>
      </Drawer.Body>
    </>
  );
}
