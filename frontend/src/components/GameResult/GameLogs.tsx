import { useTranslation } from 'react-i18next';
import { Button, Checkbox, Popover, Table } from '@heroui/react';
import { useCallback, useRef, useState } from 'react';
import { popoverVariants } from '@heroui/styles';
import { Outcome } from '../../models/GameData';
import {
  GameRoundLogs,
  LogUser,
  PreparedGameRoundLogs,
  PreparedRoundLog,
  UndoneEntry,
} from '../../models/GameResult';
import { getTimeFormatted } from '../../utils/date-tools';
import GameChart from './GameChart';
import { SignIcon } from '../Icons/SignIcon';
import { WarningIcon } from '../Icons/WarningIcon';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';

// popoverSlots style the author popover like a HeroUI Popover: it opens from
// a table row, not from a Popover's own trigger, so it has no Popover root to
// get them from.
const popoverSlots = popoverVariants();

// LogRow is a row of the log table: an entry (index into the entries) or an
// undone entry (index into the undone ones).
type LogRow =
  | { kind: 'entry'; item: PreparedRoundLog; index: number }
  | { kind: 'undone'; entry: UndoneEntry; index: number };

// LogsState is a round log as its owner has it: a table screen has it in its
// table response; the result page loads it (and can show loading, an error
// and a retry).
export type LogsState = {
  data: GameRoundLogs | null;
  isLoading?: boolean;
  isError?: boolean;
  fetchData?: () => unknown;
};

// GameLogs shows a round log as a chart or a table. It loads nothing: the
// caller passes the log it has.
export default function GameLogs({
  isChartVisible = false,
  isLogVisible = false,
  logs,
}: {
  isChartVisible: boolean;
  isLogVisible: boolean;
  logs: LogsState;
}) {
  const { t } = useTranslation();
  const { data, isError = false, isLoading = false, fetchData } = logs;
  const tableRef = useRef<HTMLDivElement>(null);
  const rowRef = useRef<HTMLElement | null>(null);
  // openKey: the row whose author popover is open (log-<index> or
  // undone-<index>), or null.
  const [openKey, setOpenKey] = useState<string | null>(null);
  // hideUndone: the viewer hides the undone entries (they show by default).
  const [hideUndone, setHideUndone] = useState(false);

  function getFinalScores() {
    if (!preparedData || !preparedData.rounds.length) {
      return [0, 0, 0, 0];
    }

    const item = preparedData.rounds.slice(-1)[0];

    if (!item.initial_scores || !item.score_changes) return;

    const matrix: [number[], number[]] = [
      item.initial_scores,
      item.score_changes,
    ];
    const result = matrix.reduce(
      (r, a) => r.map((b, i) => a[i] + b) as [number, number, number, number]
    );

    return result;
  }

  function transformGameDataForTable(
    input: GameRoundLogs | null
  ): PreparedGameRoundLogs | null {
    if (!input) return null;

    const rounds = input.round_logs.map((round, index, arr) => {
      let initial_scores = [0, 0, 0, 0];

      if (round.round !== 1) {
        initial_scores = [
          arr[index - 1].east_score,
          arr[index - 1].south_score,
          arr[index - 1].west_score,
          arr[index - 1].north_score,
        ];
      }

      let score_changes = [0, 0, 0, 0];

      if (round.outcome !== Outcome.Draw && round.round === 1) {
        score_changes = [
          round.east_score,
          round.south_score,
          round.west_score,
          round.north_score,
        ];
      }

      if (round.outcome !== Outcome.Draw && round.round !== 1) {
        score_changes = [
          round.east_score - arr[index - 1].east_score,
          round.south_score - arr[index - 1].south_score,
          round.west_score - arr[index - 1].west_score,
          round.north_score - arr[index - 1].north_score,
        ];
      }

      return {
        round: round.round,
        outcome: round.outcome,
        initial_scores: initial_scores,
        result: round.value,
        score_changes: score_changes,
      };
    });

    return {
      players: input.players,
      rounds: rounds,
    } as PreparedGameRoundLogs;
  }

  const preparedData = transformGameDataForTable(data);

  const getColoredScore = (value: number | undefined) => {
    if (value === undefined) {
      return 'text-muted';
    }

    if (value === 0) {
      return 'text-muted';
    }

    if (value < 0) {
      return 'text-danger';
    }

    if (value > 0) {
      return 'text-success';
    }
  };

  const getSignedScore = (value: number | undefined) => {
    if (value === undefined) return 'n/a';

    return (value <= 0 ? '' : '+') + value;
  };

  const getResultScore = (item: PreparedRoundLog) => {
    if (item.outcome === null) return 'n/a';

    if (item.outcome === Outcome.Deal) {
      return <span className="text-danger">{item.result}</span>;
    }

    if (item.outcome === Outcome.Self) {
      return <span className="text-success">{item.result}</span>;
    }

    if (item.outcome === Outcome.Draw) {
      return <SignIcon sign="minus" className="justify-self-end" />;
    }

    if ([Outcome.Penalty, Outcome.FalseMahjong].includes(item.outcome)) {
      return <WarningIcon className="justify-self-end text-warning" />;
    }
  };

  const getLegend = () => {
    return (
      <legend>
        <div>
          <span className="text-success">#: {t('roundEntry.self')}</span>
        </div>
        <div>
          <span className="text-danger">#: {t('roundEntry.dealIn')}</span>
        </div>
        <div>
          <span>
            <SignIcon sign="minus" className="inline" />: {t('roundEntry.draw')}
          </span>
        </div>
        <div>
          <span className="text-warning">
            <WarningIcon className="inline" />: {t('roundEntry.penalty')}
          </span>
        </div>
      </legend>
    );
  };

  const finalScores = getFinalScores();

  const onRefresh = useCallback(() => {
    fetchData?.();
  }, [fetchData]);

  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-full my-auto justify-center items-center">
        <LoadingSpinner label={t('common.loading')} />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="flex w-full h-full justify-center items-center">
        {fetchData && (
          <Button type="button" onPress={onRefresh}>
            {t('common.refresh')}
          </Button>
        )}
      </div>
    );
  }

  if (!data || !preparedData) {
    return (
      <div className="flex flex-col h-full w-full justify-center items-center gap-3">
        <p className="mb-4 text-center text-3xl">{t('logs.noData')}</p>
        {fetchData && (
          <Button type="button" onPress={onRefresh}>
            {t('common.refresh')}
          </Button>
        )}
      </div>
    );
  }

  // logRows interleaves the entries with the undone ones: each undone entry
  // shows where it was undone, before the first entry entered after that.
  const undone = data.undone ?? [];
  const logRows: LogRow[] = [];
  let nextUndone = 0;
  preparedData.rounds.forEach((item, index) => {
    const enteredAt = data.round_logs[index]?.created_at ?? '';
    while (
      nextUndone < undone.length &&
      undone[nextUndone].undone_at < enteredAt
    ) {
      logRows.push({
        kind: 'undone',
        entry: undone[nextUndone],
        index: nextUndone,
      });
      nextUndone++;
    }
    logRows.push({ kind: 'entry', item, index });
  });
  for (; nextUndone < undone.length; nextUndone++) {
    logRows.push({
      kind: 'undone',
      entry: undone[nextUndone],
      index: nextUndone,
    });
  }

  // undoneItem shows an undone entry like any other: the scores before it
  // and the change it made.
  const undoneItem = (entry: UndoneEntry): PreparedRoundLog => ({
    round: entry.round,
    outcome: entry.outcome,
    initial_scores: entry.scores_before,
    result: entry.value,
    score_changes: entry.scores_after.map(
      (score, seat) => score - entry.scores_before[seat]
    ) as [number, number, number, number],
  });

  // A row of the log: an entry, or an undone one, struck through and dimmed.
  // A tap opens the popover with who entered it (and who undid it).
  const logRow = (item: PreparedRoundLog, key: string, isUndone: boolean) => (
    <Table.Row
      key={key}
      id={key}
      className={
        'cursor-pointer ' +
        (isUndone ? 'line-through opacity-50' : 'even:bg-default/50')
      }
    >
      <Table.Cell className="text-end">{item.round}</Table.Cell>
      <Table.Cell className="text-end">
        <div className="text-base">
          {getSignedScore(item.initial_scores[0])}
        </div>
        <div className={'text-xs ' + getColoredScore(item.score_changes?.[0])}>
          {getSignedScore(item.score_changes?.[0])}
        </div>
      </Table.Cell>
      <Table.Cell className="text-end">
        <div className="text-base">
          {getSignedScore(item.initial_scores[1])}
        </div>
        <div className={'text-xs ' + getColoredScore(item.score_changes?.[1])}>
          {getSignedScore(item.score_changes?.[1])}
        </div>
      </Table.Cell>
      <Table.Cell className="text-end">
        <div className="text-base">
          {getSignedScore(item.initial_scores[2])}
        </div>
        <div className={'text-xs ' + getColoredScore(item.score_changes?.[2])}>
          {getSignedScore(item.score_changes?.[2])}
        </div>
      </Table.Cell>
      <Table.Cell className="text-end">
        <div className="text-base">
          {getSignedScore(item.initial_scores[3])}
        </div>
        <div className={'text-xs ' + getColoredScore(item.score_changes?.[3])}>
          {getSignedScore(item.score_changes?.[3])}
        </div>
      </Table.Cell>
      <Table.Cell className="text-end">
        <div className="text-base">{getResultScore(item)}</div>
      </Table.Cell>
    </Table.Row>
  );

  const shownRows = hideUndone
    ? logRows.filter((row) => row.kind === 'entry')
    : logRows;

  // The row tapped last (log-<index> or undone-<index>): a popover under it
  // says who entered it, and for an undone entry who undid it.
  const openIndex = openKey == null ? -1 : Number(openKey.split('-')[1]);
  const openEntry = openKey?.startsWith('log-')
    ? data.round_logs[openIndex]
    : undefined;
  const openUndone = openKey?.startsWith('undone-')
    ? undone[openIndex]
    : undefined;
  const enteredText = (by: LogUser, at: string) =>
    by
      ? t('logs.enteredBy', { name: by.name, time: getTimeFormatted(at) })
      : t('logs.enteredAt', { time: getTimeFormatted(at) });

  return (
    <>
      {isChartVisible ? <GameChart preparedData={preparedData} /> : null}

      {isLogVisible ? (
        <div
          ref={tableRef}
          className={
            'overflow-hidden flex flex-col justify-center mb-3 w-full ' +
            '[&_table]:w-full [&_table]:table-fixed [&_:is(th,td)]:box-border [&_:is(th,td)]:px-1 ' +
            '[&_:is(th,td):first-child]:w-10 [&_:is(th,td):first-child]:min-w-10 ' +
            '[&_:is(th,td):last-child]:w-10 [&_:is(th,td):last-child]:min-w-10 ' +
            '[&_:is(th,td):nth-child(n+2):nth-child(-n+5)]:w-[calc((100%-5rem)/4)] [&_:is(th,td):nth-child(n+2):nth-child(-n+5)]:min-w-0'
          }
        >
          {undone.length > 0 && (
            <Checkbox
              isSelected={hideUndone}
              onChange={setHideUndone}
              className="mb-2 self-start"
            >
              <Checkbox.Content>
                <Checkbox.Control>
                  <Checkbox.Indicator />
                </Checkbox.Control>
                {t('logs.hideUndone', { count: undone.length })}
              </Checkbox.Content>
            </Checkbox>
          )}
          <Table className="w-full self-center z-10">
            <Table.ScrollContainer>
              <Table.Content
                aria-label={t('game.logs')}
                onRowAction={(key) => {
                  const id = String(key);
                  if (id === 'finalScores') return;
                  rowRef.current = tableRef.current?.querySelector(
                    `[data-key="${id}"]`
                  ) as HTMLElement | null;
                  setOpenKey(id);
                }}
              >
                <Table.Header>
                  <Table.Column
                    id="round"
                    isRowHeader
                    className="min-w-0 text-end"
                  >
                    #
                  </Table.Column>
                  {preparedData.players.map((player) => (
                    <Table.Column
                      key={player.id}
                      id={String(player.id)}
                      className="min-w-0 overflow-hidden text-end"
                    >
                      <div
                        className="overflow-hidden text-ellipsis"
                        title={`ID ${player.id}`}
                      >
                        {player.name}
                      </div>
                    </Table.Column>
                  ))}
                  <Table.Column id="result" className="min-w-0 text-end">
                    {/* A popover, not a tooltip: it opens on tap on phones. */}
                    <Popover>
                      <Button
                        size="sm"
                        variant="ghost"
                        isIconOnly
                        aria-label={t('logs.legend')}
                      >
                        Δ
                      </Button>
                      <Popover.Content className="max-w-64">
                        <Popover.Dialog>{getLegend()}</Popover.Dialog>
                      </Popover.Content>
                    </Popover>
                  </Table.Column>
                </Table.Header>
                <Table.Body>
                  {[
                    // Keyed by position: round, outcome and scores can repeat
                    // (two equal penalties in one round), and React Aria
                    // throws on duplicate keys in a collection.
                    ...shownRows.map((row) =>
                      row.kind === 'entry'
                        ? logRow(row.item, `log-${row.index}`, false)
                        : logRow(
                            undoneItem(row.entry),
                            `undone-${row.index}`,
                            true
                          )
                    ),
                    ...(finalScores?.length === 4
                      ? [
                          <Table.Row
                            key="finalScores"
                            id="finalScores"
                            className="bg-accent-soft"
                          >
                            <Table.Cell className="text-end text-base font-bold">
                              Σ
                            </Table.Cell>
                            <Table.Cell className="text-end text-base font-bold">
                              {finalScores[0]}
                            </Table.Cell>
                            <Table.Cell className="text-end text-base font-bold">
                              {finalScores[1]}
                            </Table.Cell>
                            <Table.Cell className="text-end text-base font-bold">
                              {finalScores[2]}
                            </Table.Cell>
                            <Table.Cell className="text-end text-base font-bold">
                              {finalScores[3]}
                            </Table.Cell>
                            <Table.Cell> </Table.Cell>
                          </Table.Row>,
                        ]
                      : []),
                  ]}
                </Table.Body>
              </Table.Content>
            </Table.ScrollContainer>
          </Table>
          <Popover.Content
            triggerRef={rowRef}
            isOpen={openEntry != null || openUndone != null}
            onOpenChange={(isOpen) => {
              if (!isOpen) setOpenKey(null);
            }}
            placement="bottom"
            className={popoverSlots.base({ className: 'max-w-64' })}
          >
            <Popover.Dialog
              aria-label={t('logs.entryDetails')}
              className={popoverSlots.dialog()}
            >
              {openEntry &&
                enteredText(openEntry.entered_by, openEntry.created_at)}
              {openUndone && (
                <>
                  <div>
                    {enteredText(openUndone.entered_by, openUndone.entered_at)}
                  </div>
                  <div>
                    {openUndone.undone_by
                      ? t('logs.undoneBy', {
                          name: openUndone.undone_by.name,
                          time: getTimeFormatted(openUndone.undone_at),
                        })
                      : t('logs.undoneAt', {
                          time: getTimeFormatted(openUndone.undone_at),
                        })}
                  </div>
                </>
              )}
            </Popover.Dialog>
          </Popover.Content>
        </div>
      ) : null}
    </>
  );
}
