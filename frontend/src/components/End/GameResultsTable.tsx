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

// Table component for rendering game results
export default function GameResultsTable({ data }: { data: GameResult }) {
  const currentPlayer = JSON.parse(
    localStorage.getItem(AUTH_USER_DATA) || '{}'
  ).name;

  return (
    <Table
      isStriped
      defaultSelectedKeys={[currentPlayer]}
      color="primary"
      selectionMode="single"
      aria-label="Game results table"
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
  );
}
