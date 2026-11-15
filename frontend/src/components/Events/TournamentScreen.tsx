import { useTranslation } from 'react-i18next';
import { Button } from '@heroui/react';
import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import type { TourState } from '../../models/Tour';
import type { TourRating } from '../../models/TourRating';
import { useEventPage } from '../../providers/event/eventLayoutContext';
import { useFetchTours } from '../../providers/event/tournamentApi';
import { useFetchTourRating } from '../../providers/event/eventApi';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import SegmentedPicker from '../SegmentedPicker';
import { formatLeft } from '../../utils/time-left';

// WIND_KANJI: the starting seats as the table shows them.
const WIND_KANJI = { east: '東', south: '南', west: '西', north: '北' };

// REFRESH_MS: how often the screen reloads the tour and the rating (it
// refreshes by itself, area 8); ALTERNATE_MS: how long "auto" shows the
// seating or the rating before switching.
const REFRESH_MS = 30_000;
const ALTERNATE_MS = 20_000;

type View = 'auto' | 'seating' | 'rating';

// TourHeader is the screen's top line: the tour and its state, and the
// countdown while it is live.
function TourHeader({ state, now }: { state: TourState; now: number }) {
  const { t } = useTranslation();
  const tour = state.current;
  if (!tour || tour.status === 'draft') {
    return (
      <p className="text-[1.6em] font-semibold">
        {tour ? t('screen.preparing', { n: tour.number }) : t('screen.noTours')}
      </p>
    );
  }

  let clock: string | null = null;
  let isUp = false;
  if (tour.status === 'live' && tour.started_at && tour.minutes) {
    const left =
      new Date(tour.started_at).getTime() + tour.minutes * 60_000 - now;
    isUp = left <= 0;
    clock = isUp ? t('screen.timeUp') : formatLeft(left);
  }

  return (
    <div className="flex flex-wrap items-baseline justify-between gap-x-6 gap-y-1">
      <p className="text-[1.6em] font-semibold">
        {t('screen.tour', { n: tour.number })}
        <span className="font-normal text-muted">
          {' · '}
          {t(`screen.status.${tour.status}`)}
        </span>
      </p>
      {clock && (
        <p
          className={`text-[3em] font-semibold leading-none tabular-nums ${isUp ? 'text-danger' : ''}`}
          aria-label={t('timer.left')}
        >
          {clock}
        </p>
      )}
    </div>
  );
}

// Seating is the tour's tables in a grid, so players find their table: the
// four seats east first, with participant numbers; during the tour also
// which tables have ended and the time added to a table.
function Seating({ state }: { state: TourState }) {
  const { t } = useTranslation();
  const tour = state.current;
  if (!tour || tour.status === 'draft' || state.tables.length === 0) {
    return null;
  }
  return (
    <ul
      aria-label={t('screen.seating')}
      className="grid grid-cols-[repeat(auto-fill,minmax(13em,1fr))] gap-[0.75em]"
    >
      {state.tables.map((table) => (
        <li
          key={table.table}
          className={`rounded-xl p-[0.6em] ${table.status === 'ended' ? 'bg-default/40 text-muted' : 'bg-surface'}`}
        >
          <p className="mb-[0.3em] flex justify-between gap-2 font-semibold">
            <span>{t('tours.table', { n: table.table })}</span>
            <span className="font-normal text-muted">
              {tour.status === 'live' && table.status === 'ended'
                ? t('screen.tableEnded')
                : table.extra_minutes > 0
                  ? t('tours.extra', { n: table.extra_minutes })
                  : ''}
            </span>
          </p>
          <ul>
            {table.seats.map((seat) => (
              <li key={seat.user_id} className="flex items-center gap-[0.5em]">
                <span className="w-[1.2em] shrink-0 text-center">
                  {WIND_KANJI[seat.wind]}
                </span>
                <span className="w-[1.8em] shrink-0 text-end tabular-nums text-muted">
                  {seat.number}
                </span>
                <span className="min-w-0 truncate">{seat.name}</span>
              </li>
            ))}
          </ul>
        </li>
      ))}
    </ul>
  );
}

// Standings is the rating after the last finished tour, in columns so a
// long list fits a TV: place, name, TP, GP.
function Standings({ rating }: { rating: TourRating }) {
  const { t } = useTranslation();
  const last = rating.tours[rating.tours.length - 1];
  if (last == null) return null;
  const ranked = rating.rows.filter((row) => row.rank != null);
  return (
    <section className="flex flex-col gap-[0.5em]">
      <h2 className="text-[1.2em] font-semibold">
        {t('tourRating.after', { n: last })}
      </h2>
      <ol
        aria-label={t('events.rating')}
        className="columns-[16em] gap-[1.5em] tabular-nums"
      >
        {ranked.map((row) => (
          <li
            key={row.player_id}
            className="flex break-inside-avoid gap-[0.5em] border-b border-border py-[0.15em]"
          >
            <span className="w-[2em] shrink-0 text-end text-muted">
              {row.rank}
            </span>
            <span className="min-w-0 flex-1 truncate">{row.name}</span>
            <span className="w-[2.5em] shrink-0 text-end font-semibold">
              {row.tp}
            </span>
            <span className="w-[3em] shrink-0 text-end text-muted">
              {row.gp}
            </span>
          </li>
        ))}
      </ol>
    </section>
  );
}

// TournamentScreen is /event/:eventId/screen, the common screen for the
// hall (projector, TV), for the tournament's staff: the current tour's
// seating, its countdown and the rating. Unlike the rest of the app it
// refreshes by itself (area 8), and "auto" alternates the seating and the
// rating. "Full screen" leaves out the app's header.
export default function TournamentScreen() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { eventId, event, canView } = useEventPage('view_tournament');
  const tours = useFetchTours(canView ? eventId : null);
  const rating = useFetchTourRating(canView ? eventId : null);
  const screen = useRef<HTMLDivElement>(null);
  const [view, setView] = useState<View>('auto');
  const [now, setNow] = useState(() => Date.now());

  // The countdown ticks every second; the data reloads every REFRESH_MS.
  const { refetch: refetchTours } = tours;
  const { refetch: refetchRating } = rating;
  useEffect(() => {
    const tick = window.setInterval(() => setNow(Date.now()), 1000);
    const reload = window.setInterval(() => {
      refetchTours();
      refetchRating();
    }, REFRESH_MS);
    return () => {
      window.clearInterval(tick);
      window.clearInterval(reload);
    };
  }, [refetchTours, refetchRating]);

  const state = tours.data;
  const hasSeating =
    !!state?.current &&
    state.current.status !== 'draft' &&
    state.tables.length > 0;
  const hasRating = (rating.data?.tours.length ?? 0) > 0;
  // "auto" alternates when there are both, else shows what there is.
  const autoShowsRating =
    hasRating && (!hasSeating || Math.floor(now / ALTERNATE_MS) % 2 === 1);
  const showRating = view === 'rating' || (view === 'auto' && autoShowsRating);

  const fullScreen = () => {
    screen.current?.requestFullscreen?.().catch(() => undefined);
  };

  return (
    <LoadingTemplate
      isLoading={tours.isLoading || rating.isLoading}
      isError={tours.isError || rating.isError}
      isForbidden={!canView}
      onRefresh={() => Promise.all([refetchTours(), refetchRating()])}
      onGoHome={() => navigate(`/event/${eventId}`)}
    >
      {event && state && (
        <div className="flex w-full min-w-0 flex-col gap-3 px-3 md:my-3 md:px-6">
          <div className="flex flex-wrap items-center gap-3">
            <SegmentedPicker
              ariaLabel={t('screen.viewLabel')}
              options={(['auto', 'seating', 'rating'] as const).map(
                (value) => ({ value, label: t(`screen.view.${value}`) })
              )}
              selectedKey={view}
              onSelectionChange={setView}
            />
            <Button variant="tertiary" onPress={fullScreen}>
              {t('screen.fullScreen')}
            </Button>
          </div>

          {/* The part that goes full screen: sized by the screen's width, so
              a TV gets big type and a phone a readable one. */}
          <div
            ref={screen}
            className="flex min-h-0 flex-col gap-[1em] overflow-auto bg-background p-[1em] text-[clamp(14px,1.35vw,40px)]"
          >
            <p className="text-[0.9em] text-muted">{event.name}</p>
            <TourHeader state={state} now={now} />
            {showRating && rating.data ? (
              <Standings rating={rating.data} />
            ) : hasSeating ? (
              <Seating state={state} />
            ) : (
              <p className="text-muted">{t('screen.nothingYet')}</p>
            )}
          </div>
        </div>
      )}
    </LoadingTemplate>
  );
}
