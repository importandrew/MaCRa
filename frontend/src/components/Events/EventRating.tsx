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

// RatingTable is one block of the rating: the ranked players, or the ones
// below the event's minimum games (no place: "—").
function RatingTable({
  label,
  items,
  sortDescriptor,
  onSortChange,
  emptyText,
}: {
  label: string;
  items: EventRatingEntry[];
  sortDescriptor: SortDescriptor;
  onSortChange: (descriptor: SortDescriptor) => void;
  emptyText?: string;
}) {
  const { t } = useTranslation();
  return (
    <Table className="z-0">
      <Table.ScrollContainer>
        <Table.Content
          aria-label={label}
          sortDescriptor={sortDescriptor}
          onSortChange={onSortChange}
        >
          <Table.Header>
            <Table.Column
              id="rank"
              allowsSorting
              className="text-center w-[10%]"
            >
              {({ sortDirection }) => (
                <Table.SortableColumnHeader sortDirection={sortDirection}>
                  #
                </Table.SortableColumnHeader>
              )}
            </Table.Column>
            <Table.Column
              id="name"
              allowsSorting
              isRowHeader
              className="text-center "
            >
              {({ sortDirection }) => (
                <Table.SortableColumnHeader sortDirection={sortDirection}>
                  {t('rating.player')}
                </Table.SortableColumnHeader>
              )}
            </Table.Column>
            <Table.Column
              id="games"
              allowsSorting
              className="text-center w-[15%]"
            >
              {({ sortDirection }) => (
                <Table.SortableColumnHeader sortDirection={sortDirection}>
                  {t('rating.games')}
                </Table.SortableColumnHeader>
              )}
            </Table.Column>
            <Table.Column id="tp" allowsSorting className="text-center w-[15%]">
              {({ sortDirection }) => (
                <Table.SortableColumnHeader sortDirection={sortDirection}>
                  {t('rating.tp')}
                </Table.SortableColumnHeader>
              )}
            </Table.Column>
            <Table.Column id="gp" allowsSorting className="text-center w-[15%]">
              {({ sortDirection }) => (
                <Table.SortableColumnHeader sortDirection={sortDirection}>
                  {t('rating.gp')}
                </Table.SortableColumnHeader>
              )}
            </Table.Column>
          </Table.Header>
          <Table.Body
            items={items}
            renderEmptyState={() => (
              <p className="py-6 text-center text-muted">{emptyText}</p>
            )}
          >
            {(item: EventRatingEntry) => (
              <Table.Row
                id={item.player_id}
                className="text-center even:bg-default/50"
              >
                <Table.Cell>{item.rank ?? '—'}</Table.Cell>
                <Table.Cell>{item.name}</Table.Cell>
                <Table.Cell>{item.games}</Table.Cell>
                <Table.Cell>{item.tp}</Table.Cell>
                <Table.Cell>{item.gp}</Table.Cell>
              </Table.Row>
            )}
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

  // Both blocks sort by the same column. Rows below the minimum have no
  // place: sorting by it keeps their order (the rating's).
  const sortedRatings = useMemo(
    () => sortRatings(data?.ratings ?? [], sortDescriptor),
    [data, sortDescriptor]
  );
  const sortedBelow = useMemo(
    () => sortRatings(data?.below_minimum ?? [], sortDescriptor),
    [data, sortDescriptor]
  );

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
            items={sortedRatings}
            sortDescriptor={sortDescriptor}
            onSortChange={setSortDescriptor}
            emptyText={t('rating.noGames')}
          />
          {sortedBelow.length > 0 && (
            <section className="flex flex-col gap-2">
              <h2 className="font-medium">
                {t('rating.belowMinimum', { count: data?.min_games ?? 0 })}
              </h2>
              <RatingTable
                label={t('rating.belowMinimum', {
                  count: data?.min_games ?? 0,
                })}
                items={sortedBelow}
                sortDescriptor={sortDescriptor}
                onSortChange={setSortDescriptor}
              />
            </section>
          )}
        </div>
      )}
    </LoadingTemplate>
  );
}
