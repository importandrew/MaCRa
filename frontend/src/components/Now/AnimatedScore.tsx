import { useEffect, useRef, useState } from 'react';

interface AnimatedScoreProps {
  score: number;
  className?: string;
}

export default function AnimatedScore({
  score,
  className = '',
}: AnimatedScoreProps) {
  const [displayScore, setDisplayScore] = useState(score);
  const [animationClass, setAnimationClass] = useState('');
  const prevScoreRef = useRef(score);
  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const animationFrameRef = useRef<number | null>(null);
  const startTimeRef = useRef<number | null>(null);
  const animationDuration = 2000; // 3 second

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
    const prevScore = prevScoreRef.current;

    // Only trigger animation if score actually changed
    if (score !== prevScore) {
      // Clean up any existing animations
      cleanup();

      // Determine animation direction
      const isIncreasing = score > prevScore;
      setAnimationClass(isIncreasing ? 'score-increase' : 'score-decrease');

      // Start the counting animation
      animateScoreChange(prevScore, score);

      // Update previous score reference
      prevScoreRef.current = score;
    }

    // Cleanup on unmount
    return () => {
      cleanup();
    };
  }, [score]);

  return (
    <span className={`${className} ${animationClass}`}>{displayScore}</span>
  );
}
