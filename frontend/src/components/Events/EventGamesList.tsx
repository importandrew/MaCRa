import { Button, CircularProgress, Listbox, ListboxItem } from '@heroui/react';
import { useNavigate } from 'react-router-dom';
import {
  EventGameListItem,
  EventGameListPlayer,
} from '../../models/EventGames';
import { ItemCounter, ListboxWrapper } from '../History/History';

type EventGamesListProps = {
  games: EventGameListItem[];
  emptyContent: string;
  ariaLabel: string;
  mode: 'live' | 'finished';
  onSelect: (gameId: number) => void;
  hasMore: boolean;
  isLoadingMore: boolean;
  onLoadMore: () => void;
};

const TOTAL_ROUNDS = 16;

function roundProgress(game: EventGameListItem, mode: 'live' | 'finished') {
  if (mode === 'finished' || game.status === 'ended' || !game.round) {
    return { value: 100, label: '✓', ariaLabel: 'Game ended' };
  }

  return {
    value: ((game.round.num - 1) / TOTAL_ROUNDS) * 100,
    label: game.round.kanji,
    ariaLabel: `Round ${game.round.kanji}`,
  };
}

function gameTextValue(game: EventGameListItem) {
  const players = game.players.map((p) => p.name).join(', ');
  return players ? `Game #${game.id}: ${players}` : `Game #${game.id}`;
}

function PlayerRow({ player }: { player: EventGameListPlayer }) {
  const hasWind = Boolean(player.wind);

  return (
    <div
      className={
        hasWind
          ? 'grid grid-cols-[1.25rem_minmax(0,1fr)_auto] gap-x-2 items-center text-base leading-5 min-w-0'
          : 'grid grid-cols-[minmax(0,1fr)_auto] gap-x-2 items-center text-base leading-5 min-w-0'
      }
    >
      {hasWind && (
        <span className="text-default-500 shrink-0">{player.wind}</span>
      )}
      <span className="truncate whitespace-nowrap min-w-0">{player.name}</span>
      <span className="tabular-nums text-right shrink-0 text-default-700">
        {player.gp}
      </span>
    </div>
  );
}

function GameListItemContent({
  game,
  mode,
}: {
  game: EventGameListItem;
  mode: 'live' | 'finished';
}) {
  const { value, label, ariaLabel } = roundProgress(game, mode);

  return (
    <div className="grid grid-cols-[auto_minmax(0,1fr)] gap-3 min-w-0 w-full items-center">
      {/* h-0: don't inflate row height; min-h-full + aspect-square: match players column */}
      <div className="h-0 min-h-full aspect-square">
        <CircularProgress
          aria-label={ariaLabel}
          value={value}
          color="primary"
          showValueLabel
          valueLabel={label}
          classNames={{
            base: 'h-full w-full max-w-none',
            svgWrapper: 'h-full w-full',
            svg: '!w-full !h-full',
            value: 'text-lg font-medium',
          }}
        />
      </div>
      <div className="flex flex-col min-w-0">
        {game.players.length > 0 ? (
          game.players.map((player) => (
            <PlayerRow key={player.id} player={player} />
          ))
        ) : (
          <span className="text-small text-default-500 leading-5">
            No players
          </span>
        )}
      </div>
    </div>
  );
}

export function EventGamesList({
  games,
  emptyContent,
  ariaLabel,
  mode,
  onSelect,
  hasMore,
  isLoadingMore,
  onLoadMore,
}: EventGamesListProps) {
  const items = games.map((game, index) => ({
    key: game.id,
    game,
    isLastInList: index === games.length - 1,
  }));

  return (
    <div className="flex flex-col w-full gap-3">
      <ListboxWrapper>
        <Listbox
          emptyContent={<p className="text-center">{emptyContent}</p>}
          aria-label={ariaLabel}
          variant="flat"
          color="primary"
          items={items}
          onAction={(key) => onSelect(Number(key))}
        >
          {(item) => (
            <ListboxItem
              key={item.key}
              textValue={gameTextValue(item.game)}
              showDivider={!item.isLastInList}
              endContent={<ItemCounter />}
            >
              <GameListItemContent game={item.game} mode={mode} />
            </ListboxItem>
          )}
        </Listbox>
      </ListboxWrapper>
      {hasMore && (
        <Button
          variant="flat"
          color="primary"
          isLoading={isLoadingMore}
          onPress={onLoadMore}
          className="self-center"
        >
          Load more
        </Button>
      )}
    </div>
  );
}

export function useEventGamesNavigate(eventId: number) {
  const navigate = useNavigate();
  return (gameId: number) => {
    navigate(`/event/${eventId}/game/${gameId}`);
  };
}
