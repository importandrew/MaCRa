import { useTranslation } from 'react-i18next';
import { Card } from '@heroui/react';
import type { ReactNode } from 'react';
import type { PlayerStats } from '../../models/PlayerStats';
import { currentLanguage } from '../../i18n';

// PLACES: the places in the order the statistics list them.
const PLACES = [1, 2, 3, 4] as const;

function StatRow({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="flex items-baseline justify-between gap-3 py-1.5">
      <dt className="text-muted">{label}</dt>
      <dd className="text-end tabular-nums">{children}</dd>
    </div>
  );
}

// PlayerStatsView shows a player's statistics in an event in two cards: per
// game (places, TP, GP) and per hand (wins, deal-ins, draws, penalties).
// Shares are of the hands played.
export default function PlayerStatsView({ stats }: { stats: PlayerStats }) {
  const { t } = useTranslation();

  if (stats.games === 0) {
    return <p className="py-6 text-center text-muted">{t('rating.noGames')}</p>;
  }

  const number = (value: number, digits = 1) =>
    value.toLocaleString(currentLanguage(), { maximumFractionDigits: digits });
  const signed = (value: number) => (value > 0 ? '+' : '') + number(value);
  const share = (count: number) =>
    stats.hands > 0
      ? t('stats.countShare', {
          n: count,
          share: number((count * 100) / stats.hands),
        })
      : String(count);
  const wins = stats.self_draws + stats.deal_in_wins;

  return (
    <div className="flex w-full flex-col gap-4">
      <Card>
        <Card.Content>
          <h2 className="font-medium">{t('stats.gamesTitle')}</h2>
          <dl className="divide-y divide-separator">
            <StatRow label={t('stats.games')}>{stats.games}</StatRow>
            <StatRow label={t('stats.places')}>
              <span className="flex flex-wrap justify-end gap-x-3">
                {PLACES.map((place) => (
                  <span key={place}>
                    {t('stats.place', {
                      place,
                      n: stats.places[place - 1],
                    })}
                  </span>
                ))}
              </span>
            </StatRow>
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
          </dl>
        </Card.Content>
      </Card>

      <Card>
        <Card.Content>
          <h2 className="font-medium">{t('stats.handsTitle')}</h2>
          <dl className="divide-y divide-separator">
            <StatRow label={t('stats.hands')}>{stats.hands}</StatRow>
            <StatRow label={t('stats.wins')}>
              <div>{share(wins)}</div>
              <div className="text-xs text-muted">
                {t('stats.winsBy', {
                  self: stats.self_draws,
                  dealIn: stats.deal_in_wins,
                })}
              </div>
            </StatRow>
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
            <StatRow label={t('stats.penalties')}>{stats.penalties}</StatRow>
            <StatRow label={t('stats.falseMahjongs')}>
              {stats.false_mahjongs}
            </StatRow>
          </dl>
        </Card.Content>
      </Card>
    </div>
  );
}
