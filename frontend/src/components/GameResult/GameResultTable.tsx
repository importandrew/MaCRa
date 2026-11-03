import { Table } from '@heroui/react';
import { GameResult } from '../../models/GameResult';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';
import { getDateTimeFormatted } from '../../utils/date-tools';

// Table component for rendering game results
export default function GameResultTable({ data }: { data: GameResult }) {
  const currentPlayer = JSON.parse(
    localStorage.getItem(AUTH_USER_DATA) || '{}'
  ).name;

  return (
    <div className="flex flex-col gap-1">
      <div className="flex justify-between items-center">
        <span className="text-sm font-medium text-muted">Game #{data.id}</span>
        <span className="text-sm text-muted">
          {getDateTimeFormatted(data.finished_at)}
        </span>
      </div>
      <Table>
        <Table.ScrollContainer>
          <Table.Content aria-label="Game results table">
            <Table.Header>
              <Table.Column isRowHeader>NAME</Table.Column>
              <Table.Column className="text-end">GP</Table.Column>
              <Table.Column className="text-end">TP</Table.Column>
            </Table.Header>
            <Table.Body items={data.player_results}>
              {(item) => (
                <Table.Row
                  id={item.name}
                  // Highlight the signed-in player's row.
                  className={
                    item.name === currentPlayer
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
