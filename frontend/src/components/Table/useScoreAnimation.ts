import { useEffect, useRef, useState } from 'react';

type SeatScore = { id: number; score: number };

// useScoreAnimation keeps the scores a table screen shows now and the ones it
// showed before they last changed, by player ID, for AnimatedScore to count
// from one to the other. players may be a new array on every render: only a
// change of the scores counts.
export function useScoreAnimation(players: SeatScore[] | null) {
  const [previousScores, setPreviousScores] = useState<Record<number, number>>(
    {}
  );
  const [currentScores, setCurrentScores] = useState<Record<number, number>>(
    {}
  );
  const currentScoresRef = useRef<Record<number, number>>({});
  const key = players?.map((p) => `${p.id}:${p.score}`).join(',') ?? null;

  useEffect(() => {
    if (key == null) return;
    const scores: Record<number, number> = {};
    for (const pair of key.split(',')) {
      const [id, score] = pair.split(':').map(Number);
      scores[id] = score;
    }
    setPreviousScores(currentScoresRef.current);
    currentScoresRef.current = scores;
    setCurrentScores(scores);
  }, [key]);

  return { previousScores, currentScores };
}
