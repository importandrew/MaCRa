import { Button, Form, Spinner } from '@heroui/react';
import { FormEvent, useCallback, useEffect, useMemo, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { showToast } from '../Toast';
import { Player } from '../../models/Player';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';
import { useFetchEventMembers } from '../../providers/event/eventApi';
import { getStoredEventId } from '../../providers/event/eventConstants';
import { useEventCapabilities } from '../../providers/event/useEventCapabilities';
import { WIND_LABELS, WIND_SYMBOLS } from '../../utils/constants';
import { shuffleArray } from '../../utils/shuffleArray';
import { authFetch, PostRequest, eventPath } from '../Request';
import PlayerSelector from './PlayerSelector';

export default function Draft() {
  const navigate = useNavigate();
  const location = useLocation();
  const [players, setPlayers] = useState<(number | undefined)[]>([]);
  const eventId = getStoredEventId();
  const { data: members, isLoading: membersLoading } =
    useFetchEventMembers(eventId);
  const {
    hasCapability,
    isLoading: capsLoading,
    me,
  } = useEventCapabilities(eventId);

  const suggestions: Player[] = useMemo(
    () =>
      (members ?? [])
        .filter((member) => member.roles.includes('player'))
        .map((member) => ({
          id: member.user_id,
          name: member.name,
        })),
    [members]
  );

  useEffect(() => {
    if (capsLoading || me === null) return;
    if (!hasCapability('create_game')) {
      showToast({
        title: 'You cannot start a game in this event',
        color: 'warning',
      });
      navigate('/');
    }
  }, [capsLoading, me, hasCapability, navigate]);

  // Derived state for selectedKeys
  const selectedKeys = useMemo(
    () => players.map((playerId) => playerId?.toString() || ''),
    [players]
  );

  // Set initial player IDs if provided in location state
  useEffect(() => {
    if (location.state?.playerIds) {
      const playerIds = location.state.playerIds as number[];
      // Pad with undefined values if less than 4 players
      const paddedPlayerIds: (number | undefined)[] = [...playerIds];
      while (paddedPlayerIds.length < 4) {
        paddedPlayerIds.push(undefined);
      }
      setPlayers(paddedPlayerIds.slice(0, 4));
    }
  }, [location.state]);

  // Handle player selection changes
  const handleSuggestionsChange = useCallback(() => {
    setPlayers([...players]); // Ensure immutability
  }, [players]);

  // Clear players and selectedKeys
  const onClear = () => {
    setPlayers([undefined, undefined, undefined, undefined]);
  };

  // Shuffle players
  const onShuffle = useCallback(() => {
    const shuffledPlayers = shuffleArray([...players]); // Shuffle a copy of the array
    setPlayers(shuffledPlayers);
  }, [players]);

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const currentPlayerId = JSON.parse(
      localStorage.getItem(AUTH_USER_DATA) || '{}'
    ).id;

    if (!players.includes(currentPlayerId)) {
      showToast({
        title: 'You should also pick youself to start the game',
        color: 'warning',
      });
      return;
    }

    const eventId = getStoredEventId();
    if (!eventId) {
      showToast({
        title: 'No event selected',
        description: 'Please select an event first',
        color: 'warning',
      });
      return;
    }

    const draftRequest = PostRequest(eventPath(eventId, '/new'), { players });
    const data = await authFetch<{ id: number }>(draftRequest);

    if (!data || 'error' in data) return;

    navigate('/now');
  };

  return (
    <div className="flex flex-col justify-center grow w-screen h-full my-auto">
      {membersLoading || capsLoading ? (
        <Spinner label="Loading..." size="lg" />
      ) : (
        <Form
          validationBehavior="native"
          onSubmit={onSubmit}
          className="items-center"
        >
          <div className="flex flex-col gap-3 p-3">
            {WIND_LABELS.map((wind: string, index: number) => (
              <div
                key={wind}
                className="relative flex justify-center items-center"
              >
                <span className="absolute left-[-24px] text-base pr-3">
                  {WIND_SYMBOLS[index]}
                </span>
                <PlayerSelector
                  wind={wind}
                  index={index}
                  players={players}
                  suggestions={suggestions}
                  selectedKeys={selectedKeys}
                  handleSuggestionsChange={handleSuggestionsChange}
                />
              </div>
            ))}

            <div className="flex justify-center gap-3">
              <Button
                variant="flat"
                type="button"
                color="default"
                onPress={onClear}
              >
                Clear
              </Button>
              <Button
                variant="flat"
                type="button"
                color="secondary"
                onPress={onShuffle}
              >
                Shuffle
              </Button>
              <Button
                variant="flat"
                type="submit"
                color="primary"
                isDisabled={
                  players.length !== 4 || selectedKeys.some((key) => !key)
                }
              >
                Play
              </Button>
            </div>
          </div>
        </Form>
      )}
    </div>
  );
}
