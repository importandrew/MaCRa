import type { TFunction } from 'i18next';
import { useTranslation } from 'react-i18next';
import { useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { Outcome, type RoundData } from '../../models/GameData';
import {
  FALSE_MAHJONG_AMOUNTS,
  PENALTY_AMOUNTS,
  STANDARD_RULES,
  type GameRules,
} from '../../models/GameRules';
import {
  fetchWithStatus,
  PostRequest,
  eventGamePath,
  eventPath,
} from '../Request';
import { showToast } from '../Toast';

// START_FINE: the penalty picked when the form opens, the most common one.
const START_FINE = 10;

export type EntryOutcome =
  Outcome.Draw | Outcome.Self | Outcome.Deal | Outcome.Penalty;
export type PenaltyType = Outcome.Penalty | Outcome.FalseMahjong;

// EntryTable is what the entry needs of the game on screen: the player's
// table (GameData) or a game's live view (GameView).
type EntryTable = {
  round: RoundData;
  log_size: number;
  rules: GameRules;
  status: string;
};

export type RoundEntryOptions<T extends EntryTable> = {
  // table is the game as the screen shows it.
  table: T | null;
  eventId: number | null;
  // gameId: enter into this game from its view, as a referee does; without
  // it, into the player's own game at their table.
  gameId?: number;
  // penaltyOnly: only a penalty or a false mahjong (the game's view); the
  // other outcomes are shown but disabled, so the layout stays the same.
  penaltyOnly?: boolean;
  // canPenalize: whether Penalty may be picked (the event's settings may
  // leave penalties to referees); shown disabled when it may not.
  canPenalize: boolean;
  onSubmitted: (data: T) => void;
  // onStale refreshes the table after an entry was refused because the
  // game has moved on since, and returns it (null if it can't be loaded).
  onStale: () => Promise<T | null>;
  // onGone: the game is no longer live (ended, finished or aborted
  // elsewhere); the screen tells the user and moves on.
  onGone: () => void | Promise<unknown>;
  // onEnded: this entry has ended the game.
  onEnded: () => void;
};

// useRoundEntry holds the round entry on a game screen: the outcome and
// values in the center of the table, and the +/− marks next to each seat.
// An entry is refused when the server's game has moved on since the table
// was loaded; onStale then refreshes it.
export function useRoundEntry<T extends EntryTable>({
  table,
  eventId,
  gameId,
  penaltyOnly = false,
  canPenalize,
  onSubmitted,
  onStale,
  onGone,
  onEnded,
}: RoundEntryOptions<T>) {
  const { t } = useTranslation();
  const startOutcome = penaltyOnly ? Outcome.Penalty : null;

  const [outcome, setOutcome] = useState<EntryOutcome | null>(startOutcome);
  const [winner, setWinner] = useState<number | null>(null);
  const [loser, setLoser] = useState<number | null>(null);
  // The game's rules; the standard ones only until the table is loaded.
  const rules = table?.rules ?? STANDARD_RULES;
  const minHand = rules.min_hand;
  // handInput is null until the player changes it: the hand starts at the
  // game's minimum.
  const [handInput, setHandInput] = useState<number | null>(null);
  const hand = handInput ?? minHand;
  const [penaltyType, setPenaltyType] = useState<PenaltyType>(Outcome.Penalty);
  const [fine, setFine] = useState(START_FINE);
  const [falseMahjongFine, setFalseMahjongFine] = useState<number | null>(null);
  const [forceNextRound, setForceNextRound] = useState(false);
  const { isPending, run } = useSubmitOnce();

  const isWin = outcome === Outcome.Self || outcome === Outcome.Deal;
  const isPenalty = outcome === Outcome.Penalty;
  const isFalseMahjong = isPenalty && penaltyType === Outcome.FalseMahjong;
  const showLose = outcome === Outcome.Deal || isPenalty;

  // Winner and loser are kept when the outcome changes (switching Self ⇄
  // Deal-in shouldn't make the user pick again), so check only the ones
  // the chosen outcome uses.
  const isComplete =
    outcome === Outcome.Draw ||
    (outcome === Outcome.Self && winner !== null && hand >= minHand) ||
    (outcome === Outcome.Deal &&
      winner !== null &&
      loser !== null &&
      winner !== loser &&
      hand >= minHand) ||
    (isPenalty &&
      canPenalize &&
      loser !== null &&
      (isFalseMahjong ? !!falseMahjongFine : PENALTY_AMOUNTS.includes(fine)));

  // A Self winner and a Penalty loser can be the same player; Deal-in can't
  // have both on one player, so switching to it drops the − mark.
  const changeOutcome = (next: EntryOutcome) => {
    setOutcome(next);
    if (next === Outcome.Deal && winner !== null && winner === loser) {
      setLoser(null);
    }
  };

  const pickWinner = (id: number) => {
    setWinner(winner === id ? null : id);
    // One player can't both win and deal in.
    if (loser === id && outcome === Outcome.Deal) setLoser(null);
  };

  const pickLoser = (id: number) => {
    setLoser(loser === id ? null : id);
    if (winner === id && outcome === Outcome.Deal) setWinner(null);
  };

  // isOutcomeDisabled: the outcomes this entry can't take, shown disabled.
  const isOutcomeDisabled = (option: EntryOutcome) =>
    option === Outcome.Penalty ? !canPenalize : penaltyOnly;

  const reset = () => {
    setOutcome(startOutcome);
    setWinner(null);
    setLoser(null);
    setHandInput(null);
    setPenaltyType(Outcome.Penalty);
    setFine(START_FINE);
    setFalseMahjongFine(null);
    setForceNextRound(false);
  };

  const submit = () => run(send);

  const send = async () => {
    if (!eventId || !outcome || !table) return;

    const penaltyPath = isFalseMahjong ? '/false-mahjong' : '/penalty';
    const request = isPenalty
      ? PostRequest(
          gameId != null
            ? eventGamePath(eventId, gameId, penaltyPath)
            : eventPath(eventId, penaltyPath),
          {
            loser,
            value: isFalseMahjong ? falseMahjongFine : fine,
            force_next_round: forceNextRound,
            log_size: table.log_size,
          }
        )
      : PostRequest(eventPath(eventId, '/round'), {
          outcome,
          won: isWin ? winner : null,
          dealt: outcome === Outcome.Deal ? loser : null,
          score: isWin ? hand : 0,
          log_size: table.log_size,
        });
    const { data, status } = await fetchWithStatus<T>(request, {
      quietStatuses: [404, 409],
    });

    // 404: the game is no longer live: someone else entered the last round,
    // finished or aborted it. onGone tells the user which.
    if (status === 404) {
      reset();
      await onGone();
      return;
    }

    // 409: someone else has changed the round log since this table was
    // loaded. Don't store this entry on top of it; show the table as it is
    // now and let the player decide whether to enter it again.
    if (status === 409) {
      reset();
      const fresh = await onStale();
      if (!fresh) return;
      showToast({
        title: staleTitle(table, fresh, t),
        description: t('game.checkTableAndRetry'),
        variant: 'warning',
      });
      return;
    }

    if (!data) return;

    if (data.status === 'ended' || data.status === 'aborted') {
      onEnded();
      return;
    }

    if (isPenalty) {
      showToast({ title: t('game.penaltyAdded'), variant: 'success' });
    }
    reset();
    onSubmitted(data);
  };

  return {
    outcome,
    setOutcome: changeOutcome,
    winner,
    loser,
    hand,
    setHand: setHandInput,
    minHand,
    falseMahjongAmounts: FALSE_MAHJONG_AMOUNTS,
    penaltyType,
    setPenaltyType,
    fine,
    setFine,
    falseMahjongFine,
    setFalseMahjongFine,
    forceNextRound,
    setForceNextRound,
    isPending,
    isWin,
    isPenalty,
    isFalseMahjong,
    showLose,
    isComplete,
    isOutcomeDisabled,
    pickWinner,
    pickLoser,
    reset,
    submit,
  };
}

export type RoundEntry = ReturnType<typeof useRoundEntry>;

// staleTitle says what changed between the table the player saw and the
// table now: a round entered, a penalty entered, or an entry undone.
function staleTitle(seen: EntryTable, now: EntryTable, t: TFunction): string {
  if (now.round.num > seen.round.num) {
    return t('game.roundAlreadyEntered', { round: seen.round.kanji });
  }
  if (now.round.num < seen.round.num) {
    return t('game.lastRoundUndone');
  }
  if (now.log_size > seen.log_size) {
    return t('game.penaltyJustEntered');
  }
  return t('game.lastEntryUndone');
}
