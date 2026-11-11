import type { TFunction } from 'i18next';
import { useTranslation } from 'react-i18next';
import { Card } from '@heroui/react';
import { KeyboardEvent, useCallback, useEffect, useRef, useState } from 'react';
import {
  CartesianGrid,
  Line,
  LineChart,
  ReferenceDot,
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

const ACTIVE_DOT = { r: 4, stroke: 'var(--surface)', strokeWidth: 2 };

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
  // For a round with several entries of one kind (penalties): which one this
  // is and how many there are.
  sameKindNumber: number;
  sameKindCount: number;
  scores: number[];
}

function transformGameDataForChart(input: PreparedGameRoundLogs): ChartPoint[] {
  const totalAfter = (round: PreparedRoundLog, index: number) =>
    round.initial_scores[index] + (round.score_changes?.[index] ?? 0);
  const entriesInRound = new Map<number, number>();
  const kindKey = (round: PreparedRoundLog) =>
    `${round.round}:${round.outcome}`;
  const kindCount = new Map<string, number>();
  for (const round of input.rounds) {
    entriesInRound.set(round.round, (entriesInRound.get(round.round) ?? 0) + 1);
    kindCount.set(kindKey(round), (kindCount.get(kindKey(round)) ?? 0) + 1);
  }
  const kindSeen = new Map<string, number>();

  const points: ChartPoint[] = [
    {
      x: 0,
      round: 0,
      outcome: null,
      isLastInRound: true,
      sameKindNumber: 1,
      sameKindCount: 1,
      scores: input.players.map(() => 0),
    },
  ];
  let indexInRound = 0;
  input.rounds.forEach((round, index) => {
    const previous = input.rounds[index - 1];
    indexInRound = previous?.round === round.round ? indexInRound + 1 : 1;
    const count = entriesInRound.get(round.round) ?? 1;
    kindSeen.set(kindKey(round), (kindSeen.get(kindKey(round)) ?? 0) + 1);
    points.push({
      x: round.round - 1 + indexInRound / count,
      round: round.round,
      outcome: round.outcome,
      isLastInRound: indexInRound === count,
      sameKindNumber: kindSeen.get(kindKey(round)) ?? 1,
      sameKindCount: kindCount.get(kindKey(round)) ?? 1,
      scores: input.players.map((_, seat) => totalAfter(round, seat)),
    });
  });
  return points;
}

function describePoint(point: ChartPoint, t: TFunction): string {
  if (point.round === 0) return t('chart.start');
  const nth =
    point.sameKindCount > 1
      ? t('chart.nth', { n: point.sameKindNumber, total: point.sameKindCount })
      : '';
  if (point.outcome === Outcome.Penalty)
    return t('chart.penalty', { round: point.round, nth });
  if (point.outcome === Outcome.FalseMahjong)
    return t('chart.falseMahjong', { round: point.round, nth });
  return point.isLastInRound
    ? t('chart.afterRound', { round: point.round })
    : t('chart.round', { round: point.round });
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
  const { t } = useTranslation();
  const chartData = transformGameDataForChart(preparedData);
  const last = chartData[chartData.length - 1];
  const [chartWidth, setChartWidth] = useState(0);

  // Two ways to pick a point: hovering (mouse, finger), and a selected point
  // that stays: set by a click or tap, moved by the arrow keys. A key press
  // shows the selected point even while the pointer rests on the chart;
  // moving the pointer to another point shows that one again.
  const [hoverX, setHoverX] = useState<number | null>(null);
  const [selectedX, setSelectedX] = useState<number | null>(null);
  const [preferSelected, setPreferSelected] = useState(false);
  const onHover = useCallback((x: number | null) => {
    setHoverX(x);
    setPreferSelected(false);
  }, []);
  const shownX = preferSelected ? selectedX : (hoverX ?? selectedX);
  const shown = chartData.find((point) => point.x === shownX) ?? last;
  const isFinal = shown === last;
  // Our own marker for the selected point; recharts draws the hover one.
  const marked =
    shownX !== null &&
    shownX === selectedX &&
    (preferSelected || hoverX === null)
      ? shown
      : null;

  const chartRef = useRef<HTMLDivElement>(null);
  const onClick = () => {
    if (hoverX !== null) setSelectedX(hoverX);
    // So the arrow keys work right after a click (without a focus ring).
    chartRef.current?.focus({ preventScroll: true });
  };
  const onKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
    const index = chartData.indexOf(shown);
    const next =
      event.key === 'ArrowLeft'
        ? Math.max(index - 1, 0)
        : event.key === 'ArrowRight'
          ? Math.min(index + 1, chartData.length - 1)
          : event.key === 'Home'
            ? 0
            : event.key === 'End'
              ? chartData.length - 1
              : null;
    if (next === null) return;
    event.preventDefault();
    setSelectedX(chartData[next].x);
    setPreferSelected(true);
  };

  return (
    <Card className="mb-3">
      <Card.Header className="flex-row items-baseline justify-between gap-3">
        <Card.Title>{t('chart.title')}</Card.Title>
        <span className="text-sm text-muted">
          {isFinal ? t('chart.final') : describePoint(shown, t)}
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
        {/* Focusable for the arrow keys; the ring shows for keyboard focus only,
            as in HeroUI. recharts' own keyboard layer is off: it only worked
            when the SVG itself had focus, and lost to the hovering pointer. */}
        <div
          ref={chartRef}
          tabIndex={0}
          // A slider over the chart's points: the ARIA pattern for picking one
          // of an ordered set with the arrow keys, Home and End.
          role="slider"
          aria-label={t('chart.title')}
          aria-valuemin={0}
          aria-valuemax={chartData.length - 1}
          aria-valuenow={chartData.indexOf(shown)}
          aria-valuetext={`${isFinal ? t('chart.final') : describePoint(shown, t)}: ${preparedData.players.map((player, index) => `${player.name} ${shown.scores[index]}`).join(', ')}`}
          onClick={onClick}
          onKeyDown={onKeyDown}
          className="h-60 sm:h-72 -mx-2 rounded-md select-none touch-pan-y outline-none focus-visible:outline-2 focus-visible:outline-solid focus-visible:outline-(--focus) [&_:focus]:outline-none"
        >
          <ResponsiveContainer
            width="100%"
            height="100%"
            onResize={(width) => setChartWidth(width)}
          >
            <LineChart
              data={chartData}
              margin={{ top: 8, right: 8, bottom: 0, left: 0 }}
              accessibilityLayer={false}
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
              {/* Only the crosshair and active dots: values are in the readout.
                  Hidden while a key press shows the selected point. */}
              <Tooltip
                content={() => null}
                cursor={
                  preferSelected
                    ? false
                    : { stroke: 'var(--muted)', strokeDasharray: '3 3' }
                }
              />
              <ActivePointReporter onChange={onHover} />
              {preparedData.players.map((player, index) => (
                <Line
                  key={player.id}
                  name={player.name}
                  dataKey={(point: ChartPoint) => point.scores[index]}
                  stroke={SERIES_COLORS[index]}
                  strokeWidth={2}
                  dot={false}
                  activeDot={preferSelected ? false : ACTIVE_DOT}
                  isAnimationActive={false}
                />
              ))}
              {marked && (
                <ReferenceLine
                  x={marked.x}
                  stroke="var(--muted)"
                  strokeDasharray="3 3"
                />
              )}
              {marked &&
                marked.scores.map((score, index) => (
                  <ReferenceDot
                    key={index}
                    x={marked.x}
                    y={score}
                    fill={SERIES_COLORS[index]}
                    {...ACTIVE_DOT}
                  />
                ))}
            </LineChart>
          </ResponsiveContainer>
        </div>
      </Card.Content>
    </Card>
  );
}
