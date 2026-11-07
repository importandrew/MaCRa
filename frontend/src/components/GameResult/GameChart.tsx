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
import { Outcome } from '../../models/GameData';
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

// A chart point: one log entry and each player's total after it, by seat.
// A round can have several entries (a penalty or false mahjong besides its
// outcome); each gets its own point, so `step` (the entry's position) is the
// x value, not the round number, which repeats.
interface ChartPoint {
  step: number;
  round: number;
  outcome: Outcome | null;
  scores: number[];
}

// Penalties and false mahjong are extra entries inside a round.
const isExtraEntry = (point: ChartPoint) =>
  point.outcome === Outcome.Penalty || point.outcome === Outcome.FalseMahjong;

function transformGameDataForChart(input: PreparedGameRoundLogs): ChartPoint[] {
  const totalAfter = (round: PreparedRoundLog, index: number) =>
    round.initial_scores[index] + (round.score_changes?.[index] ?? 0);

  return [
    { step: 0, round: 0, outcome: null, scores: input.players.map(() => 0) },
    ...input.rounds.map((round, index) => ({
      step: index + 1,
      round: round.round,
      outcome: round.outcome,
      scores: input.players.map((_, seat) => totalAfter(round, seat)),
    })),
  ];
}

function describePoint(point: ChartPoint): string {
  if (point.step === 0) return 'Start';
  if (point.outcome === Outcome.Penalty) return `Round ${point.round}, penalty`;
  if (point.outcome === Outcome.FalseMahjong)
    return `Round ${point.round}, false mahjong`;
  return `After round ${point.round}`;
}

// The x-axis label: the round number, or a dot for an extra entry.
const tickLabel = (point: ChartPoint | undefined) =>
  !point ? '' : isExtraEntry(point) ? '•' : String(point.round);

// Which points get an x-axis label. recharts thins crowded labels from the end
// on its own, which dropped e.g. the second-to-last round and gave uneven runs
// like "1 3 • 6". Instead, label every round, or every 2nd, 5th or 10th one,
// whichever fits; dots for extra entries only while every round fits.
const Y_AXIS_WIDTH = 36;
const DIGIT_WIDTH = 7; // px per digit of a 12px label, with some slack
const LABEL_GAP = 6;
function pickTicks(points: ChartPoint[], chartWidth: number): number[] {
  const plotWidth = chartWidth - Y_AXIS_WIDTH - 16;
  const stepWidth = plotWidth / Math.max(points.length - 1, 1);
  const lastRound = points[points.length - 1].round;
  const labelWidth = String(lastRound).length * DIGIT_WIDTH + LABEL_GAP;
  // Rounds between labels; the distance in px also counts extra entries.
  const roundsPerStep = (points.length - 1) / Math.max(lastRound, 1);
  const every =
    [1, 2, 5, 10].find((n) => (n * stepWidth) / roundsPerStep >= labelWidth) ??
    20;
  return points
    .filter((point) =>
      isExtraEntry(point)
        ? every === 1 && stepWidth >= labelWidth
        : point.round % every === 0
    )
    .map((point) => point.step);
}

// Reports the step under the finger or pointer (recharts tracks it for both
// mouse and touch), or null when none is.
function ActiveStepReporter({
  onChange,
}: {
  onChange: (step: number | null) => void;
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
  // The entry under the finger or pointer; the final scores otherwise.
  const [activeStep, setActiveStep] = useState<number | null>(null);
  const shown =
    chartData.find((point) => point.step === activeStep) ??
    chartData[chartData.length - 1];
  const isFinal = shown === chartData[chartData.length - 1];
  const [chartWidth, setChartWidth] = useState(0);

  return (
    <Card className="mb-3">
      <Card.Header className="flex-row items-baseline justify-between gap-3">
        <Card.Title>Score by round</Card.Title>
        <span className="text-sm text-muted">
          {isFinal ? 'Final' : describePoint(shown)}
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
        {/* recharts makes the chart focusable for keyboard users (arrow keys move
            between rounds). Show the ring for keyboard focus only, as HeroUI
            does; a click or tap focuses it too, which drew a box around it. */}
        <div className="h-60 sm:h-72 -mx-2 select-none touch-pan-y [&_:focus]:outline-none [&_.recharts-surface:focus-visible]:outline-2 [&_.recharts-surface:focus-visible]:outline-solid [&_.recharts-surface:focus-visible]:outline-(--focus) [&_.recharts-surface:focus-visible]:rounded-md">
          <ResponsiveContainer
            width="100%"
            height="100%"
            onResize={(width) => setChartWidth(width)}
          >
            <LineChart
              data={chartData}
              margin={{ top: 8, right: 8, bottom: 0, left: 0 }}
            >
              <CartesianGrid vertical={false} stroke="var(--separator)" />
              <XAxis
                dataKey="step"
                ticks={pickTicks(chartData, chartWidth)}
                interval={0}
                tickFormatter={(step: number) => tickLabel(chartData[step])}
                axisLine={false}
                tickLine={false}
                tick={{ fill: 'var(--muted)', fontSize: 12 }}
              />
              <YAxis
                axisLine={false}
                tickLine={false}
                tick={{ fill: 'var(--muted)', fontSize: 12 }}
                width={Y_AXIS_WIDTH}
              />
              <ReferenceLine y={0} stroke="var(--muted)" />
              {/* Only the crosshair and active dots: values are in the readout. */}
              <Tooltip
                content={() => null}
                cursor={{ stroke: 'var(--muted)', strokeDasharray: '3 3' }}
              />
              <ActiveStepReporter onChange={setActiveStep} />
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
