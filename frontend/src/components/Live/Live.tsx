import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useRouteIds } from '../../hooks/useRouteIds';
import { useFetchLiveGameData } from '../../providers/game/gameDataApi';
import { useGameCapabilities } from '../../providers/game/useGameCapabilities';
import type { PlayerScore } from '../../models/GameData';
import { RoundEntryPanel, SeatMarks } from '../CurrentGame/RoundEntry';
import { type RoundEntry, useRoundEntry } from '../CurrentGame/useRoundEntry';
import ContextMenu from '../ContextMenu/ContextMenu';
import RefreshButton from '../RefreshButton';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import AnimatedScore from '../Table/AnimatedScore';
import SeatName from '../Table/SeatName';
import TableLayout from '../Table/TableLayout';
import { useScoreAnimation } from '../Table/useScoreAnimation';
import LiveMenu from './Menu/LiveMenu';

// LiveSeat is a seat of the live view: the name and the score, or, while a
// penalty is entered, the − mark to pick the penalized player instead of
// the score (as at the table).
function LiveSeat({
  player,
  entry,
  isEntryOpen,
  isSideSeat = false,
  previousScores,
  currentScores,
}: {
  player: PlayerScore;
  entry: RoundEntry;
  isEntryOpen: boolean;
  isSideSeat?: boolean;
  previousScores: Record<number, number>;
  currentScores: Record<number, number>;
}) {
  return (
    <>
      {isEntryOpen && (
        <SeatMarks entry={entry} player={player} isSideSeat={isSideSeat} />
      )}
      <SeatName player={player} />
      <span
        className={`text-ellipsis overflow-hidden ${isEntryOpen ? 'hidden' : ''}`}
      >
        <AnimatedScore
          previousScores={previousScores}
          currentScores={currentScores}
          playerId={player.id}
        />
      </span>
    </>
  );
}

// onGameEnded is called once this user has ended the game (Finish, or a
// penalty that ended it), or found it no longer live: the game page reloads
// its status and shows the result.
export default function Live({ onGameEnded }: { onGameEnded: () => void }) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { eventId, gameId } = useRouteIds();
  const { data, isLoading, isError, isForbidden, refetch, mutate } =
    useFetchLiveGameData(eventId, gameId);
  // A referee (or a seated player, as the event's settings allow) enters
  // penalties here: the round kanji opens the entry, penalties only.
  const { hasCapability } = useGameCapabilities(eventId, gameId);
  const canEnterPenalty = hasCapability('enter_penalty');
  const [isEntryOpen, setIsEntryOpen] = useState(false);
  const entry = useRoundEntry({
    table: data,
    eventId,
    gameId: gameId ?? undefined,
    penaltyOnly: true,
    canPenalize: true,
    onSubmitted: (view) => {
      mutate(view);
      setIsEntryOpen(false);
    },
    onStale: refetch,
    onGone: onGameEnded,
    onEnded: onGameEnded,
  });
  const [isContextMenuOpen, setIsContextMenuOpen] = useState(false);
  const [rotationOffset, setRotationOffset] = useState(0);
  const { previousScores, currentScores } = useScoreAnimation(
    data?.player_scores ?? null
  );

  const closeContextMenuOpen = () => {
    setIsContextMenuOpen(false);
  };

  // Refresh clears the entry, as at the table.
  const onRefresh = () => {
    entry.reset();
    setIsEntryOpen(false);
    return refetch();
  };

  const getRotatedIndex = (index: number): number => {
    const rotatedIndex = (index - rotationOffset + 4) % 4;
    return rotatedIndex;
  };

  const rotateCounterClockwise = () => {
    setRotationOffset((prev) => (prev + 1 + 4) % 4);
  };

  const rotateClockwise = () => {
    setRotationOffset((prev) => (prev - 1 + 4) % 4);
  };

  const bot = data && data.player_scores[getRotatedIndex(0)];
  const right = data && data.player_scores[getRotatedIndex(1)];
  const top = data && data.player_scores[getRotatedIndex(2)];
  const left = data && data.player_scores[getRotatedIndex(3)];

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      isForbidden={isForbidden}
      subject="game"
      onRefresh={() => void onRefresh()}
      onGoHome={() => navigate('/')}
    >
      {data && eventId && bot && right && top && left ? (
        // overflow-x-clip: while a penalty is entered, a side seat with its −
        // mark sticks out a few px on a 320px phone, as at the table
        // (CurrentGame).
        <TableLayout
          className="overflow-x-clip"
          top={
            <LiveSeat
              player={top}
              entry={entry}
              isEntryOpen={isEntryOpen}
              previousScores={previousScores}
              currentScores={currentScores}
            />
          }
          left={
            <LiveSeat
              player={left}
              entry={entry}
              isEntryOpen={isEntryOpen}
              isSideSeat
              previousScores={previousScores}
              currentScores={currentScores}
            />
          }
          right={
            <LiveSeat
              player={right}
              entry={entry}
              isEntryOpen={isEntryOpen}
              isSideSeat
              previousScores={previousScores}
              currentScores={currentScores}
            />
          }
          bottom={
            <LiveSeat
              player={bot}
              entry={entry}
              isEntryOpen={isEntryOpen}
              previousScores={previousScores}
              currentScores={currentScores}
            />
          }
          center={
            <div className="relative z-10 py-3 justify-self-center self-center items-center w-full text-5xl">
              <div className="flex flex-col items-center gap-2">
                {canEnterPenalty ? (
                  <Button
                    size="lg"
                    variant="secondary"
                    className={
                      isEntryOpen ? 'h-11 px-3 text-3xl' : 'p-4 h-16 text-5xl'
                    }
                    aria-label={t('game.enterPenalty', {
                      round: data.round.kanji,
                    })}
                    aria-expanded={isEntryOpen}
                    onPress={() => setIsEntryOpen(!isEntryOpen)}
                  >
                    {data.round.kanji}
                  </Button>
                ) : (
                  <div className="p-4 h-16 text-5xl">{data.round.kanji}</div>
                )}
                {isEntryOpen && (
                  <div className="w-full px-1 text-base">
                    <RoundEntryPanel entry={entry} />
                  </div>
                )}
                <div className={`flex gap-2 ${isEntryOpen ? 'hidden' : ''}`}>
                  <Button
                    size="sm"
                    variant="tertiary"
                    type="button"
                    aria-label={t('game.rotateClockwise')}
                    onPress={rotateClockwise}
                    className="text-xl"
                  >
                    ↻
                  </Button>
                  <Button
                    size="sm"
                    variant="tertiary"
                    type="button"
                    aria-label={t('game.rotateCounterclockwise')}
                    onPress={rotateCounterClockwise}
                    className="text-xl"
                  >
                    ↺
                  </Button>
                </div>
              </div>
            </div>
          }
          refresh={
            <RefreshButton
              onRefresh={onRefresh}
              className="absolute z-1 left-0 bottom-0"
            />
          }
          menu={
            <ContextMenu
              isOpen={isContextMenuOpen}
              onOpen={() => setIsContextMenuOpen(true)}
              onClose={closeContextMenuOpen}
            >
              <LiveMenu
                gameId={data.id}
                eventId={eventId}
                onUndoSuccess={refetch}
                onStale={async () => (await refetch())?.log ?? null}
                log={data.log}
                onGameEnded={onGameEnded}
              />
            </ContextMenu>
          }
        />
      ) : null}
    </LoadingTemplate>
  );
}
