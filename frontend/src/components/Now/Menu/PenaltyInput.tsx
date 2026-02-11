import {
  Button,
  Checkbox,
  Form,
  Input,
  Radio,
  RadioGroup,
} from '@heroui/react';
import { FormEvent, useEffect, useState } from 'react';
import { showToast } from '../../Toast';
import { PlayerOrder } from '../../../models/GameData';
import { globalFetch, PostRequest } from '../../Request';
import { GameResult } from '../../../models/GameResult';

export default function PenaltyInput({
  players = [],
  finedPlayerId = null,
  setFinedPlayerId,
}: {
  players: PlayerOrder[] | null;
  finedPlayerId: number | null;
  setFinedPlayerId: (value: number | null) => void;
}) {
  const [loser, setLoser] = useState<number | null>(finedPlayerId);
  const [distributed, setDistibuted] = useState(false);
  const [forceNextRound, setForceNextRound] = useState(false);
  const [value, setValue] = useState<number>(0);

  const resetPenaltyResults = () => {
    setLoser(null);
    setFinedPlayerId(0);
    window.requestAnimationFrame(() => setFinedPlayerId(null));
    setDistibuted(false);
    setForceNextRound(false);
    setValue(0);
  };

  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const roundResultRequest = PostRequest(
      `/api/${distributed ? 'false-mahjong' : 'penalty'}`,
      {
        loser,
        value,
        force_next_round: forceNextRound,
      }
    );
    const data = await globalFetch<GameResult>(roundResultRequest);

    if (!data) return;

    showToast({
      title: 'Penalty successfully added',
      color: 'warning',
    });
    resetPenaltyResults();
  };

  const handleIsDistributed = (value: boolean) => {
    setValue(0);
    setDistibuted(value);
  };

  useEffect(() => {
    if (distributed === false) {
      setValue(0);
    }
  }, [distributed]);

  return (
    <Form
      validationBehavior="native"
      onSubmit={onSubmit}
      className="mb-3 items-center"
    >
      <div className="flex flex-col gap-4 w-full sm:w-[300px]">
        <div>
          <RadioGroup
            label="Select player to fine:"
            value={loser?.toString()}
            onValueChange={(value: string) => setLoser(Number(value))}
            color="warning"
          >
            {players &&
              players.map((player) => (
                <Radio
                  key={player.id}
                  value={player.id.toString()}
                  description={`[${player.id}]`}
                >
                  {player.name}
                </Radio>
              ))}
          </RadioGroup>
        </div>

        <div>
          <div className="text-3x1 text-gray-500 mb-2">Is false mahjong?</div>
          <Checkbox
            radius="full"
            color="warning"
            isSelected={distributed}
            onValueChange={handleIsDistributed}
          >
            Yes
          </Checkbox>
        </div>

        <div>
          <div className="text-3x1 text-gray-500 mb-2">Force next round?</div>
          <Checkbox
            radius="full"
            color="warning"
            isSelected={forceNextRound}
            onValueChange={setForceNextRound}
          >
            Yes
          </Checkbox>
        </div>

        <div className={distributed ? '' : 'hidden'}>
          <RadioGroup
            label="Fine value:"
            value={value.toString()}
            onValueChange={(value) => setValue(Number(value))}
            validate={(value: string | null) => {
              if (!distributed) return null;

              return Number(value) === 30 || Number(value) === 60
                ? null
                : 'Fine value is required';
            }}
            color="warning"
            orientation="horizontal"
          >
            <Radio className="mr-3" value="30">
              30
            </Radio>
            <Radio className="mr-3" value="60">
              60
            </Radio>
          </RadioGroup>
        </div>

        <div className={'sm:max-w-[300px] ' + (distributed ? 'hidden' : '')}>
          <Input
            type="number"
            isClearable={true}
            value={value.toString()}
            min="0"
            pattern="\d*"
            label="Fine value"
            className="w-full inline-flex"
            onChange={(event) => {
              setValue(+event.target.value);
            }}
            validate={(value) => {
              if (!distributed) return null;

              return Number(value) === 0 ? 'Fine value is required' : null;
            }}
          />
        </div>
        <div className="flex justify-center sm:max-w-[300px]">
          <Button
            type="submit"
            variant="flat"
            color="warning"
            isDisabled={!loser || !value}
          >
            Submit
          </Button>
        </div>
      </div>
    </Form>
  );
}
