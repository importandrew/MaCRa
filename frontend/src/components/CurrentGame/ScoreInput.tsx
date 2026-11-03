import {
  Button,
  FieldError,
  Form,
  Label,
  ListBox,
  Select,
  Input,
  Spinner,
  ToggleButton,
  TextField,
  ToggleButtonGroup,
  type Key,
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

// PlayerSelect is a compact player dropdown. The round form sits in the
// middle cell of the table, which is too narrow for one button per player.
// TODO.local.md item 10 replaces this form with seat-based input.
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

    const roundResult = PostRequest(eventPath(eventId, '/round'), {
      outcome,
      won: Number(winner) || null,
      dealt: Number(loser) || null,
      score,
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

  // firstKey returns the selected key of a single-selection group, or ''.
  const firstKey = (keys: Set<Key>) => String([...keys][0] ?? '');

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
          <div className="flex flex-col gap-1">
            <Label id="outcome-label">Outcome</Label>
            <ToggleButtonGroup
              aria-labelledby="outcome-label"
              selectionMode="single"
              disallowEmptySelection
              fullWidth
              selectedKeys={outcome ? [outcome] : []}
              onSelectionChange={(keys) => {
                resetRoundResults();
                setOutcome(firstKey(keys) as Outcome);
              }}
            >
              {[Outcome.Draw, Outcome.Self, Outcome.Deal].map((key, index) => (
                <ToggleButton key={key} id={key}>
                  {index > 0 && <ToggleButtonGroup.Separator />}
                  {OUTCOME_LABELS[index]}
                </ToggleButton>
              ))}
            </ToggleButtonGroup>
          </div>

          {isWinOutcome && (
            <PlayerSelect
              label="Player won"
              players={players}
              selectedId={winner}
              onChange={setWinner}
              disabledId={loser}
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
            <Button
              type="submit"
              variant="primary"
              isDisabled={!isComplete}
              isPending={isLoading}
            >
              {({ isPending }) => (
                <>
                  {isPending && <Spinner color="current" size="sm" />}
                  OK
                </>
              )}
            </Button>
          </div>
        </div>
      </Form>
    </div>
  );
}
