import { useTranslation } from 'react-i18next';
import { Table } from '@heroui/react';
import { GameResult } from '../../models/GameResult';
import { UserData } from '../../models/UserData';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';
import { getDateTimeFormatted } from '../../utils/date-tools';

// Table component for rendering game results
export default function GameResultTable({ data }: { data: GameResult }) {
  const { t } = useTranslation();
  const currentPlayerId = (
    JSON.parse(
      localStorage.getItem(AUTH_USER_DATA) || '{}'
    ) as Partial<UserData>
  ).id;

  return (
    <div className="flex flex-col gap-1">
      <div className="flex justify-between items-center">
        <span className="text-sm font-medium text-muted">
          {t('games.gameNumber', { id: data.id })}
        </span>
        <span className="text-sm text-muted">
          {getDateTimeFormatted(data.finished_at)}
        </span>
      </div>
      <Table>
        <Table.ScrollContainer>
          <Table.Content aria-label={t('result.table')}>
            <Table.Header>
              <Table.Column isRowHeader>{t('result.name')}</Table.Column>
              <Table.Column className="text-end">{t('rating.gp')}</Table.Column>
              <Table.Column className="text-end">{t('rating.tp')}</Table.Column>
            </Table.Header>
            <Table.Body items={data.player_results}>
              {(item) => (
                <Table.Row
                  // By id, not name: players can share a name, and React Aria
                  // throws on duplicate keys in a collection.
                  id={item.id}
                  // Highlight the signed-in player's row.
                  className={
                    item.id === currentPlayerId
                      ? 'bg-accent-soft font-medium'
                      : undefined
                  }
                >
                  <Table.Cell>{item.name}</Table.Cell>
                  <Table.Cell className="text-end">{item.gp}</Table.Cell>
                  <Table.Cell className="text-end">{item.tp}</Table.Cell>
                </Table.Row>
              )}
            </Table.Body>
          </Table.Content>
        </Table.ScrollContainer>
      </Table>
    </div>
  );
}
