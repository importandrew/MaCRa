import { Dropdown, Label } from '@heroui/react';
import { PlayerOrder } from '../../../models/GameData';
import { UserIcon } from '../../Icons/UserIcon';
import UserAvatar from '../../UserAvatar';

export default function PlayerMenu({
  player,
}: {
  player: PlayerOrder | undefined;
}) {
  return (
    <Dropdown>
      <Dropdown.Trigger
        className="shrink-0 rounded-full"
        aria-label={`${player?.name ?? 'Player'} menu`}
      >
        <UserAvatar name={player?.name || ''} />
      </Dropdown.Trigger>
      {player && (
        <Dropdown.Popover>
          <Dropdown.Menu aria-label="Player menu" disabledKeys={['profile']}>
            <Dropdown.Item id="profile" textValue="Profile">
              <UserIcon className="text-accent" />
              <Label>Profile</Label>
            </Dropdown.Item>
          </Dropdown.Menu>
        </Dropdown.Popover>
      )}
    </Dropdown>
  );
}
