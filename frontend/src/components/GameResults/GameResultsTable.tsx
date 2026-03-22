import {
  Table,
  TableBody,
  TableCell,
  TableColumn,
  TableHeader,
  TableRow,
} from '@heroui/react';
import { GameResult } from '../../models/GameResult';
import { AUTH_USER_DATA } from '../../providers/auth/authConstants';
import { getDateTimeFormatted } from '../../utils/date-tools';

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
          {getDateTimeFormatted(data.finished_at)}
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
