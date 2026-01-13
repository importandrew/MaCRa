import {
  Dropdown,
  DropdownTrigger,
  Avatar,
  DropdownMenu,
  DropdownItem,
} from '@heroui/react';
import { PlayerOrder } from '../../../models/GameData';
import { WarningIcon } from '../../Icons/WarningIcon';
import { UserIcon } from '../../Icons/UserIcon';

export default function PlayerMenu({
  player,
  setFinedPlayerId,
}: {
  player: PlayerOrder | undefined;
  setFinedPlayerId: (id: number) => void;
}) {
  const assignPenalty = (_event: unknown, id: number) => {
    setFinedPlayerId(id);
  };

  return (
    <Dropdown>
      <DropdownTrigger>
        <Avatar
          as="button"
          className="transition-transform"
          color="default"
          size="md"
          name={player?.name || ''}
        />
      </DropdownTrigger>
      {player && (
        <DropdownMenu
          disabledKeys={['profile']}
          aria-label="Player menu"
          variant="flat"
        >
          <DropdownItem
            key="penalty"
            color="warning"
            onPress={(event) => assignPenalty(event, player.id)}
            startContent={<WarningIcon className="text-warning" />}
          >
            Penalty
          </DropdownItem>
          <DropdownItem
            key="profile"
            startContent={<UserIcon className="text-primary" />}
          >
            Profile
          </DropdownItem>
        </DropdownMenu>
      )}
    </Dropdown>
  );
}
