import { Button, Popover, Table } from '@heroui/react';
import { useCallback } from 'react';
import { Outcome } from '../../models/GameData';
import {
  GameRoundLogs,
  PreparedGameRoundLogs,
  PreparedRoundLog,
} from '../../models/GameResult';
import { useFetchGameLogsData } from '../../providers/game/gameDataApi';
import GameChart from './GameChart';
import { MinusIcon } from '../Icons/MinusIcon';
import { WarningIcon } from '../Icons/WarningIcon';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';

export default function GameLogs({
  isChartVisible = false,
  isLogVisible = false,
  gameId,
  eventId,
}: {
  isChartVisible: boolean;
  isLogVisible: boolean;
  gameId: number | null;
  eventId?: number;
}) {
  const { data, isError, isLoading, fetchData } = useFetchGameLogsData(
    gameId,
    eventId
  );

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
      return <MinusIcon className="justify-self-end" />;
    }

    if ([Outcome.Penalty, Outcome.FalseMahjong].includes(item.outcome)) {
      return <WarningIcon className="justify-self-end text-warning" />;
    }
  };

  const getRowKey = (item: PreparedRoundLog) => {
    let key = item.round.toString();

    if (item.outcome) {
      key += item.outcome;
    }

    if (item.result) {
      key += item.result;
    }

    if (item.score_changes) {
      key += item.initial_scores.join('');
    }

    return key;
  };

  const getLegend = () => {
    return (
      <legend>
        <div>
          <span className="text-success">#: {Outcome.Self}</span>
        </div>
        <div>
          <span className="text-danger">#: {Outcome.Deal}</span>
        </div>
        <div>
          <span>
            <MinusIcon className="inline" />: {Outcome.Draw}
          </span>
        </div>
        <div>
          <span className="text-warning">
            <WarningIcon className="inline" />: {Outcome.Penalty}
          </span>
        </div>
      </legend>
    );
  };

  const finalScores = getFinalScores();

  const onRefresh = useCallback(() => {
    fetchData();
  }, [fetchData]);

  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-full my-auto justify-center items-center">
        <LoadingSpinner label="Loading..." />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="flex w-full h-full justify-center items-center">
        <Button type="button" onPress={onRefresh}>
          Refresh
        </Button>
      </div>
    );
  }

  if (!data || !preparedData) {
    return (
      <div className="flex flex-col h-full w-full justify-center items-center gap-3">
        <p className="mb-4 text-center text-3xl">No data available</p>
        <Button type="button" onPress={onRefresh}>
          Refresh
        </Button>
      </div>
    );
  }

  return (
    <>
      {isChartVisible ? <GameChart preparedData={preparedData} /> : null}

      {isLogVisible ? (
        <div
          className={
            'overflow-hidden flex flex-col justify-center mb-3 w-full ' +
            '[&_table]:w-full [&_table]:table-fixed [&_:is(th,td)]:box-border [&_:is(th,td)]:px-1 ' +
            '[&_:is(th,td):first-child]:w-10 [&_:is(th,td):first-child]:min-w-10 ' +
            '[&_:is(th,td):last-child]:w-10 [&_:is(th,td):last-child]:min-w-10 ' +
            '[&_:is(th,td):nth-child(n+2):nth-child(-n+5)]:w-[calc((100%-5rem)/4)] [&_:is(th,td):nth-child(n+2):nth-child(-n+5)]:min-w-0'
          }
        >
          <Table className="w-full self-center z-10">
            <Table.ScrollContainer>
              <Table.Content aria-label="Game logs">
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
                        aria-label="Legend"
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
                    ...preparedData.rounds.map((item) => (
                      <Table.Row
                        key={getRowKey(item)}
                        id={getRowKey(item)}
                        className="even:bg-default/50"
                      >
                        <Table.Cell className="text-end">
                          {item.round}
                        </Table.Cell>
                        <Table.Cell className="text-end">
                          <div className="text-base">
                            {getSignedScore(item.initial_scores[0])}
                          </div>
                          <div
                            className={
                              'text-xs ' +
                              getColoredScore(item.score_changes?.[0])
                            }
                          >
                            {getSignedScore(item.score_changes?.[0])}
                          </div>
                        </Table.Cell>
                        <Table.Cell className="text-end">
                          <div className="text-base">
                            {getSignedScore(item.initial_scores[1])}
                          </div>
                          <div
                            className={
                              'text-xs ' +
                              getColoredScore(item.score_changes?.[1])
                            }
                          >
                            {getSignedScore(item.score_changes?.[1])}
                          </div>
                        </Table.Cell>
                        <Table.Cell className="text-end">
                          <div className="text-base">
                            {getSignedScore(item.initial_scores[2])}
                          </div>
                          <div
                            className={
                              'text-xs ' +
                              getColoredScore(item.score_changes?.[2])
                            }
                          >
                            {getSignedScore(item.score_changes?.[2])}
                          </div>
                        </Table.Cell>
                        <Table.Cell className="text-end">
                          <div className="text-base">
                            {getSignedScore(item.initial_scores[3])}
                          </div>
                          <div
                            className={
                              'text-xs ' +
                              getColoredScore(item.score_changes?.[3])
                            }
                          >
                            {getSignedScore(item.score_changes?.[3])}
                          </div>
                        </Table.Cell>
                        <Table.Cell className="text-end">
                          <div className="text-base">
                            {getResultScore(item)}
                          </div>
                        </Table.Cell>
                      </Table.Row>
                    )),
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
        </div>
      ) : null}
    </>
  );
}
