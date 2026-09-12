import {
  SortDescriptor,
  Table,
  TableBody,
  TableCell,
  TableColumn,
  TableHeader,
  TableRow,
} from '@heroui/react';
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

  const ratings = data?.ratings ?? [];

  const sortedRatings = useMemo(() => {
    const column = (sortDescriptor.column as RatingColumn) || 'rank';
    const direction = sortDescriptor.direction === 'descending' ? -1 : 1;
    return [...ratings].sort(
      (a, b) => compareRatings(a, b, column) * direction
    );
  }, [ratings, sortDescriptor]);

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
            title={`Rating — ${event.name}`}
          >
            Rating — {event.name}
          </h1>

          <Table
            isStriped
            removeWrapper
            className="z-0"
            aria-label="Event rating"
            sortDescriptor={sortDescriptor}
            onSortChange={setSortDescriptor}
          >
            <TableHeader>
              <TableColumn key="rank" allowsSorting align="center" width="10%">
                #
              </TableColumn>
              <TableColumn key="name" allowsSorting align="center">
                Player
              </TableColumn>
              <TableColumn key="games" allowsSorting align="center" width="15%">
                Games
              </TableColumn>
              <TableColumn key="tp" allowsSorting align="center" width="15%">
                TP
              </TableColumn>
              <TableColumn key="gp" allowsSorting align="center" width="15%">
                GP
              </TableColumn>
            </TableHeader>
            <TableBody emptyContent="No games played yet" items={sortedRatings}>
              {(item: EventRatingEntry) => (
                <TableRow key={item.player_id}>
                  <TableCell>{item.rank}</TableCell>
                  <TableCell>{item.name}</TableCell>
                  <TableCell>{item.games}</TableCell>
                  <TableCell>{item.tp}</TableCell>
                  <TableCell>{item.gp}</TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      )}
    </LoadingTemplate>
  );
}
