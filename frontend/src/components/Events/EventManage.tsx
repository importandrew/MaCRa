import {
  Accordion,
  AccordionItem,
  Button,
  Card,
  CardBody,
  Tab,
  Tabs,
  Checkbox,
  CheckboxGroup,
  Form,
  Input,
  Spinner,
  Table,
  TableBody,
  TableCell,
  TableColumn,
  TableHeader,
  TableRow,
  Textarea,
} from '@heroui/react';
import { FormEvent, type ReactNode, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Player } from '../../models/Player';
import { EventRole } from '../../models/EventRole';
import {
  EventVisibility,
  ContentReach,
  CONTENT_REACH_OPTIONS,
} from '../../models/Event';
import {
  addEventMember,
  removeEventMember,
  updateEvent,
  useFetchEventMembers,
} from '../../providers/event/eventApi';
import { useEventLayout } from '../../providers/event/eventLayoutContext';
import { CrownIcon } from '../Icons/CrownIcon';
import { EyeIcon } from '../Icons/EyeIcon';
import { MahjongTileIcon } from '../Icons/MahjongTileIcon';
import { ScalesIcon } from '../Icons/ScalesIcon';
import { GetRequest, authFetch } from '../Request';
import { showToast } from '../Toast';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import PlayerAutocomplete from '../PlayerAutocomplete';

const ALL_ROLES: EventRole[] = ['player', 'referee', 'observer'];

const ROLE_META: Record<
  EventRole,
  { label: string; icon: (props: { className?: string }) => ReactNode }
> = {
  player: {
    label: 'Player',
    icon: (props) => <MahjongTileIcon {...props} />,
  },
  referee: {
    label: 'Referee',
    icon: (props) => <ScalesIcon {...props} />,
  },
  observer: {
    label: 'Observer',
    icon: (props) => <EyeIcon {...props} />,
  },
};

const VISIBILITY_OPTIONS: { value: EventVisibility; label: string }[] = [
  { value: 'private', label: 'Private' },
  { value: 'public', label: 'Public' },
];

type VisibilityTabsProps<T extends string> = {
  label: string;
  ariaLabel: string;
  options: { value: T; label: string }[];
  selectedKey: T;
  onSelectionChange: (key: T) => void;
  isOptionDisabled?: (value: T) => boolean;
};

function VisibilityTabs<T extends string>({
  label,
  ariaLabel,
  options,
  selectedKey,
  onSelectionChange,
  isOptionDisabled,
}: VisibilityTabsProps<T>) {
  return (
    <div className="flex flex-col gap-1.5">
      <p className="text-sm">{label}</p>
      <Tabs
        aria-label={ariaLabel}
        fullWidth
        selectedKey={selectedKey}
        onSelectionChange={(key) => onSelectionChange(key as T)}
      >
        {options.map((opt) => (
          <Tab
            key={opt.value}
            title={opt.label}
            isDisabled={isOptionDisabled?.(opt.value)}
          />
        ))}
      </Tabs>
    </div>
  );
}

export default function EventManage() {
  const navigate = useNavigate();
  const { eventId, event, hasCapability, refetch, isForbidden } =
    useEventLayout();

  const canManageEvent = hasCapability('manage_event');
  const canManageMembers = hasCapability('manage_members');
  const hasMultipleManagingCapabilities = canManageEvent && canManageMembers;

  const {
    data: members,
    isLoading: membersLoading,
    refetch: refetchMembers,
  } = useFetchEventMembers(canManageMembers ? eventId : null);

  const [editName, setEditName] = useState('');
  const [editDescription, setEditDescription] = useState('');
  const [editVisibility, setEditVisibility] =
    useState<EventVisibility>('private');
  const [editRatingsReach, setEditRatingsReach] =
    useState<ContentReach>('members');
  const [editGamesReach, setEditGamesReach] = useState<ContentReach>('members');
  const [isSaving, setIsSaving] = useState(false);

  const [players, setPlayers] = useState<Player[]>([]);
  const [selectedUserId, setSelectedUserId] = useState<string>('');
  const [selectedRoles, setSelectedRoles] = useState<EventRole[]>([]);
  const [isAddingMember, setIsAddingMember] = useState(false);
  const [activeTab, setActiveTab] = useState<'overview' | 'members'>(
    canManageEvent ? 'overview' : 'members'
  );

  const publicReachDisabled = editVisibility === 'private';
  const isReachOptionDisabled = (value: ContentReach) =>
    value === 'public' && publicReachDisabled;

  useEffect(() => {
    if (event) {
      setEditName(event.name);
      setEditDescription(event.description || '');
      setEditVisibility(event.visibility);
      setEditRatingsReach(event.ratings_reach);
      setEditGamesReach(event.games_reach);
    }
  }, [event]);

  useEffect(() => {
    if (!canManageMembers) return;
    const fetchPlayers = async () => {
      const request = GetRequest('/api/players');
      const data = await authFetch<Player[]>(request);
      if (data) setPlayers(data);
    };
    fetchPlayers();
  }, [canManageMembers]);

  const applyVisibilityPreset = (value: EventVisibility) => {
    const reach: ContentReach = value === 'public' ? 'public' : 'members';
    setEditRatingsReach(reach);
    setEditGamesReach(reach);
  };

  const handleSave = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    setIsSaving(true);
    const ok = await updateEvent(eventId, {
      name: editName.trim(),
      description: editDescription.trim() || undefined,
      visibility: editVisibility,
      ratings_reach: editRatingsReach,
      games_reach: editGamesReach,
    });
    setIsSaving(false);

    if (ok) {
      showToast({ title: 'Event updated', color: 'success' });
      refetch();
    }
  };

  const handleAddMember = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!selectedUserId || !selectedRoles.length) return;

    setIsAddingMember(true);
    const ok = await addEventMember(
      eventId,
      parseInt(selectedUserId, 10),
      selectedRoles
    );
    setIsAddingMember(false);

    if (ok) {
      showToast({ title: 'Member added', color: 'success' });
      setSelectedUserId('');
      setSelectedRoles([]);
      refetchMembers();
    }
  };

  const handleRemoveMember = async (userId: number) => {
    const ok = await removeEventMember(eventId, userId);
    if (ok) {
      showToast({ title: 'Member removed', color: 'success' });
      refetchMembers();
    }
  };

  if (isForbidden || (!canManageEvent && !canManageMembers)) {
    return (
      <LoadingTemplate
        isLoading={false}
        isError={false}
        isForbidden
        onRefresh={() => { }}
        onGoHome={() => navigate(`/event/${eventId}`)}
      >
        {null}
      </LoadingTemplate>
    );
  }

  if (!event) return null;

  return (
    <div className="flex flex-col w-full max-w-[600px] mx-auto px-3 md:px-0 md:my-3 gap-4 min-w-0">
      <h1
        className="text-2xl font-semibold truncate min-w-0"
        title={`Manage — ${event.name}`}
      >
        Manage — {event.name}
      </h1>

      {hasMultipleManagingCapabilities && (
        <Tabs
          aria-label="Manage event sections"
          fullWidth
          selectedKey={activeTab}
          onSelectionChange={(key) =>
            setActiveTab(key as 'overview' | 'members')
          }
        >
          {canManageEvent && <Tab key="overview" title="Overview" />}
          {canManageMembers && <Tab key="members" title="Members" />}
        </Tabs>
      )}

      {activeTab === 'overview' &&
        canManageEvent && (
          <>
            <Card>
              <CardBody className="gap-2">
                <p className="text-default-500">General</p>
                <Form
                  id="edit-event-form"
                  onSubmit={handleSave}
                >
                  <Input
                    label="Name"
                    value={editName}
                    onValueChange={setEditName}
                    isRequired
                    variant="bordered"
                  />
                  <Textarea
                    label="Description"
                    value={editDescription}
                    onValueChange={setEditDescription}
                    variant="bordered"
                  />
                  <div className="w-full flex flex-col gap-2">
                    <p className="text-default-500">Visibility</p>
                    <VisibilityTabs
                      label="General"
                      ariaLabel="Event visibility"
                      options={VISIBILITY_OPTIONS}
                      selectedKey={editVisibility}
                      onSelectionChange={(value) => {
                        setEditVisibility(value);
                        applyVisibilityPreset(value);
                      }}
                    />
                    <div className="flex flex-col gap-5">
                      <VisibilityTabs
                        label="Ratings"
                        ariaLabel="Ratings visibility"
                        options={CONTENT_REACH_OPTIONS}
                        selectedKey={editRatingsReach}
                        onSelectionChange={setEditRatingsReach}
                        isOptionDisabled={isReachOptionDisabled}
                      />
                      <VisibilityTabs
                        label="Live games & results"
                        ariaLabel="Live games and results visibility"
                        options={CONTENT_REACH_OPTIONS}
                        selectedKey={editGamesReach}
                        onSelectionChange={setEditGamesReach}
                        isOptionDisabled={isReachOptionDisabled}
                      />
                    </div>
                  </div>
                </Form>
              </CardBody>
            </Card>
            <Button
              type="submit"
              form="edit-event-form"
              color="primary"
              isLoading={isSaving}
              isDisabled={!editName.trim()}
              className="w-3/4 self-center"
            >
              Save changes
            </Button>
          </>
        )}

      {activeTab === 'members' &&
        canManageMembers && (
          <Accordion
            selectionMode="multiple"
          // defaultExpandedKeys={['members', 'add']}
          >
            <AccordionItem
              key="add"
              aria-label="Add members"
              title="Add members"
            >
              <p className="text-sm text-default-500 mb-3">
                All 4 players need the player role to start a game.
              </p>
              <Form onSubmit={handleAddMember} className="gap-4">
                <PlayerAutocomplete
                  label="User"
                  items={players}
                  selectedKey={selectedUserId || null}
                  onSelectionChange={(id) =>
                    setSelectedUserId(id != null ? String(id) : '')
                  }
                />
                <CheckboxGroup
                  label="Roles"
                  value={selectedRoles}
                  onValueChange={(v) => setSelectedRoles(v as EventRole[])}
                >
                  {ALL_ROLES.map((role) => {
                    const meta = ROLE_META[role];
                    return (
                      <Checkbox
                        key={role}
                        value={role}
                        classNames={{ label: 'flex items-center' }}
                      >
                        <span className="flex items-center gap-2 leading-none">
                          {meta.icon({ className: 'w-6 h-6 shrink-0' })}
                          {meta.label}
                        </span>
                      </Checkbox>
                    );
                  })}
                </CheckboxGroup>
                <Button
                  type="submit"
                  color="primary"
                  isLoading={isAddingMember}
                  isDisabled={!selectedUserId || !selectedRoles.length}
                  className="w-3/4 self-center"
                >
                  Add member
                </Button>
              </Form>
            </AccordionItem>

            <AccordionItem
              key="members"
              aria-label="Members"
              title={`Members [${members?.length}]`}
            >
              {membersLoading ? (
                <Spinner size="sm" />
              ) : !members?.length ? (
                <p className="text-sm text-default-500">No members yet</p>
              ) : (
                <Table
                  hideHeader
                  isStriped
                  removeWrapper
                  className="z-0"
                  aria-label="Event members"
                >
                  <TableHeader>
                    <TableColumn>Member</TableColumn>
                    <TableColumn width={1}>Actions</TableColumn>
                  </TableHeader>
                  <TableBody emptyContent={"No members yet"}>
                    {members.map((member) => {
                      const orderedRoles = ALL_ROLES.filter((role) =>
                        member.roles.includes(role)
                      );
                      const isOwner = member.user_id === event.owner_user_id;

                      return (
                        <TableRow key={member.user_id}>
                          <TableCell>
                            <div className="flex justify-between items-center gap-2 min-w-0">
                              <span className="min-w-0 truncate">
                                {member.name}
                              </span>
                              <div className="flex items-center gap-2 shrink-0 text-default-500">
                                {orderedRoles.map((role) => {
                                  const meta = ROLE_META[role];
                                  return (
                                    <span
                                      key={role}
                                      className="inline-flex shrink-0"
                                      title={meta.label}
                                      aria-label={meta.label}
                                      role="img"
                                    >
                                      {meta.icon({ className: 'w-6 h-6' })}
                                    </span>
                                  );
                                })}
                              </div>
                            </div>
                          </TableCell>
                          <TableCell className="w-0 whitespace-nowrap">
                            {isOwner ? (
                              <span
                                className="flex justify-center text-default-500"
                                title="Owner"
                                aria-label="Owner"
                                role="img"
                              >
                                <CrownIcon className="w-6 h-6" />
                              </span>
                            ) : (
                              <Button
                                size="sm"
                                color="danger"
                                variant="flat"
                                onPress={() =>
                                  handleRemoveMember(member.user_id)
                                }
                              >
                                Remove
                              </Button>
                            )}
                          </TableCell>
                        </TableRow>
                      );
                    })}
                  </TableBody>
                </Table>
              )}
            </AccordionItem>
          </Accordion>
        )}
    </div>
  );
}
