import {
  Brush,
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ReferenceLine,
  ResponsiveContainer,
  Tooltip as TooltipChart,
  XAxis,
  YAxis,
} from 'recharts';
import {
  PreparedGameRoundLogs,
  PreparedRoundLog,
} from '../../models/GameResult';

function transformGameDataForChart(input: PreparedGameRoundLogs | null) {
  if (!input) return [];

  const [player1, player2, player3, player4] = input.players.map(
    (player) => player.name
  );

  const getRoundResult = (round: PreparedRoundLog, index: number) => {
    return (
      round.initial_scores?.[index] +
      (round.score_changes ? round.score_changes?.[index] : 0)
    );
  };

  const result = input.rounds.map((round) => {
    return {
      round: round.round,
      [player1]: getRoundResult(round, 0),
      [player2]: getRoundResult(round, 1),
      [player3]: getRoundResult(round, 2),
      [player4]: getRoundResult(round, 3),
    };
  });

  return [
    {
      round: 0,
      [player1]: 0,
      [player2]: 0,
      [player3]: 0,
      [player4]: 0,
    },
    ...result,
  ];
}

export default function GameChart({
  preparedData,
}: {
  preparedData: PreparedGameRoundLogs;
}) {
  const chartData = transformGameDataForChart(preparedData);

  return (
    <ResponsiveContainer
      width="100%"
      height={400}
      minWidth={300}
      className="mb-3 text-base"
    >
      <LineChart data={chartData} width={600} height={400}>
        <CartesianGrid />
        <XAxis dataKey="round" />
        <YAxis />
        <ReferenceLine y={0} stroke="#71717A" strokeWidth={1} />
        <Legend />
        <Line
          dataKey={preparedData.players[0].name}
          activeDot={{
            r: 6,
          }}
          stroke="#006FEE"
        />
        <Line
          dataKey={preparedData.players[1].name}
          activeDot={{
            r: 6,
          }}
          stroke="#7828C8"
        />
        <Line
          dataKey={preparedData.players[2].name}
          activeDot={{
            r: 6,
          }}
          stroke="#17C964"
        />
        <Line
          dataKey={preparedData.players[3].name}
          activeDot={{
            r: 6,
          }}
          stroke="#F31260"
        />
        <TooltipChart cursor={{ stroke: '#71717A', strokeWidth: 1 }} />
        <Brush />
      </LineChart>
    </ResponsiveContainer>
  );
}
