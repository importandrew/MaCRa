import { Button, Form, Spinner } from '@heroui/react';
import { FormEvent, useCallback, useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Player } from '../../models/Player';
import { WIND_LABELS, WIND_SYMBOLS } from '../../utils/constants';
import { shuffleArray } from '../../utils/shuffleArray';
import { GetRequest, globalFetch, PostRequest } from '../Request';
import PlayerSelector from './PlayerSelector';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';

// Custom hook for fetching players
const useFetchPlayers = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [suggestions, setSuggestions] = useState<Player[]>([]);

  useEffect(() => {
    const fetchPlayers = async () => {
      setIsLoading(true);
      try {
        const playersRequest = GetRequest('/api/players/list');
        const response = await globalFetch(playersRequest);
        const result: Player[] | {  error: string  } = await response.json();

        if (!response.ok && 'error' in result) {
          toast.error(result.error);
          return;
        }

        setSuggestions((result as Player[]) || []);
      } catch (error) {
        toast.error('Failed to fetch players. Please try again.');
      } finally {
        setIsLoading(false);
      }
    };

    fetchPlayers();
  }, []);

  return { suggestions, isLoading };
};

export default function Draft() {
  const navigate = useNavigate();
  const [players, setPlayers] = useState<(number | undefined)[]>([]);
  const { suggestions, isLoading } = useFetchPlayers();

  // Derived state for selectedKeys
  const selectedKeys = useMemo(
    () => players.map((playerId) => playerId?.toString() || ''),
    [players]
  );

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
      toast.warning('You should also pick youself to start the game');
      return;
    }

    try {
      const draftRequest = PostRequest('/api/new', { players });
      const response = await globalFetch(draftRequest);

      if (!response.ok) {
        const result = await response.json();
        toast.error(`Failed to submit draft: ${result.error}`);
        return;
      }

      navigate('/now');
    } catch (error) {
      toast.error('Failed to submit draft. Please try again.');
    }
  };

  return (
    <div className="flex flex-col justify-center grow w-screen h-full my-auto">
      {isLoading ? (
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
