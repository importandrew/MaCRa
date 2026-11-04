import { Button, Switch, ToggleButton, ToggleButtonGroup } from '@heroui/react';
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

// SeatMarks are the 👍 (won) and 👎 (dealt in / fined) toggles next to a
// player's name on the table, shown only while the outcome needs them.
export function SeatMarks({
  entry,
  player,
}: {
  entry: RoundEntry;
  player: PlayerScore;
}) {
  if (!entry.isWin && !entry.showLose) return null;

  return (
    <span className="flex shrink-0 gap-1.5">
      {entry.isWin && (
        <ToggleButton
          isIconOnly
          size="sm"
          aria-label={`${player.name} won`}
          isSelected={entry.winner === player.id}
          onChange={() => entry.pickWinner(player.id)}
          className="size-10 text-lg"
        >
          👍
        </ToggleButton>
      )}
      {entry.showLose && (
        <ToggleButton
          isIconOnly
          size="sm"
          aria-label={`${player.name} ${entry.isPenalty ? 'fined' : 'dealt in'}`}
          isSelected={entry.loser === player.id}
          onChange={() => entry.pickLoser(player.id)}
          className="size-10 text-lg"
        >
          👎
        </ToggleButton>
      )}
    </span>
  );
}

// RoundEntryPanel is the center of the table while a round is entered:
// outcome, hand value or fine, and OK/Cancel. Players are marked on their
// seats (SeatMarks), so no names are picked from a list.
export function RoundEntryPanel({
  entry,
  onCancel,
}: {
  entry: RoundEntry;
  onCancel: () => void;
}) {
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

      <div className="flex justify-center gap-2">
        <Button size="sm" variant="secondary" onPress={onCancel}>
          Cancel
        </Button>
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
