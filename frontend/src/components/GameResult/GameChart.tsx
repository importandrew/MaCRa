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
// The x axis counts rounds: round r ends at x = r. A round can have several
// entries (a penalty or false mahjong besides its outcome); they share the
// stretch from r - 1 to r evenly, so rounds stay evenly spaced and the last
// entry of round r is the one at x = r.
interface ChartPoint {
  x: number;
  round: number;
  outcome: Outcome | null;
  isLastInRound: boolean;
  scores: number[];
}

function transformGameDataForChart(input: PreparedGameRoundLogs): ChartPoint[] {
  const totalAfter = (round: PreparedRoundLog, index: number) =>
    round.initial_scores[index] + (round.score_changes?.[index] ?? 0);
  const entriesInRound = new Map<number, number>();
  for (const round of input.rounds) {
    entriesInRound.set(round.round, (entriesInRound.get(round.round) ?? 0) + 1);
  }

  const points: ChartPoint[] = [
    {
      x: 0,
      round: 0,
      outcome: null,
      isLastInRound: true,
      scores: input.players.map(() => 0),
    },
  ];
  let indexInRound = 0;
  input.rounds.forEach((round, index) => {
    const previous = input.rounds[index - 1];
    indexInRound = previous?.round === round.round ? indexInRound + 1 : 1;
    const count = entriesInRound.get(round.round) ?? 1;
    points.push({
      x: round.round - 1 + indexInRound / count,
      round: round.round,
      outcome: round.outcome,
      isLastInRound: indexInRound === count,
      scores: input.players.map((_, seat) => totalAfter(round, seat)),
    });
  });
  return points;
}

function describePoint(point: ChartPoint): string {
  if (point.round === 0) return 'Start';
  if (point.outcome === Outcome.Penalty) return `Round ${point.round}, penalty`;
  if (point.outcome === Outcome.FalseMahjong)
    return `Round ${point.round}, false mahjong`;
  return point.isLastInRound
    ? `After round ${point.round}`
    : `Round ${point.round}`;
}

// The labeled rounds: every round, or every 2nd, 5th or 10th one, whichever
// fits. recharts' own thinning drops labels from the end unevenly.
const Y_AXIS_WIDTH = 36;
const DIGIT_WIDTH = 7; // px per digit of a 12px label, with some slack
const LABEL_GAP = 6;
function pickTicks(lastRound: number, chartWidth: number): number[] {
  const roundWidth = (chartWidth - Y_AXIS_WIDTH - 16) / Math.max(lastRound, 1);
  const labelWidth = String(lastRound).length * DIGIT_WIDTH + LABEL_GAP;
  const every = [1, 2, 5, 10].find((n) => n * roundWidth >= labelWidth) ?? 20;
  const ticks = [];
  for (let round = 0; round <= lastRound; round += every) ticks.push(round);
  return ticks;
}

// Reports the x of the point under the finger or pointer (recharts tracks it
// for both mouse and touch), or null when there is none.
function ActivePointReporter({
  onChange,
}: {
  onChange: (x: number | null) => void;
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
  const [activeX, setActiveX] = useState<number | null>(null);
  const last = chartData[chartData.length - 1];
  const shown = chartData.find((point) => point.x === activeX) ?? last;
  const isFinal = shown === last;
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
                dataKey="x"
                type="number"
                domain={[0, last.round]}
                ticks={pickTicks(last.round, chartWidth)}
                interval={0}
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
              <ActivePointReporter onChange={setActiveX} />
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
