import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { GameData, Outcome } from '../../models/GameData';
import { getStoredEventId } from '../../providers/event/eventConstants';
import { authFetch, PostRequest, eventPath } from '../Request';
import { showToast } from '../Toast';

export const MIN_HAND = 8;
const START_FINE = 10;

export type EntryOutcome =
  Outcome.Draw | Outcome.Self | Outcome.Deal | Outcome.Penalty;
export type PenaltyType = Outcome.Penalty | Outcome.FalseMahjong;

// useRoundEntry holds the round entry on the game screen: the outcome and
// values in the center of the table, and the +/− marks next to each seat.
export function useRoundEntry(onSubmitted: (data: GameData) => void) {
  const navigate = useNavigate();

  const [outcome, setOutcome] = useState<EntryOutcome | null>(null);
  const [winner, setWinner] = useState<number | null>(null);
  const [loser, setLoser] = useState<number | null>(null);
  const [hand, setHand] = useState(MIN_HAND);
  const [penaltyType, setPenaltyType] = useState<PenaltyType>(Outcome.Penalty);
  const [fine, setFine] = useState(START_FINE);
  const [falseMahjongFine, setFalseMahjongFine] = useState<number | null>(null);
  const [forceNextRound, setForceNextRound] = useState(false);
  const [isPending, setIsPending] = useState(false);

  const isWin = outcome === Outcome.Self || outcome === Outcome.Deal;
  const isPenalty = outcome === Outcome.Penalty;
  const isFalseMahjong = isPenalty && penaltyType === Outcome.FalseMahjong;
  const showLose = outcome === Outcome.Deal || isPenalty;

  // Winner and loser are kept when the outcome changes (switching Self ⇄
  // Deal-in shouldn't make the user pick again), so check only the ones
  // the chosen outcome uses.
  const isComplete =
    outcome === Outcome.Draw ||
    (outcome === Outcome.Self && winner !== null) ||
    (outcome === Outcome.Deal &&
      winner !== null &&
      loser !== null &&
      winner !== loser) ||
    (isPenalty && loser !== null && (!isFalseMahjong || !!falseMahjongFine));

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

  const submit = async () => {
    const eventId = getStoredEventId();
    if (!eventId || !outcome) return;

    setIsPending(true);
    const request = isPenalty
      ? PostRequest(
          eventPath(eventId, isFalseMahjong ? '/false-mahjong' : '/penalty'),
          {
            loser,
            value: isFalseMahjong ? falseMahjongFine : fine,
            force_next_round: forceNextRound,
          }
        )
      : PostRequest(eventPath(eventId, '/round'), {
          outcome,
          won: isWin ? winner : null,
          dealt: outcome === Outcome.Deal ? loser : null,
          score: isWin ? hand : 0,
        });
    const data = await authFetch<GameData>(request, { isToastEnabled: true });
    setIsPending(false);

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
    setOutcome,
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
