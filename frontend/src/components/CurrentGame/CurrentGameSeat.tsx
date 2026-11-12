import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';
import type { ComponentProps } from 'react';
import type { GameData } from '../../models/GameData';
import AnimatedScore from '../Table/AnimatedScore';
import SeatName from '../Table/SeatName';
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
  const { t } = useTranslation();
  return (
    <>
      {isEntryOpen && (
        <SeatMarks entry={entry} player={player} isSideSeat={isSideSeat} />
      )}
      <SeatName player={player} />
      <div className="flex items-center justify-center gap-1">
        <span
          className={`text-ellipsis overflow-hidden ${isEntryOpen ? 'hidden' : ''}`}
        >
          <AnimatedScore
            previousScores={previousScores}
            currentScores={currentScores}
            playerId={player.id}
          />
        </span>
        {!isEntryOpen && (
          <Button
            isIconOnly
            variant="tertiary"
            size="sm"
            aria-label={t('delta.showFor', { name: player.name })}
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
        )}
      </div>
    </>
  );
}
