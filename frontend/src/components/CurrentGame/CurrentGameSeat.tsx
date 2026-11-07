import { Button } from '@heroui/react';
import type { ComponentProps } from 'react';
import type { GameData } from '../../models/GameData';
import AnimatedScore from './AnimatedScore';
import UserAvatar from '../UserAvatar';
import { SeatMarks } from './RoundEntry';

type SeatPlayer = GameData['me'];
type SeatEntry = ComponentProps<typeof SeatMarks>['entry'];

export type CurrentGameSeatProps = {
  player: SeatPlayer;
  entry: SeatEntry;
  isEntryOpen: boolean;
  isSideSeat?: boolean;
  previousScores: Record<number, number>;
  currentScores: Record<number, number>;
  onOpenDelta: (playerId: number) => void;
};

export default function CurrentGameSeat({
  player,
  entry,
  isEntryOpen,
  isSideSeat = false,
  previousScores,
  currentScores,
  onOpenDelta,
}: CurrentGameSeatProps) {
  return (
    <>
      {isEntryOpen && (
        <SeatMarks entry={entry} player={player} isSideSeat={isSideSeat} />
      )}
      <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
        <span className="shrink-0">{player.wind}</span>
        <UserAvatar name={player.name} className="shrink-0" />
        <span className="min-w-0 truncate">{player.name}</span>
      </span>
      <div className="flex items-center justify-center gap-1">
        {!isEntryOpen && (
          <>
            <span className="text-ellipsis overflow-hidden">
              <AnimatedScore
                previousScores={previousScores}
                currentScores={currentScores}
                playerId={player.id}
              />
            </span>
            <Button
              isIconOnly
              variant="tertiary"
              size="sm"
              aria-label={`Show position thresholds for ${player.name}`}
              onPress={() => onOpenDelta(player.id)}
              className="shrink-0"
            >
              <svg
                aria-hidden="true"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
              >
                <path
                  d="M7 20V7M3.5 10.5 7 6l3.5 4.5"
                  stroke="var(--color-success, var(--success))"
                  strokeWidth="2.8"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
                <path
                  d="M17 4v13m-3.5-3.5L17 18l3.5-4.5"
                  stroke="var(--color-danger, var(--danger))"
                  strokeWidth="2.8"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            </Button>
          </>
        )}
      </div>
    </>
  );
}
