import { useTranslation } from 'react-i18next';
import { Button, Table } from '@heroui/react';
import { useState } from 'react';
import type { TourRating, TourRatingRow } from '../../models/TourRating';
import SegmentedPicker from '../SegmentedPicker';
import { downloadText, tourRatingCsv } from './tourRatingCsv';

// CELL: the narrow padding of the rating's cells (as the event rating).
const CELL = 'px-1';

// signed writes GP with its sign: "+32", "-3", "0".
function signed(gp: number): string {
  return gp > 0 ? `+${gp}` : `${gp}`;
}

// Row is a row of the tables: a player, or a separator above the players
// who haven't played yet or who withdrew.
type Row =
  | { id: number; kind: 'player'; row: TourRatingRow; isMuted: boolean }
  | { id: string; kind: 'separator'; label: string };

// ratingRows lists the players with separators before the ones who haven't
// played a finished tour yet and the ones who withdrew (rows come grouped).
function ratingRows(rating: TourRating, labels: [string, string]): Row[] {
  const rows: Row[] = [];
  let group = 0; // 0 ranked, 1 not played, 2 withdrew
  for (const row of rating.rows) {
    const rowGroup = row.withdrawn ? 2 : row.rank == null ? 1 : 0;
    if (rowGroup > group) {
      group = rowGroup;
      rows.push({
        id: `separator-${group}`,
        kind: 'separator',
        label: labels[group - 1],
      });
    }
    rows.push({
      id: row.player_id,
      kind: 'player',
      row,
      isMuted: rowGroup > 0,
    });
  }
  return rows;
}

// TourRatingView is a tournament's rating (area 8, "Rating"): the totals, or
// TP and GP of each finished tour (the place and name stay while the tours
// scroll sideways), "Tour N in progress", and the CSV export.
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
  const rows = ratingRows(rating, [
    t('tourRating.notPlayed'),
    t('tourRating.withdrew'),
  ]);
  const byTour = view === 'tours' && rating.tours.length > 0;
  const columns = 2 + (byTour ? rating.tours.length : 0) + 2;

  const exportCsv = () =>
    downloadText(
      tourRatingCsv(rating, t),
      `${eventName.replace(/[\\/:*?"<>|]+/g, '_')} — ${t('tourRating.fileName')}.csv`,
      'text/csv;charset=utf-8'
    );

  // Sticky cells keep the place and the name in view while the tours
  // scroll; they need a background of their own: the header's (the table's
  // frame) or the rows'.
  const stickyHead = byTour ? 'sticky z-10 bg-surface-secondary' : '';
  const sticky = byTour ? 'sticky z-10 bg-surface' : '';

  return (
    <div className="flex flex-col gap-3">
      <p className="text-sm text-muted">
        {rating.tours.length > 0
          ? t('tourRating.after', { n: rating.tours[rating.tours.length - 1] })
          : t('tourRating.noTours')}
        {rating.in_progress > 0 &&
          ` · ${t('tourRating.inProgress', { n: rating.in_progress })}`}
      </p>
      {rating.tours.length > 0 && (
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

      <Table className="z-0">
        <Table.ScrollContainer>
          <Table.Content
            // Totals fit the screen (the name takes the room left); by tour
            // the table is as wide as its columns and scrolls.
            className={byTour ? 'min-w-max' : 'table-fixed'}
            aria-label={t('events.rating')}
          >
            <Table.Header>
              <Table.Column
                id="rank"
                className={`${CELL} w-8 left-0 ${stickyHead}`}
              >
                #
              </Table.Column>
              <Table.Column
                id="name"
                isRowHeader
                className={`${CELL} ${byTour ? 'left-8 w-28 max-w-28' : ''} ${stickyHead}`}
              >
                {t('rating.name')}
              </Table.Column>
              {byTour
                ? rating.tours.map((tour) => (
                    <Table.Column
                      key={tour}
                      id={`tour-${tour}`}
                      className={`${CELL} text-center`}
                    >
                      {t('tourRating.tourColumn', { n: tour })}
                    </Table.Column>
                  ))
                : []}
              <Table.Column id="tp" className={`${CELL} w-11 text-center`}>
                {t('rating.tp')}
              </Table.Column>
              <Table.Column id="gp" className={`${CELL} w-12 text-center`}>
                {t('rating.gp')}
              </Table.Column>
            </Table.Header>
            <Table.Body
              items={rows}
              renderEmptyState={() => (
                <p className="py-6 text-center text-muted">
                  {t('rating.noGames')}
                </p>
              )}
            >
              {(item: Row) => {
                if (item.kind === 'separator') {
                  return (
                    <Table.Row id={item.id}>
                      <Table.Cell
                        colSpan={columns}
                        className={`${CELL} bg-default py-1.5 text-start text-xs font-medium text-muted`}
                      >
                        {/* Sticky, so the label stays in view while the
                            tours scroll. */}
                        <span className="sticky left-1">{item.label}</span>
                      </Table.Cell>
                    </Table.Row>
                  );
                }
                const { row } = item;
                const cell = `${CELL} ${item.isMuted ? 'text-muted' : ''}`;
                return (
                  <Table.Row
                    id={item.id}
                    href={playerHref(row.player_id)}
                    className="cursor-pointer text-center"
                  >
                    <Table.Cell className={`${cell} left-0 ${sticky}`}>
                      {row.rank ?? '—'}
                    </Table.Cell>
                    <Table.Cell className={`${cell} left-8 ${sticky}`}>
                      <span
                        title={row.name}
                        className={`block truncate ${byTour ? 'max-w-28 text-start' : 'text-center'}`}
                      >
                        {row.name}
                      </span>
                    </Table.Cell>
                    {byTour
                      ? row.tours.map((result, i) => (
                          <Table.Cell
                            key={rating.tours[i]}
                            className={`${cell} whitespace-nowrap tabular-nums`}
                          >
                            {result ? (
                              <>
                                {result.tp}
                                <span className="text-muted">
                                  {' / '}
                                  {signed(result.gp)}
                                </span>
                              </>
                            ) : (
                              '—'
                            )}
                          </Table.Cell>
                        ))
                      : []}
                    <Table.Cell className={`${cell} tabular-nums`}>
                      {row.tp}
                    </Table.Cell>
                    <Table.Cell className={`${cell} tabular-nums`}>
                      {row.gp}
                    </Table.Cell>
                  </Table.Row>
                );
              }}
            </Table.Body>
          </Table.Content>
        </Table.ScrollContainer>
      </Table>
      {byTour && (
        <p className="text-xs text-muted">{t('tourRating.byTourHint')}</p>
      )}

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
