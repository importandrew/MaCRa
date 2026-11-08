import { useEffect, useRef } from 'react';
import {
  GameLogsState,
  useFetchGameLogsData,
} from '../../providers/game/gameDataApi';

// useTableLogs loads the round log of the game a table screen shows, together
// with the table: again each time the table changes (load, Refresh, an entry,
// an undo), not when the menu opens. So the log in the table menu shows the
// same state as the table, and an undo sends the size of both.
//
// The table and the log are two requests. If someone enters or undoes an
// entry between them, the log's size differs from the table's; the table
// then reloads (reloadTable), and the log with it. Once per pair of sizes, so
// a mismatch that doesn't go away can't reload forever.
export function useTableLogs(
  table: { id: number; log_size: number } | null,
  reloadTable: () => unknown,
  eventId?: number
): GameLogsState {
  const logs = useFetchGameLogsData(table?.id ?? null, eventId, table);
  const retriedRef = useRef<string | null>(null);

  useEffect(() => {
    if (!table || !logs.data || logs.data.id !== table.id) return;
    const sizes = `${table.log_size}/${logs.data.round_logs.length}`;
    if (table.log_size === logs.data.round_logs.length) return;
    if (retriedRef.current === sizes) return;
    retriedRef.current = sizes;
    void reloadTable();
    // Compare when a log arrives, with the table shown at that moment.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [logs.data]);

  return logs;
}
