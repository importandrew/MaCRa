import { useTranslation } from 'react-i18next';
import { Button, Table } from '@heroui/react';
import { useState } from 'react';
import type { TourRating } from '../../models/TourRating';
import SegmentedPicker from '../SegmentedPicker';
import { downloadText, tourRatingCsv } from './tourRatingCsv';

// CELL: the narrow padding of the rating's cells (as the event rating).
const CELL = 'px-1';

// Line is a row of the table: a player's place (null: none), TP and GP, in
// total or in one tour; isOut: at the bottom, greyed (withdrew, or no result
// to rank).
type Line = {
  id: number;
  name: string;
  rank: number | null;
  tp: number;
  gp: number;
  isOut: boolean;
};

// totalLines are the totals as the server ranks them; players without a
// place (haven't played a finished tour, or withdrew) come last.
function totalLines(rating: TourRating): Line[] {
  return rating.rows.map((row) => ({
    id: row.player_id,
    name: row.name,
    rank: row.rank,
    tp: row.tp,
    gp: row.gp,
    isOut: row.rank == null,
  }));
}

// tourLines ranks the players by their result in one tour (the index in
// rating.tours): TP, then GP, equal in both sharing the place, so the best
// result of the tour is on top. Players who didn't play it come last; those
// who withdrew stay greyed.
function tourLines(rating: TourRating, index: number): Line[] {
  const played: Line[] = [];
  const out: Line[] = [];
  for (const row of rating.rows) {
    const result = row.tours[index];
    const line: Line = {
      id: row.player_id,
      name: row.name,
      rank: null,
      tp: result?.tp ?? 0,
      gp: result?.gp ?? 0,
      isOut: row.withdrawn || result == null,
    };
    (result ? played : out).push(line);
  }
  played.sort((a, b) => b.tp - a.tp || b.gp - a.gp);
  played.forEach((line, i) => {
    const before = played[i - 1];
    line.rank =
      before && before.tp === line.tp && before.gp === line.gp
        ? before.rank
        : i + 1;
  });
  return [...played, ...out];
}

// TourRatingView is a tournament's rating (area 8, "Rating"): the totals, or
// one finished tour's results ranked (an organizer often rewards the best
// result of a tour), "tour N in progress", and the CSV export.
export default function TourRatingView({
  rating,
  eventName,
  playerHref,
}: {
  rating: TourRating;
  eventName: string;
  // playerHref: the link of a player's statistics (a tap on their row).
  playerHref: (playerId: number) => string;
}) {
  const { t } = useTranslation();
  const [view, setView] = useState<'totals' | 'tours'>('totals');
  const lastTour = rating.tours[rating.tours.length - 1];
  // pickedTour: the tour shown by tour; null: the last finished one.
  const [pickedTour, setPickedTour] = useState<number | null>(null);
  const tour =
    pickedTour != null && rating.tours.includes(pickedTour)
      ? pickedTour
      : lastTour;
  const byTour = view === 'tours' && tour != null;
  const lines = byTour
    ? tourLines(rating, rating.tours.indexOf(tour))
    : totalLines(rating);

  const exportCsv = () =>
    downloadText(
      tourRatingCsv(rating, t),
      `${eventName.replace(/[\\/:*?"<>|]+/g, '_')} — ${t('tourRating.fileName')}.csv`,
      'text/csv;charset=utf-8'
    );

  return (
    <div className="flex flex-col gap-3">
      <p className="text-sm text-muted">
        {lastTour != null
          ? t('tourRating.after', { n: lastTour })
          : t('tourRating.noTours')}
        {rating.in_progress > 0 &&
          ` · ${t('tourRating.inProgress', { n: rating.in_progress })}`}
      </p>
      {lastTour != null && (
        <SegmentedPicker
          ariaLabel={t('tourRating.viewLabel')}
          options={[
            { value: 'totals', label: t('tourRating.totals') },
            { value: 'tours', label: t('tourRating.byTour') },
          ]}
          selectedKey={view}
          onSelectionChange={setView}
        />
      )}
      {byTour && (
        <SegmentedPicker
          ariaLabel={t('tourRating.tourLabel')}
          options={rating.tours.map((n) => ({
            value: String(n),
            label: String(n),
          }))}
          selectedKey={String(tour)}
          onSelectionChange={(n) => setPickedTour(Number(n))}
          optionClassName="px-0"
        />
      )}

      <Table className="z-0">
        <Table.ScrollContainer>
          {/* Fixed layout: the name takes the room the others leave. */}
          <Table.Content
            className="table-fixed"
            aria-label={
              byTour
                ? t('tourRating.tourColumn', { n: tour })
                : t('events.rating')
            }
          >
            <Table.Header>
              <Table.Column id="rank" className={`${CELL} w-8`}>
                #
              </Table.Column>
              <Table.Column id="name" isRowHeader className={CELL}>
                {t('rating.name')}
              </Table.Column>
              <Table.Column id="tp" className={`${CELL} w-11 text-center`}>
                {t('rating.tp')}
              </Table.Column>
              <Table.Column id="gp" className={`${CELL} w-12 text-center`}>
                {t('rating.gp')}
              </Table.Column>
            </Table.Header>
            <Table.Body
              items={lines}
              renderEmptyState={() => (
                <p className="py-6 text-center text-muted">
                  {t('rating.noGames')}
                </p>
              )}
            >
              {(line: Line) => {
                // The cells carry the colors: HeroUI styles each cell.
                const cell = `${CELL} ${line.isOut ? 'bg-default text-muted' : ''}`;
                const noResult = byTour && line.rank == null;
                return (
                  <Table.Row
                    id={line.id}
                    href={playerHref(line.id)}
                    className="cursor-pointer text-center"
                  >
                    <Table.Cell className={cell}>{line.rank ?? '—'}</Table.Cell>
                    <Table.Cell className={cell}>
                      <span title={line.name} className="block truncate">
                        {line.name}
                      </span>
                    </Table.Cell>
                    <Table.Cell className={`${cell} tabular-nums`}>
                      {noResult ? '—' : line.tp}
                    </Table.Cell>
                    <Table.Cell className={`${cell} tabular-nums`}>
                      {noResult ? '—' : line.gp}
                    </Table.Cell>
                  </Table.Row>
                );
              }}
            </Table.Body>
          </Table.Content>
        </Table.ScrollContainer>
      </Table>

      {rating.rows.length > 0 && (
        <Button
          variant="tertiary"
          className="w-3/4 self-center"
          onPress={exportCsv}
        >
          {t('tourRating.exportCsv')}
        </Button>
      )}
    </div>
  );
}
