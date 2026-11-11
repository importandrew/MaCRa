import { useTranslation } from 'react-i18next';
import { SortDescriptor, Table } from '@heroui/react';
import { useCallback, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { EventRatingEntry } from '../../models/EventRating';
import { useFetchEventRatings } from '../../providers/event/eventApi';
import { useEventLayout } from '../../providers/event/eventLayoutContext';
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
      return a.rank - b.rank;
    case 'games':
      return a.games - b.games;
    case 'tp':
      return a.tp - b.tp;
    case 'gp':
      return a.gp - b.gp;
  }
}

export default function EventRating() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const {
    eventId,
    event,
    hasCapability,
    isLoading: layoutLoading,
    isForbidden: layoutForbidden,
  } = useEventLayout();

  const canViewRatings = hasCapability('view_ratings');
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
  const forbidden =
    layoutForbidden ||
    ratingsForbidden ||
    (!layoutLoading && !!event && !canViewRatings);

  const sortedRatings = useMemo(() => {
    const ratings = data?.ratings ?? [];
    const column = (sortDescriptor.column as RatingColumn) || 'rank';
    const direction = sortDescriptor.direction === 'descending' ? -1 : 1;
    return [...ratings].sort(
      (a, b) => compareRatings(a, b, column) * direction
    );
  }, [data, sortDescriptor]);

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

          <Table className="z-0">
            <Table.ScrollContainer>
              <Table.Content
                aria-label={t('events.rating')}
                sortDescriptor={sortDescriptor}
                onSortChange={setSortDescriptor}
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
                  <Table.Column
                    id="tp"
                    allowsSorting
                    className="text-center w-[15%]"
                  >
                    {({ sortDirection }) => (
                      <Table.SortableColumnHeader sortDirection={sortDirection}>
                        {t('rating.tp')}
                      </Table.SortableColumnHeader>
                    )}
                  </Table.Column>
                  <Table.Column
                    id="gp"
                    allowsSorting
                    className="text-center w-[15%]"
                  >
                    {({ sortDirection }) => (
                      <Table.SortableColumnHeader sortDirection={sortDirection}>
                        {t('rating.gp')}
                      </Table.SortableColumnHeader>
                    )}
                  </Table.Column>
                </Table.Header>
                <Table.Body
                  items={sortedRatings}
                  renderEmptyState={() => (
                    <p className="py-6 text-center text-muted">
                      {t('rating.noGames')}
                    </p>
                  )}
                >
                  {(item: EventRatingEntry) => (
                    <Table.Row
                      id={item.player_id}
                      className="text-center even:bg-default/50"
                    >
                      <Table.Cell>{item.rank}</Table.Cell>
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
        </div>
      )}
    </LoadingTemplate>
  );
}
