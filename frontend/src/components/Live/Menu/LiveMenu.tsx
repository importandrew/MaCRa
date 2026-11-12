import { useTranslation } from 'react-i18next';
import { Drawer } from '@heroui/react';
import { useNavigate } from 'react-router-dom';
import LazyAccordion from '../../LazyAccordion';
import { useEventCapabilities } from '../../../providers/event/useEventCapabilities';
import { BookIcon } from '../../Icons/BookIcon';
import { ListIcon } from '../../Icons/ListIcon';
import LogsMenu from '../../GameMenu/LogsMenu';
import { GameRoundLogs } from '../../../models/GameResult';
import AbortMenu from '../../GameMenu/AbortMenu';
import FinishMenu from '../../GameMenu/FinishMenu';
import { FansMenu } from '../../GameMenu/FansMenu';

export default function LiveMenu({
  gameId = null,
  eventId,
  onUndoSuccess,
  onStale,
  log,
  onGameEnded,
}: {
  gameId: number | null;
  eventId: number;
  // onUndoSuccess reloads the screen after an undo.
  onUndoSuccess?: () => unknown;
  // onStale reloads the screen after an undo refused because the log has
  // changed, and returns its fresh log.
  onStale: () => Promise<GameRoundLogs | null>;
  // log is the round log from the screen's table response.
  log: GameRoundLogs;
  // onGameEnded reloads the game page, which then shows the result.
  onGameEnded: () => void;
}) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { hasCapability } = useEventCapabilities(eventId);
  const canUndoLastRound = hasCapability('undo_round');
  const canAbort = hasCapability('abort_game');
  const canFinish = hasCapability('finish_game');

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
                <LogsMenu
                  gameId={gameId}
                  eventId={eventId}
                  canUndoLastRound={canUndoLastRound}
                  onUndoSuccess={onUndoSuccess}
                  onStale={onStale}
                  logs={{ data: log }}
                />
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
            eventId={eventId}
            canFinish={canFinish}
            onFinished={onGameEnded}
          />
          <AbortMenu
            gameId={gameId}
            eventId={eventId}
            canAbort={canAbort}
            onAborted={() => navigate(`/event/${eventId}/games/live`)}
          />
        </div>
      </Drawer.Body>
    </>
  );
}
