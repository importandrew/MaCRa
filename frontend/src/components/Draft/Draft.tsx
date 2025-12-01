import { Button, Form } from '@heroui/react';
import { FormEvent, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Player } from '../../models/Player';
import { WIND_LABELS, WIND_SYMBOLS } from '../../utils/constants';
import { GetRequest, PostRequest } from '../Request';
import PlayerSelector from './PlayerSelector';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';

interface SelectorState {
  inputValue: string;
  suggestions: Player[];
  selectedPlayerId: number | undefined;
  isFocused: boolean;
  needsInitialFetch: boolean;
}

const createInitialSelectorState = (): SelectorState => ({
  inputValue: '',
  suggestions: [],
  selectedPlayerId: undefined,
  isFocused: false,
  needsInitialFetch: true,
});

export default function Draft() {
  const navigate = useNavigate();
  const [selectors, setSelectors] = useState<SelectorState[]>([
    createInitialSelectorState(),
    createInitialSelectorState(),
    createInitialSelectorState(),
    createInitialSelectorState(),
  ]);
  const debounceTimers = useRef<{ [key: number]: ReturnType<typeof setTimeout> }>({});

  // Derived state for selectedKeys (all selected player IDs across all selectors)
  const selectedKeys = useMemo(
    () => selectors.map((s) => s.selectedPlayerId?.toString() || ''),
    [selectors]
  );

  // Fetch players from API
  const fetchPlayers = useCallback(async (name: string, limit: number = 10) => {
    try {
      const playersRequest = GetRequest('/api/players/list', { name, limit });
      const response = await fetch(playersRequest);
      const result = await response.json();

      if (!response.ok) {
        toast.error(result.error || 'Failed to fetch players');
        return [];
      }

      return result || [];
    } catch (error) {
      toast.error('Failed to fetch players. Please try again.');
      return [];
    }
  }, []);

  // Initial fetch on mount - single API call with empty name
  useEffect(() => {
    const initialFetch = async () => {
      const players = await fetchPlayers('', 10);

      setSelectors((prev) =>
        prev.map((selector) => ({
          ...selector,
          suggestions: players,
          needsInitialFetch: false,
        }))
      );
    };

    initialFetch();
  }, [fetchPlayers]);

  // Cleanup debounce timers on unmount
  useEffect(() => {
    return () => {
      Object.values(debounceTimers.current).forEach((timer) =>
        clearTimeout(timer)
      );
    };
  }, []);

  // Handle input change with debounce
  const handleInputChange = useCallback(
    (index: number, value: string) => {
      // Clear existing timer for this selector
      if (debounceTimers.current[index]) {
        clearTimeout(debounceTimers.current[index]);
      }

      // Update input value immediately and clear selection if user is typing something different
      setSelectors((prev) => {
        const updated = [...prev];
        const currentSelector = updated[index];
        const selectedPlayer = currentSelector.suggestions.find(
          (p) => p.id === currentSelector.selectedPlayerId
        );
        
        // If user is typing something different from selected player's name, clear selection
        if (selectedPlayer && value !== selectedPlayer.name) {
          updated[index] = {
            ...currentSelector,
            inputValue: value,
            selectedPlayerId: undefined,
          };
        } else {
          updated[index] = { ...currentSelector, inputValue: value };
        }
        return updated;
      });

      // Debounce API call
      debounceTimers.current[index] = setTimeout(async () => {
        const players = await fetchPlayers(value, 10);

        setSelectors((prev) => {
          const updated = [...prev];
          updated[index] = { ...updated[index], suggestions: players };
          return updated;
        });
      }, 300);
    },
    [fetchPlayers]
  );

  // Handle player selection
  const handleSelectionChange = useCallback((index: number, playerId: number | null) => {
    setSelectors((prev) => {
      const updated = [...prev];
      const selectedPlayer = updated[index].suggestions.find((p) => p.id === playerId);
      updated[index] = {
        ...updated[index],
        selectedPlayerId: playerId || undefined,
        inputValue: selectedPlayer?.name || updated[index].inputValue,
      };
      return updated;
    });
  }, []);

  // Handle clear for a single selector
  const handleSelectorClear = useCallback(
    (index: number) => {
      const selector = selectors[index];
      const isFocused = selector.isFocused;

      setSelectors((prev) => {
        const updated = [...prev];
        updated[index] = {
          ...createInitialSelectorState(),
          needsInitialFetch: !isFocused,
        };
        return updated;
      });

      if (isFocused) {
        // If focused, fetch immediately
        fetchPlayers('', 10).then((players) => {
          setSelectors((prev) => {
            const updated = [...prev];
            updated[index] = { ...updated[index], suggestions: players };
            return updated;
          });
        });
      }
    },
    [selectors, fetchPlayers]
  );

  // Handle focus
  const handleFocus = useCallback(
    (index: number) => {
      setSelectors((prev) => {
        const updated = [...prev];
        const currentSelector = updated[index];
        updated[index] = { ...currentSelector, isFocused: true };

        // If needs initial fetch or has input value, fetch with current input value
        const fetchName = currentSelector.inputValue || '';
        if (currentSelector.needsInitialFetch || currentSelector.inputValue) {
          updated[index].needsInitialFetch = false;
          fetchPlayers(fetchName, 10).then((players) => {
            setSelectors((current) => {
              const currentUpdated = [...current];
              currentUpdated[index] = {
                ...currentUpdated[index],
                suggestions: players,
              };
              return currentUpdated;
            });
          });
        }

        return updated;
      });
    },
    [fetchPlayers]
  );

  // Handle blur
  const handleBlur = useCallback((index: number) => {
    setSelectors((prev) => {
      const updated = [...prev];
      updated[index] = { ...updated[index], isFocused: false };
      return updated;
    });
  }, []);

  // Clear all selectors
  const onClear = useCallback(() => {
    // Clear all debounce timers
    Object.values(debounceTimers.current).forEach((timer) => clearTimeout(timer));
    debounceTimers.current = {};

    // Reset all selectors and trigger single initial fetch
    fetchPlayers('', 10).then((players) => {
      setSelectors((prev) =>
        prev.map(() => ({
          ...createInitialSelectorState(),
          suggestions: players,
        }))
      );
    });
  }, [fetchPlayers]);

  // Shuffle selector states (winds stay fixed at their positions)
  const onShuffle = useCallback(() => {
    setSelectors((prev) => {
      const shuffled = [...prev];
      // Fisher-Yates shuffle - only shuffle the selector states, not the wind positions
      for (let i = shuffled.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
      }
      // Ensure all selectors are blurred to prevent dropdowns from opening
      return shuffled.map((selector) => ({
        ...selector,
        isFocused: false,
      }));
    });
  }, []);

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const playerIds = selectors
      .map((s) => s.selectedPlayerId)
      .filter((id): id is number => id !== undefined);

    if (playerIds.length !== 4) {
      toast.warning('Please select all four players');
      return;
    }

    const currentPlayerId = JSON.parse(
      localStorage.getItem(AUTH_USER_DATA) || '{}'
    ).id;

    if (!playerIds.includes(currentPlayerId)) {
      toast.warning('You should also pick yourself to start the game');
      return;
    }

    try {
      const draftRequest = PostRequest('/api/new', { players: playerIds });
      const response = await fetch(draftRequest);

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
      <Form
        validationBehavior="native"
        onSubmit={onSubmit}
        className="items-center"
      >
        <div className="flex flex-col gap-3 p-3">
          {WIND_LABELS.map((wind: string, index: number) => {
            const selector = selectors[index];
            const windSymbol = WIND_SYMBOLS[index];
            return (
              <div
                key={wind}
                className="relative flex justify-center items-center"
              >
                <span className="absolute left-[-24px] text-base pr-3">
                  {windSymbol}
                </span>
                  <PlayerSelector
                    wind={wind}
                    index={index}
                    inputValue={selector.inputValue}
                    suggestions={selector.suggestions}
                    selectedKey={selector.selectedPlayerId?.toString() || ''}
                    selectedKeys={selectedKeys}
                    isFocused={selector.isFocused}
                    onInputChange={(value) => handleInputChange(index, value)}
                    onSelectionChange={(playerId) =>
                      handleSelectionChange(index, playerId)
                    }
                    onClear={() => handleSelectorClear(index)}
                    onFocus={() => handleFocus(index)}
                    onBlur={() => handleBlur(index)}
                  />
              </div>
            );
          })}

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
                  selectors.length !== 4 ||
                  selectors.some((s) => s.selectedPlayerId === undefined)
                }
              >
                Play
              </Button>
            </div>
          </div>
        </Form>
    </div>
  );
}