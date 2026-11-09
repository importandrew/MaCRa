import { Drawer } from '@heroui/react';
import { useNavigate } from 'react-router-dom';
import LazyAccordion from '../../LazyAccordion';
import { useEventCapabilities } from '../../../providers/event/useEventCapabilities';
import { BookIcon } from '../../Icons/BookIcon';
import { ListIcon } from '../../Icons/ListIcon';
import LogsMenu from '../../CurrentGame/Menu/LogsMenu';
import { GameLogsState } from '../../../providers/game/gameDataApi';
import AbortMenu from '../../CurrentGame/Menu/AbortMenu';
import FinishMenu from '../../CurrentGame/Menu/FinishMenu';
import { FansMenu } from '../../CurrentGame/Menu/FansMenu';

export default function LiveMenu({
  gameId = null,
  eventId,
  onUndoSuccess,
  logs,
  logSize,
  onGameEnded,
}: {
  gameId: number | null;
  eventId?: number;
  // onUndoSuccess reloads the screen; it is also called after an undo
  // refused because the log has changed.
  onUndoSuccess?: () => Promise<unknown>;
  // logs and logSize are the screen's, see useTableLogs.
  logs: GameLogsState;
  logSize: number;
  // onGameEnded reloads the game page, which then shows the result.
  onGameEnded: () => void;
}) {
  const navigate = useNavigate();
  const { hasCapability } = useEventCapabilities(eventId ?? null);
  const canUndoLastRound = hasCapability('undo_round');
  const canAbort = hasCapability('abort_game');
  const canFinish = hasCapability('finish_game');

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
                <LogsMenu
                  gameId={gameId}
                  eventId={eventId}
                  canUndoLastRound={canUndoLastRound}
                  onUndoSuccess={onUndoSuccess}
                  onStale={onUndoSuccess}
                  tableLogs={logs}
                  tableLogSize={logSize}
                />
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
            eventId={eventId ?? null}
            canFinish={canFinish}
            onFinished={onGameEnded}
          />
          <AbortMenu
            gameId={gameId}
            eventId={eventId ?? null}
            canAbort={canAbort}
            onAborted={() =>
              navigate(eventId ? `/event/${eventId}/games/live` : '/')
            }
          />
        </div>
      </Drawer.Body>
    </>
  );
}
