import { Switch, ToggleButton, ToggleButtonGroup } from '@heroui/react';
import { Outcome, PlayerScore } from '../../models/GameData';
import { OUTCOME_LABELS } from '../../utils/constants';
import PendingButton from '../PendingButton';
import SegmentedPicker from '../SegmentedPicker';
import ValueStepper from '../ValueStepper';
import { EntryOutcome, MIN_HAND, RoundEntry } from './useRoundEntry';

const OUTCOME_OPTIONS = [
  ...[Outcome.Draw, Outcome.Self, Outcome.Deal].map((value, index) => ({
    value,
    // Non-breaking hyphen: "Deal-in" must not wrap in the narrow center.
    label: OUTCOME_LABELS[index].replace('-', '‑'),
  })),
  { value: Outcome.Penalty, label: 'Penalty' },
] as { value: EntryOutcome; label: string }[];

// SeatMarks are the + (won) and − (dealt in / fined) toggles at a player's
// seat while a round is entered. They sit between the name and the center
// of the table; the score line is hidden meanwhile, so the center keeps its
// room. + and − read the
// same when the seat is rotated, unlike 👍/👎.
export function SeatMarks({
  entry,
  player,
  isSideSeat = false,
}: {
  entry: RoundEntry;
  player: PlayerScore;
  // isSideSeat turns the − back on the left and right seats: rotated 90°
  // with the seat it would look like a vertical bar ("I").
  isSideSeat?: boolean;
}) {
  return (
    // Fixed height: the line keeps its size whatever the outcome shows.
    <span className="flex h-10 items-center justify-center gap-2">
      {entry.isWin && (
        <ToggleButton
          isIconOnly
          size="sm"
          aria-label={`${player.name} won`}
          isSelected={entry.winner === player.id}
          onChange={() => entry.pickWinner(player.id)}
          className="size-10 text-2xl"
        >
          +
        </ToggleButton>
      )}
      {entry.showLose && (
        <ToggleButton
          isIconOnly
          size="sm"
          aria-label={`${player.name} ${entry.isPenalty ? 'fined' : 'dealt in'}`}
          isSelected={entry.loser === player.id}
          onChange={() => entry.pickLoser(player.id)}
          className="size-10 text-2xl"
        >
          <span className={isSideSeat ? 'rotate-90' : undefined}>
            {'\u2212'}
          </span>
        </ToggleButton>
      )}
    </span>
  );
}

// RoundEntryPanel is the center of the table while a round is entered:
// outcome, hand value or fine, and OK. Players are marked on their seats
// (SeatMarks), so no names are picked from a list. The round kanji above it
// closes it, so there's no Cancel.
export function RoundEntryPanel({ entry }: { entry: RoundEntry }) {
  return (
    <div className="flex w-full flex-col items-stretch gap-2 text-base">
      {/* 2×2 instead of one row: four options don't fit the center of a
          phone in one line. */}
      <ToggleButtonGroup
        aria-label="Outcome"
        selectionMode="single"
        disallowEmptySelection
        isDetached
        className="grid w-full grid-cols-2 gap-1.5"
        selectedKeys={entry.outcome ? [entry.outcome] : []}
        onSelectionChange={(keys) => {
          const [key] = keys;
          if (key != null) entry.setOutcome(String(key) as EntryOutcome);
        }}
      >
        {OUTCOME_OPTIONS.map((opt) => (
          <ToggleButton
            key={opt.value}
            id={opt.value}
            size="sm"
            className="w-full min-w-0 not-data-selected:text-foreground"
          >
            {opt.label}
          </ToggleButton>
        ))}
      </ToggleButtonGroup>

      {entry.isPenalty && (
        <SegmentedPicker
          ariaLabel="Penalty type"
          options={[
            { value: Outcome.Penalty, label: 'Penalty' },
            { value: Outcome.FalseMahjong, label: 'False mahjong' },
          ]}
          selectedKey={entry.penaltyType}
          onSelectionChange={entry.setPenaltyType}
        />
      )}

      {entry.isWin && (
        <ValueStepper
          label="Hand value"
          value={entry.hand}
          min={MIN_HAND}
          onChange={entry.setHand}
        />
      )}

      {entry.isPenalty && !entry.isFalseMahjong && (
        <ValueStepper
          label="Fine"
          value={entry.fine}
          min={1}
          onChange={entry.setFine}
        />
      )}

      {entry.isFalseMahjong && (
        <SegmentedPicker
          ariaLabel="Fine"
          options={[
            { value: '30', label: '30' },
            { value: '60', label: '60' },
          ]}
          selectedKey={
            entry.falseMahjongFine ? String(entry.falseMahjongFine) : null
          }
          onSelectionChange={(value) =>
            entry.setFalseMahjongFine(Number(value))
          }
        />
      )}

      {entry.isPenalty && (
        <Switch
          size="sm"
          isSelected={entry.forceNextRound}
          onChange={entry.setForceNextRound}
          className="self-center"
        >
          <Switch.Content>
            <Switch.Control>
              <Switch.Thumb />
            </Switch.Control>
            Force next round
          </Switch.Content>
        </Switch>
      )}

      <div className="flex justify-center">
        <PendingButton
          size="sm"
          variant="primary"
          isDisabled={!entry.isComplete}
          isPending={entry.isPending}
          onPress={entry.submit}
        >
          OK
        </PendingButton>
      </div>
    </div>
  );
}
