import { Card } from '@heroui/react';
import { useEffect, useState } from 'react';
import {
  CartesianGrid,
  Line,
  LineChart,
  ReferenceLine,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
  useActiveTooltipLabel,
} from 'recharts';
import {
  PreparedGameRoundLogs,
  PreparedRoundLog,
} from '../../models/GameResult';

// One color per player, in seat order, from the --chart-* tokens in index.css.
const SERIES_COLORS = [
  'var(--chart-1)',
  'var(--chart-2)',
  'var(--chart-3)',
  'var(--chart-4)',
];

// A chart point: the round number and each player's total after it, by seat.
interface ChartPoint {
  round: number;
  scores: number[];
}

function transformGameDataForChart(input: PreparedGameRoundLogs): ChartPoint[] {
  const totalAfter = (round: PreparedRoundLog, index: number) =>
    round.initial_scores[index] + (round.score_changes?.[index] ?? 0);

  return [
    { round: 0, scores: input.players.map(() => 0) },
    ...input.rounds.map((round) => ({
      round: round.round,
      scores: input.players.map((_, index) => totalAfter(round, index)),
    })),
  ];
}

// Reports the round under the finger or pointer (recharts tracks it for both
// mouse and touch), or null when none is.
function ActiveRoundReporter({
  onChange,
}: {
  onChange: (round: number | null) => void;
}) {
  const label = useActiveTooltipLabel();
  useEffect(() => {
    onChange(label === undefined ? null : Number(label));
  }, [label, onChange]);
  return null;
}

export default function GameChart({
  preparedData,
}: {
  preparedData: PreparedGameRoundLogs;
}) {
  const chartData = transformGameDataForChart(preparedData);
  // The round under the finger or pointer; the final scores otherwise.
  const [activeRound, setActiveRound] = useState<number | null>(null);
  const shown =
    chartData.find((point) => point.round === activeRound) ??
    chartData[chartData.length - 1];
  const isFinal = shown === chartData[chartData.length - 1];

  return (
    <Card className="mb-3">
      <Card.Header className="flex-row items-baseline justify-between gap-3">
        <Card.Title>Score by round</Card.Title>
        <span className="text-sm text-muted">
          {isFinal
            ? 'Final'
            : shown.round === 0
              ? 'Start'
              : `After round ${shown.round}`}
        </span>
      </Card.Header>
      <Card.Content className="gap-3">
        {/* Legend and readout in one: each player's total at the shown round. */}
        <ul className="grid grid-cols-2 gap-x-4 gap-y-1 text-sm">
          {preparedData.players.map((player, index) => (
            <li key={player.id} className="flex items-center gap-2 min-w-0">
              <span
                aria-hidden
                className="h-0.5 w-3 shrink-0 rounded-full"
                style={{ backgroundColor: SERIES_COLORS[index] }}
              />
              <span className="truncate">{player.name}</span>
              <span className="ml-auto font-medium tabular-nums">
                {shown.scores[index]}
              </span>
            </li>
          ))}
        </ul>
        <div className="h-60 sm:h-72 -mx-2 select-none touch-pan-y">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart
              data={chartData}
              margin={{ top: 8, right: 8, bottom: 0, left: 0 }}
            >
              <CartesianGrid vertical={false} stroke="var(--separator)" />
              <XAxis
                dataKey="round"
                axisLine={false}
                tickLine={false}
                tick={{ fill: 'var(--muted)', fontSize: 12 }}
                minTickGap={12}
              />
              <YAxis
                axisLine={false}
                tickLine={false}
                tick={{ fill: 'var(--muted)', fontSize: 12 }}
                width={36}
              />
              <ReferenceLine y={0} stroke="var(--muted)" />
              {/* Only the crosshair and active dots: values are in the readout. */}
              <Tooltip
                content={() => null}
                cursor={{ stroke: 'var(--muted)', strokeDasharray: '3 3' }}
              />
              <ActiveRoundReporter onChange={setActiveRound} />
              {preparedData.players.map((player, index) => (
                <Line
                  key={player.id}
                  name={player.name}
                  dataKey={(point: ChartPoint) => point.scores[index]}
                  stroke={SERIES_COLORS[index]}
                  strokeWidth={2}
                  dot={false}
                  activeDot={{
                    r: 4,
                    stroke: 'var(--surface)',
                    strokeWidth: 2,
                  }}
                  isAnimationActive={false}
                />
              ))}
            </LineChart>
          </ResponsiveContainer>
        </div>
      </Card.Content>
    </Card>
  );
}
