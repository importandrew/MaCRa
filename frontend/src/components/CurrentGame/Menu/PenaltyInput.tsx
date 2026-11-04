import { Button, Form, Label, Input, Switch, TextField } from '@heroui/react';
import { FormEvent, useEffect, useState } from 'react';
import { showToast } from '../../Toast';
import { PlayerOrder } from '../../../models/GameData';
import { authFetch, PostRequest, eventPath } from '../../Request';
import { getStoredEventId } from '../../../providers/event/eventConstants';
import { GameResult } from '../../../models/GameResult';
import PlayerPicker from '../../PlayerPicker';
import SegmentedPicker from '../../SegmentedPicker';

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

    const eventId = getStoredEventId();
    if (!eventId) return;

    const roundResultRequest = PostRequest(
      eventPath(eventId, distributed ? '/false-mahjong' : '/penalty'),
      {
        loser,
        value,
        force_next_round: forceNextRound,
      }
    );
    const data = await authFetch<GameResult>(roundResultRequest);

    if (!data) return;

    showToast({
      title: 'Penalty successfully added',
      variant: 'warning',
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
        <PlayerPicker
          label="Select player to fine"
          players={players ?? []}
          selectedId={loser ? String(loser) : ''}
          onChange={(id) => setLoser(Number(id))}
        />

        <Switch isSelected={distributed} onChange={handleIsDistributed}>
          <Switch.Content>
            <Switch.Control>
              <Switch.Thumb />
            </Switch.Control>
            False mahjong
          </Switch.Content>
        </Switch>

        <Switch isSelected={forceNextRound} onChange={setForceNextRound}>
          <Switch.Content>
            <Switch.Control>
              <Switch.Thumb />
            </Switch.Control>
            Force next round
          </Switch.Content>
        </Switch>

        {distributed ? (
          <SegmentedPicker
            label="Fine value"
            options={[
              { value: '30', label: '30' },
              { value: '60', label: '60' },
            ]}
            selectedKey={value === 30 || value === 60 ? String(value) : null}
            onSelectionChange={(next) => setValue(Number(next))}
          />
        ) : (
          // A text field with a numeric keyboard: NumberField commits its
          // value only on blur, which would keep Submit disabled on the tap.
          <TextField
            fullWidth
            value={value === 0 ? '' : value.toString()}
            onChange={(text) => {
              const digits = text.replace(/\D/g, '');
              setValue(digits ? Number(digits) : 0);
            }}
          >
            <Label>Fine value</Label>
            <Input inputMode="numeric" pattern="\d*" autoComplete="off" />
          </TextField>
        )}

        <div className="flex justify-center sm:max-w-[300px]">
          <Button type="submit" variant="primary" isDisabled={!loser || !value}>
            Submit
          </Button>
        </div>
      </div>
    </Form>
  );
}
