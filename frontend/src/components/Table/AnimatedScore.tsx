import { useEffect, useState } from 'react';

interface AnimatedScoreProps {
  previousScores: Record<number, number>;
  currentScores: Record<number, number>;
  playerId: number;
  className?: string;
}

// ANIMATION_MS is how long the score counts from the old value to the new.
const ANIMATION_MS = 2000;
// CLASS_LINGER_MS keeps the increase/decrease style a moment after the end.
const CLASS_LINGER_MS = 200;

type Frame = { score: number; className: string };

// AnimatedScore shows a player's score; when it changes, it counts from the
// previous value to the new one with an increase/decrease style. Between
// animations it shows the score as it is.
export default function AnimatedScore({
  previousScores,
  currentScores,
  playerId,
  className = '',
}: AnimatedScoreProps) {
  const previousScore = previousScores[playerId] ?? 0;
  const currentScore = currentScores[playerId] ?? 0;

  // frame is the animation's current value, null when not animating. It
  // changes only from animation frames and timers, never while rendering.
  const [frame, setFrame] = useState<Frame | null>(null);

  useEffect(() => {
    if (currentScore === previousScore) return;

    const style =
      currentScore > previousScore ? 'score-increase' : 'score-decrease';
    let startTime: number | null = null;
    let request = 0;
    let timer: ReturnType<typeof setTimeout> | undefined;

    const animate = (timestamp: number) => {
      startTime ??= timestamp;
      const progress = Math.min((timestamp - startTime) / ANIMATION_MS, 1);
      const easeOutQuart = 1 - Math.pow(1 - progress, 4);
      const score = Math.round(
        previousScore + (currentScore - previousScore) * easeOutQuart
      );
      setFrame({ score, className: style });

      if (progress < 1) {
        request = requestAnimationFrame(animate);
      } else {
        timer = setTimeout(() => setFrame(null), CLASS_LINGER_MS);
      }
    };
    request = requestAnimationFrame(animate);

    return () => {
      cancelAnimationFrame(request);
      clearTimeout(timer);
    };
    // Animate only when currentScore changes; previousScore is read as the
    // starting value, and a change to it alone must not restart the animation.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentScore]);

  return (
    <span className={`${className} ${frame?.className ?? ''}`}>
      {frame?.score ?? currentScore}
    </span>
  );
}
