import { useTranslation } from 'react-i18next';
import { useEffect, useState } from 'react';
import type { GameTimer } from '../models/GameTimer';

// HAND_WINDS: the round wind of hands 1–4, 5–8, 9–12, 13–16.
const HAND_WINDS = ['東', '南', '西', '北'];

// handKanji names a hand as the table does, e.g. hand 6 → 南2.
function handKanji(hand: number): string {
  return `${HAND_WINDS[Math.floor((hand - 1) / 4) % 4]}${((hand - 1) % 4) + 1}`;
}

// formatLeft shows the time left as m:ss, or h:mm:ss from an hour.
function formatLeft(ms: number): string {
  const total = Math.max(0, Math.ceil(ms / 1000));
  const h = Math.floor(total / 3600);
  const m = Math.floor((total % 3600) / 60);
  const s = String(total % 60).padStart(2, '0');
  return h > 0 ? `${h}:${String(m).padStart(2, '0')}:${s}` : `${m}:${s}`;
}

// TourTimer is a tournament game's countdown on its table and view. It is
// the one thing on a table screen that changes by itself (area 8): it ticks
// every second. Once the time is up it says what's left to play, as far as
// the last load knows (the table refreshes to learn the last hand).
export default function TourTimer({ timer }: { timer: GameTimer }) {
  const { t } = useTranslation();
  const endsAt = new Date(timer.ends_at).getTime();
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    const id = window.setInterval(() => setNow(Date.now()), 1000);
    return () => window.clearInterval(id);
  }, []);

  const left = endsAt - now;
  if (left > 0) {
    return (
      <p className="text-lg tabular-nums" aria-label={t('timer.left')}>
        {formatLeft(left)}
      </p>
    );
  }

  return (
    <p className="text-center text-sm font-medium text-danger">
      {timer.stopped
        ? t('timer.stopped')
        : timer.last_hand
          ? t('timer.lastHand', { hand: handKanji(timer.last_hand) })
          : t('timer.up')}
    </p>
  );
}
