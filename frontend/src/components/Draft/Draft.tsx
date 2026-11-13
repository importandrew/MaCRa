import { useTranslation } from 'react-i18next';
import { Button, Form } from '@heroui/react';
import { FormEvent, useCallback, useEffect, useMemo, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { useLocation, useNavigate } from 'react-router-dom';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';
import { Player } from '../../models/Player';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';
import { useFetchEventMembers } from '../../providers/event/eventApi';
import { useEvent } from '../../providers/event/eventConstants';
import { useEventCapabilities } from '../../providers/event/useEventCapabilities';
import { WINDS, WIND_SYMBOLS } from '../../utils/constants';
import { shuffleArray } from '../../utils/shuffleArray';
import { authFetch, PostRequest, eventPath } from '../Request';
import PlayerSelector from './PlayerSelector';

export default function Draft() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  // One slot per seat from the start: Shuffle only moves slots that exist,
  // so with [] and one pick it had nothing to swap with. "Play again" passes
  // the last game's players (location state), seated as they were.
  const [players, setPlayers] = useState<(number | undefined)[]>(() => {
    const given = (location.state?.playerIds as number[] | undefined) ?? [];
    return [0, 1, 2, 3].map((i) => given[i]);
  });
  const { selectedEventId: eventId } = useEvent();
  const { data: members, isLoading: membersLoading } =
    useFetchEventMembers(eventId);
  const {
    hasCapability,
    isLoading: capsLoading,
    isError: capsError,
    isForbidden: capsForbidden,
    refetch: refetchCaps,
  } = useEventCapabilities(eventId);
  const canCreateGame = hasCapability('create_game');

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

  // 403 on /me means the user has no access to the event at all (e.g. they
  // were removed from a private event), which is the same as no create_game.
  // A network error leaves the capabilities unknown: LoadingTemplate shows
  // "Try again" instead.
  useEffect(() => {
    if (capsLoading || capsError) return;
    if (capsForbidden || !canCreateGame) {
      showToast({
        title: t('draft.cannotStart'),
        variant: 'warning',
      });
      navigate('/');
    }
  }, [capsLoading, capsError, capsForbidden, canCreateGame, navigate, t]);

  // Derived state for selectedKeys
  const selectedKeys = useMemo(
    () => players.map((playerId) => playerId?.toString() || ''),
    [players]
  );

  // setSeat picks (or clears) the player of one seat.
  const setSeat = useCallback((index: number, playerId: number | undefined) => {
    setPlayers((prev) => prev.map((id, i) => (i === index ? playerId : id)));
  }, []);

  // Clear players and selectedKeys
  const onClear = () => {
    setPlayers([undefined, undefined, undefined, undefined]);
  };

  // Shuffle players
  const onShuffle = useCallback(() => {
    const shuffledPlayers = shuffleArray([...players]); // Shuffle a copy of the array
    setPlayers(shuffledPlayers);
  }, [players]);

  const { isPending: isSubmitting, run } = useSubmitOnce();

  const onSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    void run(submit);
  };

  const submit = async () => {
    const currentPlayerId = JSON.parse(
      localStorage.getItem(AUTH_USER_DATA) || '{}'
    ).id;

    if (!players.includes(currentPlayerId)) {
      showToast({
        title: t('draft.pickYourself'),
        variant: 'warning',
      });
      return;
    }

    if (!eventId) {
      showToast({
        title: t('header.noEventSelected'),
        description: t('events.selectFirst'),
        variant: 'warning',
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
      <LoadingTemplate
        isLoading={membersLoading || capsLoading}
        isError={capsError}
        onRefresh={refetchCaps}
      >
        {/* Render the form only once create_game is confirmed, so it doesn't
            flash before the redirect above. */}
        {canCreateGame && (
          <Form
            validationBehavior="native"
            onSubmit={onSubmit}
            // v3 Form is a plain <form>, not a flex column as in v2, so
            // items-center needs flex to center the selectors.
            className="flex flex-col items-center"
          >
            <div className="flex w-full flex-col gap-3 p-3 md:w-auto">
              {WINDS.map((wind, index) => (
                <PlayerSelector
                  key={wind}
                  wind={`${WIND_SYMBOLS[index]} ${t(`winds.${wind}`)}`}
                  index={index}
                  suggestions={suggestions}
                  selectedKeys={selectedKeys}
                  onChange={setSeat}
                />
              ))}

              <div className="flex justify-center gap-3">
                <Button variant="tertiary" type="button" onPress={onClear}>
                  {t('draft.clear')}
                </Button>
                <Button variant="secondary" type="button" onPress={onShuffle}>
                  {t('draft.shuffle')}
                </Button>
                <PendingButton
                  variant="primary"
                  type="submit"
                  isPending={isSubmitting}
                  isDisabled={
                    players.length !== 4 || selectedKeys.some((key) => !key)
                  }
                >
                  {t('draft.play')}
                </PendingButton>
              </div>
            </div>
          </Form>
        )}
      </LoadingTemplate>
    </div>
  );
}
