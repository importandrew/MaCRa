import { useTranslation } from 'react-i18next';
import { Card } from '@heroui/react';
import { type ReactNode, useId, useState } from 'react';
import type { PlayerStats } from '../../models/PlayerStats';
import { currentLanguage } from '../../i18n';
import { ChevronRightIcon } from '../Icons/ChevronRightIcon';

// PLACES: the places in the order the statistics list them.
const PLACES = [1, 2, 3, 4] as const;

const ROW_CLASS = 'flex items-baseline justify-between gap-3 py-1.5';

function StatRow({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className={ROW_CLASS}>
      <span className="text-muted">{label}</span>
      <span className="text-end tabular-nums">{children}</span>
    </div>
  );
}

// StatGroup is a row with a total that a tap expands into its parts (e.g.
// wins into self-draws and deal-ins), indented under it.
function StatGroup({
  label,
  value,
  children,
}: {
  label: string;
  value: ReactNode;
  children: ReactNode;
}) {
  const [isExpanded, setIsExpanded] = useState(false);
  const partsId = useId();
  return (
    <div>
      <button
        type="button"
        aria-expanded={isExpanded}
        aria-controls={partsId}
        onClick={() => setIsExpanded(!isExpanded)}
        className={`${ROW_CLASS} w-full cursor-pointer text-start`}
      >
        <span className="flex items-center gap-1 text-muted">
          {label}
          <ChevronRightIcon
            className={`shrink-0 transition-transform ${isExpanded ? 'rotate-90' : ''}`}
          />
        </span>
        <span className="text-end tabular-nums">{value}</span>
      </button>
      {isExpanded && (
        <div id={partsId} className="pb-1.5 ps-5 text-sm">
          {children}
        </div>
      )}
    </div>
  );
}

// PlayerStatsView shows a player's statistics in an event in two cards: per
// game (places, TP, GP) and per hand (wins, deal-ins, draws, penalties).
// Places, wins and penalties expand into their parts. Shares of places are of
// the games, the others of the hands played.
export default function PlayerStatsView({ stats }: { stats: PlayerStats }) {
  const { t } = useTranslation();

  if (stats.games === 0) {
    return <p className="py-6 text-center text-muted">{t('rating.noGames')}</p>;
  }

  const number = (value: number, digits = 1) =>
    value.toLocaleString(currentLanguage(), { maximumFractionDigits: digits });
  const signed = (value: number) => (value > 0 ? '+' : '') + number(value);
  const shareOf = (count: number, total: number) =>
    total > 0
      ? t('stats.countShare', {
          n: count,
          share: number((count * 100) / total),
        })
      : String(count);
  const share = (count: number) => shareOf(count, stats.hands);
  const wins = stats.self_draws + stats.deal_in_wins;
  const averagePlace =
    PLACES.reduce((sum, place) => sum + place * stats.places[place - 1], 0) /
    stats.games;

  return (
    <div className="flex w-full flex-col gap-4">
      <Card>
        <Card.Content>
          <h2 className="font-medium">{t('stats.gamesTitle')}</h2>
          <div className="divide-y divide-separator">
            <StatRow label={t('stats.games')}>{stats.games}</StatRow>
            <StatGroup
              label={t('stats.places')}
              value={t('stats.averagePlace', { place: number(averagePlace) })}
            >
              {PLACES.map((place) => (
                <StatRow key={place} label={t('stats.place', { place })}>
                  {shareOf(stats.places[place - 1], stats.games)}
                </StatRow>
              ))}
            </StatGroup>
            <StatRow label={t('stats.avgTP')}>
              {number(stats.avg_tp, 2)}
            </StatRow>
            <StatRow label={t('stats.avgGP')}>{signed(stats.avg_gp)}</StatRow>
            <StatRow label={t('stats.bestGame')}>
              {stats.best_gp == null ? '—' : signed(stats.best_gp)}
            </StatRow>
            <StatRow label={t('stats.worstGame')}>
              {stats.worst_gp == null ? '—' : signed(stats.worst_gp)}
            </StatRow>
          </div>
        </Card.Content>
      </Card>

      <Card>
        <Card.Content>
          <h2 className="font-medium">{t('stats.handsTitle')}</h2>
          <div className="divide-y divide-separator">
            <StatRow label={t('stats.hands')}>{stats.hands}</StatRow>
            <StatGroup label={t('stats.wins')} value={share(wins)}>
              <StatRow label={t('stats.selfDraws')}>
                {share(stats.self_draws)}
              </StatRow>
              <StatRow label={t('stats.dealInWins')}>
                {share(stats.deal_in_wins)}
              </StatRow>
            </StatGroup>
            <StatRow label={t('stats.avgWin')}>
              {wins > 0 ? number(stats.avg_win) : '—'}
            </StatRow>
            <StatRow label={t('stats.maxWin')}>
              {wins > 0 ? stats.max_win : '—'}
            </StatRow>
            <StatRow label={t('stats.dealIns')}>
              {share(stats.deal_ins)}
            </StatRow>
            <StatRow label={t('stats.draws')}>{share(stats.draws)}</StatRow>
            <StatGroup
              label={t('stats.penalties')}
              value={stats.penalties + stats.false_mahjongs}
            >
              <StatRow label={t('stats.penaltiesOnly')}>
                {stats.penalties}
              </StatRow>
              <StatRow label={t('stats.falseMahjongs')}>
                {stats.false_mahjongs}
              </StatRow>
            </StatGroup>
          </div>
        </Card.Content>
      </Card>
    </div>
  );
}
