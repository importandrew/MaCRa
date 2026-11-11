import type { TFunction } from 'i18next';
import { useTranslation } from 'react-i18next';
import { ProgressCircle } from '@heroui/react';
import {
  EventGameListItem,
  EventGameListPlayer,
} from '../../models/EventGames';
import NavList from '../NavList';
import PendingButton from '../PendingButton';

type EventGamesListProps = {
  games: EventGameListItem[];
  emptyContent: string;
  ariaLabel: string;
  mode: 'live' | 'finished';
  gameHref: (gameId: number) => string;
  hasMore: boolean;
  isLoadingMore: boolean;
  onLoadMore: () => void;
};

const TOTAL_ROUNDS = 16;

function roundProgress(
  game: EventGameListItem,
  mode: 'live' | 'finished',
  t: TFunction
) {
  if (mode === 'finished' || game.status === 'ended' || !game.round) {
    return { value: 100, label: '✓', ariaLabel: t('games.ended') };
  }

  return {
    value: ((game.round.num - 1) / TOTAL_ROUNDS) * 100,
    label: game.round.kanji,
    ariaLabel: t('games.round', { round: game.round.kanji }),
  };
}

function gameTextValue(game: EventGameListItem, t: TFunction) {
  const players = game.players.map((p) => p.name).join(', ');
  return players
    ? t('games.gameWithPlayers', { id: game.id, players })
    : t('games.gameNumber', { id: game.id });
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
      {hasWind && <span className="text-muted shrink-0">{player.wind}</span>}
      <span className="truncate whitespace-nowrap min-w-0">{player.name}</span>
      <span className="tabular-nums text-right shrink-0 text-foreground">
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
  const { t } = useTranslation();
  const { value, label, ariaLabel } = roundProgress(game, mode, t);

  return (
    <div className="grid grid-cols-[auto_minmax(0,1fr)] gap-3 min-w-0 w-full items-center">
      {/* h-0: don't inflate row height; min-h-full + aspect-square: match players column */}
      <div className="h-0 min-h-full aspect-square">
        <div className="relative h-full w-full">
          <ProgressCircle
            aria-label={ariaLabel}
            value={value}
            className="h-full w-full"
          >
            <ProgressCircle.Track className="h-full w-full">
              <ProgressCircle.TrackCircle />
              <ProgressCircle.FillCircle />
            </ProgressCircle.Track>
          </ProgressCircle>
          <span className="absolute inset-0 flex items-center justify-center text-lg font-medium">
            {label}
          </span>
        </div>
      </div>
      <div className="flex flex-col min-w-0">
        {game.players.length > 0 ? (
          game.players.map((player) => (
            <PlayerRow key={player.id} player={player} />
          ))
        ) : (
          <span className="text-sm text-muted leading-5">
            {t('games.noPlayers')}
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
  gameHref,
  hasMore,
  isLoadingMore,
  onLoadMore,
}: EventGamesListProps) {
  const { t } = useTranslation();
  return (
    <div className="flex flex-col w-full gap-3">
      <NavList
        ariaLabel={ariaLabel}
        emptyContent={<p className="text-center">{emptyContent}</p>}
        items={games.map((game) => ({
          id: String(game.id),
          href: gameHref(game.id),
          textValue: gameTextValue(game, t),
          content: <GameListItemContent game={game} mode={mode} />,
        }))}
      />
      {hasMore && (
        <PendingButton
          variant="secondary"
          isPending={isLoadingMore}
          onPress={onLoadMore}
          className="self-center"
        >
          {t('games.loadMore')}
        </PendingButton>
      )}
    </div>
  );
}
