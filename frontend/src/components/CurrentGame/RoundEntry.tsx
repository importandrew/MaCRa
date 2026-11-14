import { useTranslation } from 'react-i18next';
import {
  Separator,
  Switch,
  ToggleButton,
  ToggleButtonGroup,
} from '@heroui/react';
import { Outcome, PlayerScore } from '../../models/GameData';

import { SignIcon } from '../Icons/SignIcon';
import PendingButton from '../PendingButton';
import SegmentedPicker from '../SegmentedPicker';
import ValueStepper from '../ValueStepper';
import ListStepper from '../ListStepper';
import { PENALTY_AMOUNTS } from '../../models/GameRules';
import { EntryOutcome, RoundEntry } from './useRoundEntry';

// OUTCOMES are the outcome buttons with their texts. The texts must not wrap
// in the narrow center of the table: non-breaking hyphens and spaces there.
const OUTCOMES = [
  { value: Outcome.Draw, labelKey: 'roundEntry.draw' },
  { value: Outcome.Self, labelKey: 'roundEntry.self' },
  { value: Outcome.Deal, labelKey: 'roundEntry.dealIn' },
  { value: Outcome.Penalty, labelKey: 'roundEntry.penalty' },
] as const satisfies readonly { value: EntryOutcome; labelKey: string }[];

// ToggleButton has no bordered variant, so ghost plus a border: an outline
// that fills (and turns the border accent) when picked. Larger from md up,
// for the mouse; HeroUI's own sizes shrink there instead.
const MARK_CLASS =
  'size-10 md:size-12 rounded-lg border-2 border-border data-[selected=true]:border-accent';

const MARK_ICON_CLASS = 'm-0 size-6 md:size-7 [&_path]:stroke-2';

// Outcome buttons, step buttons and OK grow from md up too.
const DESKTOP_BUTTON_CLASS = 'md:h-10 md:text-base';

// SeatMarks are the + (won) and − (dealt in / fined) toggles at a player's
// seat while a round is entered. They sit between the name and the center
// of the table; the score line is hidden meanwhile, so the center keeps its
// room. Rounded squares, not circles, so they don't look like the round
// avatars next to them; mb-2 keeps them off the name line. + and − read the
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
  const { t } = useTranslation();
  return (
    // Fixed height: the line keeps its size whatever the outcome shows.
    <span className="mb-2 flex h-10 md:h-12 items-center justify-center gap-2">
      {entry.isWin && (
        <SeatMark
          sign="plus"
          label={t('roundEntry.won', { name: player.name })}
          isSelected={entry.winner === player.id}
          onChange={() => entry.pickWinner(player.id)}
        />
      )}
      {entry.showLose && (
        <SeatMark
          sign="minus"
          label={
            entry.isPenalty
              ? t('roundEntry.fined', { name: player.name })
              : t('roundEntry.dealtIn', { name: player.name })
          }
          isSelected={entry.loser === player.id}
          onChange={() => entry.pickLoser(player.id)}
          iconClassName={isSideSeat ? 'rotate-90' : undefined}
        />
      )}
    </span>
  );
}

// SeatMark is one + or − toggle of SeatMarks.
function SeatMark({
  sign,
  label,
  isSelected,
  onChange,
  iconClassName = '',
}: {
  sign: 'plus' | 'minus';
  label: string;
  isSelected: boolean;
  onChange: () => void;
  iconClassName?: string;
}) {
  return (
    <ToggleButton
      isIconOnly
      size="sm"
      variant="ghost"
      aria-label={label}
      isSelected={isSelected}
      onChange={onChange}
      className={MARK_CLASS}
    >
      <SignIcon sign={sign} className={`${MARK_ICON_CLASS} ${iconClassName}`} />
    </ToggleButton>
  );
}

// RoundEntryPanel is the center of the table while a round is entered:
// outcome, hand value or fine, and OK. Players are marked on their seats
// (SeatMarks), so no names are picked from a list. The round kanji above it
// closes it, so there's no Cancel.
export function RoundEntryPanel({ entry }: { entry: RoundEntry }) {
  const { t } = useTranslation();
  return (
    <div className="flex w-full flex-col items-stretch gap-2 text-base">
      {/* 2×2 instead of one row: four options don't fit the center of a
          phone in one line. */}
      <ToggleButtonGroup
        aria-label={t('roundEntry.outcome')}
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
        {OUTCOMES.map((opt) => (
          <ToggleButton
            key={opt.value}
            id={opt.value}
            isDisabled={entry.isOutcomeDisabled(opt.value)}
            size="sm"
            className={`w-full min-w-0 not-data-selected:text-foreground ${DESKTOP_BUTTON_CLASS}`}
          >
            {t(opt.labelKey)}
          </ToggleButton>
        ))}
      </ToggleButtonGroup>

      {/* Sets the outcome apart from what it asks for next, so Penalty /
          False mahjong doesn't read as more outcome buttons. */}
      <Separator />

      {entry.isPenalty && (
        <SegmentedPicker
          ariaLabel={t('roundEntry.penaltyType')}
          // "False mahjong" must fit the center of a 320px phone: smaller
          // text, less padding, and options sized by their text (flex-auto)
          // rather than equal halves. Desktop keeps the usual size.
          optionClassName="flex-auto px-2 text-xs md:px-4 md:text-sm"
          options={[
            { value: Outcome.Penalty, label: t('roundEntry.penalty') },
            // The text has a non-breaking space: on two lines it makes the
            // Penalty form too tall for the center of a phone.
            {
              value: Outcome.FalseMahjong,
              label: t('roundEntry.falseMahjong'),
            },
          ]}
          selectedKey={entry.penaltyType}
          onSelectionChange={entry.setPenaltyType}
        />
      )}

      {entry.isWin && (
        <ValueStepper
          label={t('roundEntry.handValue')}
          value={entry.hand}
          min={entry.minHand}
          onChange={entry.setHand}
        />
      )}

      {entry.isPenalty && !entry.isFalseMahjong && (
        <ListStepper
          label={t('roundEntry.fine')}
          values={PENALTY_AMOUNTS}
          value={entry.fine}
          onChange={entry.setFine}
        />
      )}

      {entry.isFalseMahjong && (
        <SegmentedPicker
          ariaLabel={t('roundEntry.fine')}
          options={entry.falseMahjongAmounts.map((amount) => ({
            value: String(amount),
            label: String(amount),
          }))}
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
            {t('roundEntry.forceNextRound')}
          </Switch.Content>
        </Switch>
      )}

      <div className="flex justify-center">
        <PendingButton
          size="sm"
          variant="primary"
          className={`md:px-6 ${DESKTOP_BUTTON_CLASS}`}
          isDisabled={!entry.isComplete}
          isPending={entry.isPending}
          onPress={entry.submit}
        >
          {t('roundEntry.ok')}
        </PendingButton>
      </div>
    </div>
  );
}
