import { Button, Form, Input, Select, SelectItem } from '@heroui/react';
import { FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  GameData,
  Outcome,
  PlayerOrder,
  RoundData,
} from '../../models/GameData';
import { OUTCOME_LABELS } from '../../utils/constants';
import { globalFetch, PostRequest } from '../Request';

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

    const roundResult = PostRequest('/api/round', {
      outcome,
      won: Number(winner) || null,
      dealt: Number(loser) || null,
      score,
    });
    const data = await globalFetch<GameData>(roundResult, {
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

  return (
    <div className="relative z-10 py-3 justify-self-center self-center items-center w-full text-5xl">
      <Button
        size="lg"
        variant="flat"
        color="primary"
        className="p-4 h-16 text-5xl"
        type="button"
        onPress={() => toggleForm()}
      >
        {kanji}
      </Button>
      <Form
        validationBehavior="native"
        onSubmit={onSubmit}
        className={'mt-3 ' + (isHidden ? 'hidden' : '')}
      >
        <div className="flex flex-col gap-3 w-full">
          <div className="block">
            <Select
              label="Outcome"
              className="w-full md:w-[300px] inline-flex"
              selectedKeys={outcome ? [outcome] : []}
              onChange={(event) => {
                resetRoundResults();
                setOutcome(event.target.value as Outcome);
              }}
              validate={(value) => {
                return !value?.length ? 'Outcome is required' : null;
              }}
            >
              {Object.values(Outcome)
                .map((key, index) => (
                  <SelectItem key={key}>{OUTCOME_LABELS[index]}</SelectItem>
                ))
                .splice(0, 3)}
            </Select>
          </div>
          <>
            {outcome == Outcome.Self || outcome == Outcome.Deal ? (
              <div>
                <Select
                  label="Player won"
                  className="w-full md:w-[300px] inline-flex"
                  selectedKeys={[winner]}
                  onChange={(event) => setWinner(event.target.value)}
                  validate={(value) => {
                    return !value?.length ? 'Winner is required' : null;
                  }}
                >
                  {players.map((player) => {
                    if (player.id.toString() === loser) {
                      return null;
                    }

                    return (
                      <SelectItem key={String(player.id)}>
                        {player.name}
                      </SelectItem>
                    );
                  })}
                </Select>
              </div>
            ) : null}
          </>

          <>
            {outcome === Outcome.Deal ? (
              <div>
                <Select
                  label="Player dealt"
                  className="w-full md:w-[300px] inline-flex"
                  selectedKeys={[loser]}
                  onChange={(event) => setLoser(event.target.value)}
                  validate={(value) => {
                    return !value?.length ? 'Loser is required' : null;
                  }}
                >
                  {players.map((player) => {
                    if (player.id.toString() === winner) {
                      return null;
                    }

                    return (
                      <SelectItem key={player.id}>{player.name}</SelectItem>
                    );
                  })}
                </Select>
              </div>
            ) : null}
          </>

          <>
            {outcome == Outcome.Self || outcome == Outcome.Deal ? (
              <div>
                <Input
                  type="text"
                  inputMode="numeric"
                  isClearable={true}
                  isInvalid={error !== ''}
                  onClear={() => {
                    setScore(0);
                    setError('');
                  }}
                  onBlur={() => {
                    if (score < 8 && score > 0) {
                      setError('Hand score should be 8 or higher');
                    } else {
                      setError('');
                    }
                  }}
                  errorMessage={error}
                  pattern="\d*"
                  autoComplete="off"
                  label="Hand value"
                  className="w-full md:w-[300px] inline-flex"
                  onChange={(event) => {
                    setScore(+event.target.value);
                  }}
                  value={score === 0 ? '' : score.toString()}
                />
              </div>
            ) : null}
          </>
          <div className="flex justify-center">
            <Button
              type="submit"
              variant="flat"
              color="primary"
              isDisabled={!outcome || (outcome !== Outcome.Draw && score === 0)}
              isLoading={isLoading}
            >
              OK
            </Button>
          </div>
        </div>
      </Form>
    </div>
  );
}
