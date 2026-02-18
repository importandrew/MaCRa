import {
  Table,
  TableHeader,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
} from '@heroui/react';
import { GameResult } from '../../models/GameResult';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';

function formatFinishedAt(iso: string): string {
  if (!iso) return '—';

  const d = new Date(iso);
  const time = d.toLocaleTimeString(undefined, {
    hour: '2-digit',
    minute: '2-digit',
  });
  const date = d.toLocaleDateString(undefined, { dateStyle: 'short' });

  return `${time} ${date}`;
}

// Table component for rendering game results
export default function GameResultsTable({ data }: { data: GameResult }) {
  const currentPlayer = JSON.parse(
    localStorage.getItem(AUTH_USER_DATA) || '{}'
  ).name;

  return (
    <div className="flex flex-col gap-1">
      <div className="flex justify-between items-center">
        <span className="text-sm font-medium text-default-600">
          Game #{data.id}
        </span>
        <span className="text-sm text-default-500">
          {formatFinishedAt(data.finished_at)}
        </span>
      </div>
      <Table
        isStriped
        removeWrapper
        defaultSelectedKeys={[currentPlayer]}
        color="primary"
        selectionMode="single"
        aria-label="Game results table"
        className="z-10"
      >
        <TableHeader>
          <TableColumn align="start">NAME</TableColumn>
          <TableColumn align="end">GP</TableColumn>
          <TableColumn align="end">TP</TableColumn>
        </TableHeader>
        <TableBody items={data.player_results}>
          {(item) => (
            <TableRow key={item.name}>
              <TableCell>{item.name}</TableCell>
              <TableCell>{item.gp}</TableCell>
              <TableCell>{item.tp}</TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
