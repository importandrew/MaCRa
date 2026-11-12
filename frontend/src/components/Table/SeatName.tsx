import UserAvatar from '../UserAvatar';

// SeatName is the top line of a seat at the table: the player's wind, avatar
// and name.
export default function SeatName({
  player,
}: {
  player: { name: string; wind: string };
}) {
  return (
    <span className="flex justify-center items-center gap-3 text-ellipsis overflow-hidden">
      <span className="shrink-0">{player.wind}</span>
      <UserAvatar name={player.name} className="shrink-0" />
      <span className="min-w-0 truncate">{player.name}</span>
    </span>
  );
}
