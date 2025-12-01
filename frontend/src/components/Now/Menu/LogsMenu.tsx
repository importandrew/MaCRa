import GameLogs from '../../End/GameLogs';

export default function LogsMenu({ gameId }: { gameId: number | null }) {
  return (
    <div className="flex flex-col justify-center mb-3">
      <GameLogs
        isChartVisible={false}
        isLogVisible={true}
        gameId={gameId}
      ></GameLogs>
    </div>
  );
}
