import { useNavigate } from 'react-router-dom';

export function useEventGamesNavigate(eventId: number) {
  const navigate = useNavigate();
  return (gameId: number) => {
    navigate(`/event/${eventId}/game/${gameId}`);
  };
}
