import { useEffect, useRef, useState } from 'react';
import { AnimatedScoreProps } from '../../models/GameData';

export default function AnimatedScore({
  scoreData,
  playerId,
  className = '',
}: AnimatedScoreProps) {
  const [displayScore, setDisplayScore] = useState<number>(scoreData[playerId]);
  const [animationClass, setAnimationClass] = useState('');
  const prevScoreRef = useRef<number>(scoreData[playerId]);
  const prevPlayerIdRef = useRef<number>(playerId);
  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const animationFrameRef = useRef<number | null>(null);
  const startTimeRef = useRef<number | null>(null);
  const animationDuration = 2000; // 2 seconds

  // Clean up animation frames and timers
  const cleanup = () => {
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
    if (animationFrameRef.current) {
      cancelAnimationFrame(animationFrameRef.current);
      animationFrameRef.current = null;
    }
  };

  // Animation function for counting scores
  const animateScoreChange = (start: number, end: number) => {
    const animate = (timestamp: number) => {
      if (!startTimeRef.current) {
        startTimeRef.current = timestamp;
      }

      const elapsed = timestamp - startTimeRef.current;
      const progress = Math.min(elapsed / animationDuration, 1);

      // Easing function for smooth animation
      const easeOutQuart = 1 - Math.pow(1 - progress, 4);

      // Calculate current score based on progress
      const currentScore = Math.round(start + (end - start) * easeOutQuart);

      setDisplayScore(currentScore);

      // Continue animation if not complete
      if (progress < 1) {
        animationFrameRef.current = requestAnimationFrame(animate);
      } else {
        // Animation complete
        setDisplayScore(end);
        startTimeRef.current = null;

        // Clear animation class after a short delay
        timerRef.current = setTimeout(() => {
          setAnimationClass('');
        }, 200);
      }
    };

    animationFrameRef.current = requestAnimationFrame(animate);
  };

  useEffect(() => {
    const currentScore = scoreData[playerId];
    const prevScore = prevScoreRef.current;

    if (playerId !== prevPlayerIdRef.current) {
      setDisplayScore(currentScore);
      prevPlayerIdRef.current = playerId;
    }

    // Only trigger animation if score actually changed
    if (currentScore !== prevScore) {
      // Clean up any existing animations
      cleanup();

      // Determine animation direction
      const isIncreasing = currentScore > prevScore;
      setAnimationClass(isIncreasing ? 'score-increase' : 'score-decrease');

      // Start the counting animation
      animateScoreChange(prevScore, currentScore);
    }

    // Update ref for next render
    prevScoreRef.current = currentScore;

    // Cleanup on unmount
    return () => {
      cleanup();
    };
  }, [scoreData, playerId]);

  return (
    <span className={`${className} ${animationClass}`}>{displayScore}</span>
  );
}
