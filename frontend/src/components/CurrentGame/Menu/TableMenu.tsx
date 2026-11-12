import { useTranslation } from 'react-i18next';
import { Drawer } from '@heroui/react';
import { useNavigate } from 'react-router-dom';
import { useEventCapabilities } from '../../../providers/event/useEventCapabilities';
import { useEvent } from '../../../providers/event/eventConstants';
import { BookIcon } from '../../Icons/BookIcon';
import { ListIcon } from '../../Icons/ListIcon';
import LogsMenu from '../../GameMenu/LogsMenu';
import AbortMenu from '../../GameMenu/AbortMenu';
import FinishMenu from '../../GameMenu/FinishMenu';
import { FansMenu } from '../../GameMenu/FansMenu';
import LazyAccordion from '../../LazyAccordion';
import { GameRoundLogs } from '../../../models/GameResult';

// TableMenu is the ≡ menu of the game screen. Penalties are entered in the
// round card (tap the round kanji), not here.
export default function TableMenu({
  gameId = null,
  onRefetchNow,
  onStale,
  log,
}: {
  gameId: number | null;
  // onRefetchNow reloads the table after an undo.
  onRefetchNow?: () => unknown;
  // onStale reloads the table after an undo refused because the log has
  // changed, and returns its fresh log.
  onStale: () => Promise<GameRoundLogs | null>;
  // log is the round log from the table response.
  log: GameRoundLogs;
}) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { selectedEventId, selectedEvent } = useEvent();
  const { hasCapability } = useEventCapabilities(selectedEventId);

  const canAbort = hasCapability('abort_game');
  const canUndo = hasCapability('undo_round');
  // The event-level finish_game is computed without a game, so it leaves out
  // seated players. At the table the user is always seated, and a seated
  // player may finish in a regular event (as on the backend). Capabilities
  // computed for a given game should replace this check.
  const canFinish = selectedEvent?.type === 'regular';

  return (
    <>
      <Drawer.Header>
        <Drawer.Heading>{t('game.tableMenu')}</Drawer.Heading>
      </Drawer.Header>
      <Drawer.Body className="flex flex-col">
        <LazyAccordion
          sections={[
            {
              id: 'logs',
              title: t('game.logs'),
              icon: <ListIcon />,
              content: () => (
                <div className="flex flex-col gap-2">
                  <LogsMenu
                    gameId={gameId}
                    eventId={selectedEventId}
                    canUndoLastRound={canUndo}
                    onUndoSuccess={onRefetchNow}
                    onStale={onStale}
                    logs={{ data: log }}
                  />
                </div>
              ),
            },
            {
              id: 'fans',
              title: t('game.fans'),
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
            onAborted={() => navigate('/')}
          />
        </div>
      </Drawer.Body>
    </>
  );
}
