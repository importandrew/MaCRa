import { Button, Form, Input, Select, SelectItem } from '@heroui/react';
import { FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { GameData, Outcome, PlayerOrder } from '../../models/GameData';
import { PostRequest } from '../Request';
import { OUTCOME_LABELS, getWindLabel } from '../../utils/constants';

export default function ScroreInput({
  data,
  playersList,
  refetch,
}: {
  data: GameData;
  playersList: PlayerOrder[] | null;
  refetch: () => void;
}) {
  const navigate = useNavigate();
  const [isHidden, setIsHidden] = useState(true);
  const [outcome, setOutcome] = useState('');
  const [winner, setWinner] = useState<number | null>(null);
  const [loser, setLoser] = useState<number | null>(null);
  const [score, setScore] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const players: PlayerOrder[] = playersList || [];

  const toggleForm = () => {
    setIsHidden(!isHidden);
  };

  const resetRoundResults = () => {
    refetch();
    setOutcome('');
    setWinner(null);
    setLoser(null);
    setScore(0);
    toggleForm();
  };

  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const roundResult = PostRequest('/api/round', {
        outcome,
        won: winner,
        dealt: loser,
        score,
      });
      const response = await fetch(roundResult);

      if (!response.ok) {
        const result = await response.json();
        toast.error(`Failed to submit round results: ${result.error}`);
        return;
      }

      if (data?.round === 16) {
        navigate('/last');
        return;
      }

      resetRoundResults();
    } catch (error) {
      toast.error('Failed to submit round results. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="bg-white relative z-10 py-3 justify-self-center self-center items-center w-full text-5xl">
      <Button
        size="lg"
        variant="flat"
        color="primary"
        className="p-4 h-16 text-5xl"
        type="button"
        onPress={() => toggleForm()}
      >
        {getWindLabel(data.round)}
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
              onChange={(event) => {
                setOutcome(event.target.value);
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
          <div
            id="winner"
            className={
              !(outcome == Outcome.Self || outcome == Outcome.Deal)
                ? 'hidden'
                : 'block'
            }
          >
            <Select
              label="Player won"
              className="w-full md:w-[300px] inline-flex"
              onChange={(event) => setWinner(Number(event.target.value))}
              validate={(value) => {
                if (outcome === Outcome.Draw) return null;

                return !value?.length ? 'Winner is required' : null;
              }}
            >
              {players.map((player) => (
                <SelectItem key={player.id}>{player.name}</SelectItem>
              ))}
            </Select>
          </div>
          <div
            id="loser"
            className={outcome != Outcome.Deal || !winner ? 'hidden' : 'block'}
          >
            <Select
              label="Player dealt"
              className="w-full md:w-[300px] inline-flex"
              onChange={(event) => setLoser(Number(event.target.value))}
              validate={(value) => {
                if (outcome === Outcome.Draw || outcome === Outcome.Self)
                  return null;

                return !value?.length ? 'Loser is required' : null;
              }}
            >
              {players.map((player) => {
                if (player.id === winner) {
                  return null;
                }

                return <SelectItem key={player.id}>{player.name}</SelectItem>;
              })}
            </Select>
          </div>
          <div
            className={
              !(outcome == Outcome.Self || outcome == Outcome.Deal)
                ? 'hidden'
                : 'block'
            }
          >
            <Input
              type="number"
              isClearable={true}
              min="0"
              pattern="\d*"
              label="Hand value"
              className="w-full md:w-[300px] inline-flex"
              onChange={(event) => {
                setScore(+event.target.value);
              }}
              validate={(value) => {
                if (outcome === Outcome.Draw) return null;

                return !value.length ? 'Hand value is required' : null;
              }}
            />
          </div>
          <div className="flex justify-center">
            <Button
              type="submit"
              variant="flat"
              color="primary"
              isDisabled={!outcome}
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
