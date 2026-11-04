import { Drawer } from '@heroui/react';
import LazyAccordion from '../../LazyAccordion';
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
          <FinishMenu gameId={gameId} canFinish={canFinish} />
          <AbortMenu gameId={gameId} canAbort={canAbort} />
        </div>
      </Drawer.Body>
    </>
  );
}
