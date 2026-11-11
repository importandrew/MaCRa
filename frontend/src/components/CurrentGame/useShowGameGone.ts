import { useTranslation } from 'react-i18next';
import { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { GameStatusResponse } from '../../providers/game/gameDataApi';
import { eventGamePath, fetchWithStatus, GetRequest } from '../Request';
import { showToast } from '../Toast';

// useShowGameGone returns showGameGone, for a table whose game is no longer
// the player's live game (/now answers 404): someone else entered the last
// round, finished or aborted it. An ended game opens its result, an aborted
// one the home page.
export function useShowGameGone() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return useCallback(
    async (eventId: number, gameId: number) => {
      const { data } = await fetchWithStatus<GameStatusResponse>(
        GetRequest(eventGamePath(eventId, gameId, '/status')),
        { isToastEnabled: false }
      );

      if (data?.status === 'ended') {
        showToast({ title: t('game.ended'), variant: 'warning' });
        navigate('/last');
        return;
      }
      if (data?.status === 'aborted') {
        showToast({ title: t('game.wasAborted'), variant: 'warning' });
        navigate('/');
        return;
      }
      showToast({ title: t('game.notFound'), variant: 'danger' });
    },
    [navigate, t]
  );
}
