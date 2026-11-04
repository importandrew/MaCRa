import {
  Button,
  Modal,
  Switch,
  ToggleButton,
  useMediaQuery,
} from '@heroui/react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { GameData, Outcome, PlayerScore } from '../../models/GameData';
import { OUTCOME_LABELS } from '../../utils/constants';
import { getStoredEventId } from '../../providers/event/eventConstants';
import { authFetch, PostRequest, eventPath } from '../Request';
import { showToast } from '../Toast';
import PendingButton from '../PendingButton';
import SegmentedPicker from '../SegmentedPicker';
import ValueStepper from '../ValueStepper';

const MIN_HAND = 8;
const START_FINE = 10;

type CardOutcome = Outcome.Draw | Outcome.Self | Outcome.Deal | Outcome.Penalty;
type PenaltyType = Outcome.Penalty | Outcome.FalseMahjong;

const OUTCOME_OPTIONS = [
  ...[Outcome.Draw, Outcome.Self, Outcome.Deal].map((value, index) => ({
    value,
    // Non-breaking hyphen: "Deal-in" must not wrap on narrow phones.
    label: OUTCOME_LABELS[index].replace('-', '‑'),
  })),
  { value: Outcome.Penalty, label: 'Penalty' },
] as { value: CardOutcome; label: string }[];

// Seat shows one player where they sit at the table, with the 👍 (won) and
// 👎 (dealt in / fined) toggles the chosen outcome allows.
function Seat({
  player,
  showWin,
  showLose,
  isWinner,
  isLoser,
  loseLabel,
  onWin,
  onLose,
}: {
  player: PlayerScore;
  showWin: boolean;
  showLose: boolean;
  isWinner: boolean;
  isLoser: boolean;
  loseLabel: string;
  onWin: () => void;
  onLose: () => void;
}) {
  return (
    <div className="flex min-w-0 flex-col items-center gap-1.5 text-center">
      <span className="text-2xl leading-none">{player.wind}</span>
      <span className="line-clamp-2 w-full break-words text-sm leading-tight">
        {player.name}
      </span>
      {/* Fixed height: the seat keeps its size when the toggles change. */}
      <div className="flex h-11 items-center gap-1.5">
        {showWin && (
          <ToggleButton
            isIconOnly
            aria-label={`${player.name} won`}
            isSelected={isWinner}
            onChange={onWin}
            className="size-11 text-xl"
          >
            👍
          </ToggleButton>
        )}
        {showLose && (
          <ToggleButton
            isIconOnly
            aria-label={`${player.name} ${loseLabel}`}
            isSelected={isLoser}
            onChange={onLose}
            className="size-11 text-xl"
          >
            👎
          </ToggleButton>
        )}
      </div>
    </div>
  );
}

// RoundCard enters a round result or a penalty. Players are picked by their
// seat at the table, and values with steppers, so no keyboard is needed.
// The parent keys it to reset it (after a submit and on refresh).
export default function RoundCard({
  game,
  isOpen,
  onClose,
  onSubmitted,
}: {
  game: GameData;
  isOpen: boolean;
  onClose: () => void;
  onSubmitted: (data: GameData) => void;
}) {
  const navigate = useNavigate();
  const isDesktop = useMediaQuery('(min-width: 640px)');

  const [outcome, setOutcome] = useState<CardOutcome | null>(null);
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
    onSubmitted(data);
  };

  const seat = (player: PlayerScore) => (
    <Seat
      player={player}
      showWin={isWin}
      showLose={showLose}
      isWinner={winner === player.id}
      isLoser={loser === player.id}
      loseLabel={isPenalty ? 'fined' : 'dealt in'}
      onWin={() => pickWinner(player.id)}
      onLose={() => pickLoser(player.id)}
    />
  );

  return (
    <Modal.Backdrop
      isOpen={isOpen}
      onOpenChange={(open) => {
        if (!open) onClose();
      }}
      // No enter/exit animation: the card should just be there.
      className="animate-none!"
    >
      <Modal.Container
        placement="center"
        // Full screen on phones, a centered card on desktop.
        size={isDesktop ? 'md' : 'full'}
        scroll="inside"
        className="animate-none!"
      >
        <Modal.Dialog aria-label={`Round ${game.round.kanji}`}>
          <Modal.CloseTrigger />
          <Modal.Header>
            <Modal.Heading className="text-2xl">
              {game.round.kanji}
            </Modal.Heading>
          </Modal.Header>
          <Modal.Body className="flex flex-col gap-4">
            <SegmentedPicker
              ariaLabel="Outcome"
              options={OUTCOME_OPTIONS}
              selectedKey={outcome}
              onSelectionChange={setOutcome}
            />

            {isPenalty && (
              <SegmentedPicker
                ariaLabel="Penalty type"
                options={[
                  { value: Outcome.Penalty, label: 'Penalty' },
                  { value: Outcome.FalseMahjong, label: 'False mahjong' },
                ]}
                selectedKey={penaltyType}
                onSelectionChange={(value) => setPenaltyType(value)}
              />
            )}

            {/* The table as seen from the bottom seat (this device). */}
            <div className="grid grid-cols-3 grid-rows-3 items-center gap-x-2 gap-y-3">
              <div />
              {seat(game.opp)}
              <div />
              {seat(game.left)}
              <div />
              {seat(game.right)}
              <div />
              {seat(game.me)}
              <div />
            </div>

            {isWin && (
              <ValueStepper
                label="Hand value"
                value={hand}
                min={MIN_HAND}
                onChange={setHand}
              />
            )}

            {isPenalty && !isFalseMahjong && (
              <ValueStepper
                label="Fine"
                value={fine}
                min={1}
                onChange={setFine}
              />
            )}

            {isFalseMahjong && (
              <SegmentedPicker
                label="Fine"
                options={[
                  { value: '30', label: '30' },
                  { value: '60', label: '60' },
                ]}
                selectedKey={falseMahjongFine ? String(falseMahjongFine) : null}
                onSelectionChange={(value) =>
                  setFalseMahjongFine(Number(value))
                }
              />
            )}

            {isPenalty && (
              <Switch isSelected={forceNextRound} onChange={setForceNextRound}>
                <Switch.Content>
                  <Switch.Control>
                    <Switch.Thumb />
                  </Switch.Control>
                  Force next round
                </Switch.Content>
              </Switch>
            )}
          </Modal.Body>
          <Modal.Footer className="justify-center gap-3">
            <Button variant="secondary" onPress={onClose}>
              Cancel
            </Button>
            <PendingButton
              variant="primary"
              isDisabled={!isComplete}
              isPending={isPending}
              onPress={submit}
            >
              OK
            </PendingButton>
          </Modal.Footer>
        </Modal.Dialog>
      </Modal.Container>
    </Modal.Backdrop>
  );
}
