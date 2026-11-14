import type { ReactNode } from 'react';
import type { EventRole } from '../../models/EventRole';
import { EyeIcon } from '../Icons/EyeIcon';
import { MahjongTileIcon } from '../Icons/MahjongTileIcon';
import { ScalesIcon } from '../Icons/ScalesIcon';

// ALL_ROLES: the member roles, in the order screens show them.
export const ALL_ROLES: EventRole[] = ['player', 'referee', 'observer'];

export const ROLE_ICONS: Record<
  EventRole,
  (props: { className?: string }) => ReactNode
> = {
  player: (props) => <MahjongTileIcon {...props} />,
  referee: (props) => <ScalesIcon {...props} />,
  observer: (props) => <EyeIcon {...props} />,
};
