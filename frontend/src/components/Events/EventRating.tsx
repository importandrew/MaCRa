import { useTranslation } from 'react-i18next';
import { SortDescriptor, Table } from '@heroui/react';
import { useCallback, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { EventRatingEntry } from '../../models/EventRating';
import { useFetchEventRatings } from '../../providers/event/eventApi';
import { useEventPage } from '../../providers/event/eventLayoutContext';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';

type RatingColumn = 'rank' | 'name' | 'games' | 'tp' | 'gp';

function compareRatings(
  a: EventRatingEntry,
  b: EventRatingEntry,
  column: RatingColumn
): number {
  switch (column) {
    case 'name':
      return a.name.localeCompare(b.name);
    case 'rank':
      return (a.rank ?? 0) - (b.rank ?? 0);
    case 'games':
      return a.games - b.games;
    case 'tp':
      return a.tp - b.tp;
    case 'gp':
      return a.gp - b.gp;
  }
}

function sortRatings(
  ratings: EventRatingEntry[],
  sortDescriptor: SortDescriptor
): EventRatingEntry[] {
  const column = (sortDescriptor.column as RatingColumn) || 'rank';
  const direction = sortDescriptor.direction === 'descending' ? -1 : 1;
  return [...ratings].sort((a, b) => compareRatings(a, b, column) * direction);
}

// PlayerName is a name in the rating: one line with an ellipsis, so every row
// has the same height. The full name shows on hover (title) and as the title
// of the player's statistics, which a tap on the row opens.
function PlayerName({ name }: { name: string }) {
  return (
    <span title={name} className="block truncate text-center">
      {name}
    </span>
  );
}

// CELL_CLASS: HeroUI's px-4 leaves too little room for five columns at 320px
// with a fixed layout; the header uses the same padding.
const CELL_CLASS = 'px-1';

// RATING_COLUMNS: the numbers get fixed widths that fit their values and
// their header with the sort arrow ("Игры ˄", "-123"); the name takes the
// rest. Labels: rating.<id>.
const RATING_COLUMNS: { id: RatingColumn; width: string }[] = [
  { id: 'rank', width: 'w-8' },
  { id: 'name', width: '' },
  { id: 'games', width: 'w-[3.25rem]' },
  { id: 'tp', width: 'w-11' },
  { id: 'gp', width: 'w-12' },
];

// RatingRow is a row of the rating table: a player, or the separator above
// the players without a place.
type RatingRow =
  | { id: number; kind: 'player'; entry: EventRatingEntry; isBelow: boolean }
  | { id: 'separator'; kind: 'separator'; label: string };

// RatingTable is the rating: the ranked players, then a separator row and
// the players below the event's minimum games, greyed, without a place
// ("—").
function RatingTable({
  label,
  rows,
  sortDescriptor,
  onSortChange,
  emptyText,
  playerHref,
}: {
  label: string;
  rows: RatingRow[];
  sortDescriptor: SortDescriptor;
  onSortChange: (descriptor: SortDescriptor) => void;
  emptyText: string;
  // playerHref: the link of a player's statistics (a tap on their row).
  playerHref: (playerId: number) => string;
}) {
  const { t } = useTranslation();
  return (
    <Table className="z-0">
      <Table.ScrollContainer>
        {/* Fixed layout: the name column takes the room the others leave and
            never widens the table past the screen. */}
        <Table.Content
          className="table-fixed"
          aria-label={label}
          sortDescriptor={sortDescriptor}
          onSortChange={onSortChange}
        >
          <Table.Header>
            {RATING_COLUMNS.map((column) => (
              <Table.Column
                key={column.id}
                id={column.id}
                allowsSorting
                isRowHeader={column.id === 'name'}
                className={`${CELL_CLASS} ${column.width}`}
              >
                {({ sortDirection }) => (
                  <Table.SortableColumnHeader
                    sortDirection={sortDirection}
                    className="justify-center gap-0.5"
                  >
                    {column.id === 'rank' ? '#' : t(`rating.${column.id}`)}
                  </Table.SortableColumnHeader>
                )}
              </Table.Column>
            ))}
          </Table.Header>
          <Table.Body
            items={rows}
            renderEmptyState={() => (
              <p className="py-6 text-center text-muted">{emptyText}</p>
            )}
          >
            {(row: RatingRow) => {
              if (row.kind === 'separator') {
                return (
                  <Table.Row id={row.id}>
                    <Table.Cell
                      colSpan={5}
                      className={`${CELL_CLASS} bg-default py-1.5 text-start text-xs font-medium text-muted`}
                    >
                      {row.label}
                    </Table.Cell>
                  </Table.Row>
                );
              }
              // The cells, not the row, carry the colors: HeroUI styles
              // each cell's background and text.
              const cell = `${CELL_CLASS} ${row.isBelow ? 'text-muted' : ''}`;
              return (
                // A tap opens the player's statistics in the event.
                <Table.Row
                  id={row.id}
                  href={playerHref(row.entry.player_id)}
                  className="cursor-pointer text-center"
                >
                  <Table.Cell className={cell}>
                    {row.entry.rank ?? '—'}
                  </Table.Cell>
                  <Table.Cell className={cell}>
                    <PlayerName name={row.entry.name} />
                  </Table.Cell>
                  <Table.Cell className={cell}>{row.entry.games}</Table.Cell>
                  <Table.Cell className={cell}>{row.entry.tp}</Table.Cell>
                  <Table.Cell className={cell}>{row.entry.gp}</Table.Cell>
                </Table.Row>
              );
            }}
          </Table.Body>
        </Table.Content>
      </Table.ScrollContainer>
    </Table>
  );
}

export default function EventRating() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const {
    eventId,
    event,
    isLoading: layoutLoading,
    canView: canViewRatings,
  } = useEventPage('view_ratings');

  const {
    data,
    isLoading: ratingsLoading,
    isError,
    isForbidden: ratingsForbidden,
    refetch,
  } = useFetchEventRatings(canViewRatings ? eventId : null);

  const [sortDescriptor, setSortDescriptor] = useState<SortDescriptor>({
    column: 'rank',
    direction: 'ascending',
  });

  const onRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

  const isLoading = layoutLoading || (!!canViewRatings && ratingsLoading);
  const forbidden = !canViewRatings || ratingsForbidden;

  // A column sorts each group on its own: the players below the minimum stay
  // under the separator whatever the order. They have no place, so sorting
  // by it keeps their order (the rating's). With a minimum of 1 the group
  // holds only the players who haven't played yet.
  const minGames = data?.min_games ?? 1;
  const belowLabel =
    minGames > 1
      ? t('rating.belowMinimum', { count: minGames })
      : t('rating.notPlayedYet');
  const rows = useMemo((): RatingRow[] => {
    const ranked = sortRatings(data?.ratings ?? [], sortDescriptor);
    const below = sortRatings(data?.below_minimum ?? [], sortDescriptor);
    return [
      ...ranked.map((entry) => ({
        id: entry.player_id,
        kind: 'player' as const,
        entry,
        isBelow: false,
      })),
      ...(below.length > 0
        ? [
            {
              id: 'separator' as const,
              kind: 'separator' as const,
              label: belowLabel,
            },
          ]
        : []),
      ...below.map((entry) => ({
        id: entry.player_id,
        kind: 'player' as const,
        entry,
        isBelow: true,
      })),
    ];
  }, [data, sortDescriptor, belowLabel]);

  return (
    <LoadingTemplate
      isLoading={isLoading}
      isError={isError}
      isForbidden={forbidden}
      onRefresh={onRefresh}
      onGoHome={() => navigate(`/event/${event?.id}`)}
    >
      {event && (
        <div className="flex flex-col w-full max-w-[600px] mx-auto px-3 md:px-0 md:my-3 gap-4 min-w-0">
          <h1
            className="text-2xl font-semibold truncate min-w-0"
            title={t('events.ratingTitle', { name: event.name })}
          >
            {t('events.ratingTitle', { name: event.name })}
          </h1>

          <RatingTable
            label={t('events.rating')}
            rows={rows}
            sortDescriptor={sortDescriptor}
            onSortChange={setSortDescriptor}
            emptyText={t('rating.noGames')}
            playerHref={(playerId) => `/event/${event.id}/player/${playerId}`}
          />
        </div>
      )}
    </LoadingTemplate>
  );
}
