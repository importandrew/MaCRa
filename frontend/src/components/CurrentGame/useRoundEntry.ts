import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { GameData, Outcome } from '../../models/GameData';
import { getStoredEventId } from '../../providers/event/eventConstants';
import { fetchWithStatus, PostRequest, eventPath } from '../Request';
import { showToast } from '../Toast';

export const MIN_HAND = 8;
const START_FINE = 10;
export const MIN_FINE = 1;

export type EntryOutcome =
  Outcome.Draw | Outcome.Self | Outcome.Deal | Outcome.Penalty;
export type PenaltyType = Outcome.Penalty | Outcome.FalseMahjong;

// useRoundEntry holds the round entry on the game screen: the outcome and
// values in the center of the table, and the +/− marks next to each seat.
// table is the game as the screen shows it; an entry is refused when the
// server's game has moved on since, and onStale then refreshes the table and
// returns it (null if it can't be loaded, e.g. the game has ended: onStale
// tells the player then).
export function useRoundEntry(
  table: GameData | null,
  onSubmitted: (data: GameData) => void,
  onStale: () => Promise<GameData | null>
) {
  const navigate = useNavigate();

  const [outcome, setOutcome] = useState<EntryOutcome | null>(null);
  const [winner, setWinner] = useState<number | null>(null);
  const [loser, setLoser] = useState<number | null>(null);
  const [hand, setHand] = useState(MIN_HAND);
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
    (outcome === Outcome.Self && winner !== null && hand >= MIN_HAND) ||
    (outcome === Outcome.Deal &&
      winner !== null &&
      loser !== null &&
      winner !== loser &&
      hand >= MIN_HAND) ||
    (isPenalty &&
      loser !== null &&
      (isFalseMahjong ? !!falseMahjongFine : fine >= MIN_FINE));

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

  const reset = () => {
    setOutcome(null);
    setWinner(null);
    setLoser(null);
    setHand(MIN_HAND);
    setPenaltyType(Outcome.Penalty);
    setFine(START_FINE);
    setFalseMahjongFine(null);
    setForceNextRound(false);
  };

  const submit = () => run(send);

  const send = async () => {
    const eventId = getStoredEventId();
    if (!eventId || !outcome || !table) return;

    const request = isPenalty
      ? PostRequest(
          eventPath(eventId, isFalseMahjong ? '/false-mahjong' : '/penalty'),
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
    const { data, status } = await fetchWithStatus<GameData>(request, {
      quietStatuses: [404, 409],
    });

    // 404: the player has no live game any more: someone else entered the
    // last round, finished or aborted it. The refresh tells them which.
    if (status === 404) {
      await onStale();
      return;
    }

    // 409: someone else has changed the round log since this table was
    // loaded. Don't store this entry on top of it; show the table as it is
    // now and let the player decide whether to enter it again.
    if (status === 409) {
      const fresh = await onStale();
      if (!fresh) return;
      showToast({
        title: staleTitle(table, fresh),
        description: 'Check the refreshed table and enter again if needed.',
        variant: 'warning',
      });
      return;
    }

    if (!data) return;

    if (data.status === 'ended' || data.status === 'aborted') {
      navigate('/last');
      return;
    }

    if (isPenalty) {
      showToast({ title: 'Penalty added', variant: 'warning' });
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
    setHand,
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
    pickWinner,
    pickLoser,
    reset,
    submit,
  };
}

export type RoundEntry = ReturnType<typeof useRoundEntry>;

// staleTitle says what changed between the table the player saw and the
// table now: a round entered, a penalty entered, or an entry undone.
function staleTitle(seen: GameData, now: GameData): string {
  if (now.round.num > seen.round.num) {
    return `Round ${seen.round.kanji} has already been entered`;
  }
  if (now.round.num < seen.round.num) {
    return 'The last round has been undone';
  }
  if (now.log_size > seen.log_size) {
    return 'A penalty has just been entered';
  }
  return 'The last entry has been undone';
}
