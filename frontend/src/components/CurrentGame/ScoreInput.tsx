import {
  Button,
  FieldError,
  Form,
  Label,
  ListBox,
  Select,
  Input,
  TextField,
} from '@heroui/react';
import { FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  GameData,
  Outcome,
  PlayerOrder,
  RoundData,
} from '../../models/GameData';
import { OUTCOME_LABELS } from '../../utils/constants';
import { authFetch, PostRequest, eventPath } from '../Request';
import { getStoredEventId } from '../../providers/event/eventConstants';
import PendingButton from '../PendingButton';
import SegmentedPicker from '../SegmentedPicker';

// PlayerSelect is a compact player dropdown. The round form sits in the
// middle cell of the table, which is too narrow for one button per player.
function PlayerSelect({
  label,
  players,
  selectedId,
  onChange,
  disabledId,
}: {
  label: string;
  players: PlayerOrder[];
  selectedId: string;
  onChange: (id: string) => void;
  disabledId: string;
}) {
  return (
    <Select
      fullWidth
      placeholder="Select player"
      value={selectedId || null}
      disabledKeys={disabledId ? [disabledId] : []}
      onChange={(key) => onChange(key == null ? '' : String(key))}
    >
      <Label>{label}</Label>
      <Select.Trigger>
        <Select.Value />
        <Select.Indicator />
      </Select.Trigger>
      <Select.Popover>
        <ListBox>
          {players.map((player) => (
            <ListBox.Item
              key={player.id}
              id={String(player.id)}
              textValue={player.name}
            >
              {player.name}
              <ListBox.ItemIndicator />
            </ListBox.Item>
          ))}
        </ListBox>
      </Select.Popover>
    </Select>
  );
}

export default function ScroreInput({
  kanji,
  playersList,
  setRoundData,
}: {
  kanji: RoundData['kanji'];
  playersList: PlayerOrder[] | null;
  setRoundData: (_data: GameData) => void;
}) {
  const navigate = useNavigate();
  const [isHidden, setIsHidden] = useState(true);
  const [outcome, setOutcome] = useState('');
  const [winner, setWinner] = useState('');
  const [loser, setLoser] = useState('');
  const [score, setScore] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const players: PlayerOrder[] = playersList || [];

  const toggleForm = () => {
    setIsHidden(!isHidden);
  };

  const resetRoundResults = () => {
    setOutcome('');
    setWinner('');
    setLoser('');
    setScore(0);
    setError('');
  };

  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);

    const eventId = getStoredEventId();
    if (!eventId) {
      setIsLoading(false);
      return;
    }

    // Winner, loser and hand value are kept when the outcome changes, so
    // send only the ones that apply to the chosen outcome.
    const isWin = outcome === Outcome.Self || outcome === Outcome.Deal;
    const roundResult = PostRequest(eventPath(eventId, '/round'), {
      outcome,
      won: isWin ? Number(winner) || null : null,
      dealt: outcome === Outcome.Deal ? Number(loser) || null : null,
      score: isWin ? score : 0,
    });
    const data = await authFetch<GameData>(roundResult, {
      isToastEnabled: true,
    });

    if (!data) {
      setIsLoading(false);
      return;
    }

    if (data?.status === 'ended' || data?.status === 'aborted') {
      navigate('/last');
      return;
    }

    setRoundData(data);
    resetRoundResults();
    toggleForm();
    setIsLoading(false);
  };

  const isWinOutcome = outcome === Outcome.Self || outcome === Outcome.Deal;
  const isHandTooLow = score > 0 && score < 8;
  const isComplete =
    outcome === Outcome.Draw ||
    (outcome === Outcome.Self && winner !== '' && score >= 8) ||
    (outcome === Outcome.Deal && winner !== '' && loser !== '' && score >= 8);

  return (
    <div className="relative z-10 py-3 justify-self-center self-center items-center w-full text-5xl">
      <Button
        size="lg"
        variant="secondary"
        className="p-4 h-16 text-5xl"
        type="button"
        onPress={() => toggleForm()}
      >
        {kanji}
      </Button>
      {/* Toggle button groups instead of selects: one tap per choice and no
          popover on phones. */}
      <Form
        validationBehavior="native"
        onSubmit={onSubmit}
        className={'mt-3 text-base ' + (isHidden ? 'hidden' : '')}
      >
        <div className="flex flex-col gap-3 w-full md:w-[300px] mx-auto text-start">
          <SegmentedPicker
            label="Outcome"
            options={[Outcome.Draw, Outcome.Self, Outcome.Deal].map(
              (value, index) => ({
                value,
                // Non-breaking hyphen: "Deal-in" must not wrap in the narrow form.
                label: OUTCOME_LABELS[index].replace('-', '\u2011'),
              })
            )}
            selectedKey={outcome ? (outcome as Outcome) : null}
            onSelectionChange={(value) => {
              // Keep the picked players and hand value: switching Self ⇄
              // Deal-in shouldn't make the user pick the winner again.
              setOutcome(value);
              setError('');
            }}
          />

          {isWinOutcome && (
            <PlayerSelect
              label="Player won"
              players={players}
              selectedId={winner}
              onChange={(id) => {
                setWinner(id);
                // One player can't both win and deal in.
                if (id === loser) setLoser('');
              }}
              disabledId={outcome === Outcome.Deal ? loser : ''}
            />
          )}

          {outcome === Outcome.Deal && (
            <PlayerSelect
              label="Player dealt"
              players={players}
              selectedId={loser}
              onChange={setLoser}
              disabledId={winner}
            />
          )}

          {isWinOutcome && (
            // A text field with a numeric keyboard, not NumberField:
            // NumberField commits its value only on blur, so OK would still
            // be disabled when the user taps it right after typing.
            <TextField
              fullWidth
              value={score === 0 ? '' : score.toString()}
              isInvalid={error !== ''}
              onChange={(text) => {
                const digits = text.replace(/\D/g, '');
                setScore(digits ? Number(digits) : 0);
                setError('');
              }}
              onBlur={() => {
                setError(
                  isHandTooLow ? 'Hand score should be 8 or higher' : ''
                );
              }}
            >
              <Label>Hand value</Label>
              <Input inputMode="numeric" pattern="\d*" autoComplete="off" />
              <FieldError>{error}</FieldError>
            </TextField>
          )}

          <div className="flex justify-center">
            <PendingButton
              type="submit"
              variant="primary"
              isDisabled={!isComplete}
              isPending={isLoading}
            >
              OK
            </PendingButton>
          </div>
        </div>
      </Form>
    </div>
  );
}
