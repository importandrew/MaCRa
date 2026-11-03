import {
  Button,
  Card,
  Checkbox,
  CheckboxGroup,
  Form,
  Input,
  Label,
  Spinner,
  TextArea,
  TextField,
} from '@heroui/react';
import { FormEvent, type ReactNode, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import LazyAccordion from '../LazyAccordion';
import SegmentedPicker from '../SegmentedPicker';
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
import PendingButton from '../PendingButton';

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
      showToast({ title: 'Event updated', variant: 'success' });
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
      showToast({ title: 'Member added', variant: 'success' });
      setSelectedUserId('');
      setSelectedRoles([]);
      refetchMembers();
    }
  };

  const handleRemoveMember = async (userId: number) => {
    const ok = await removeEventMember(eventId, userId);
    if (ok) {
      showToast({ title: 'Member removed', variant: 'success' });
      refetchMembers();
    }
  };

  if (isForbidden || (!canManageEvent && !canManageMembers)) {
    return (
      <LoadingTemplate
        isLoading={false}
        isError={false}
        isForbidden
        onRefresh={() => {}}
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
        <SegmentedPicker
          ariaLabel="Manage event sections"
          options={[
            { value: 'overview', label: 'Overview' },
            { value: 'members', label: 'Members' },
          ]}
          selectedKey={activeTab}
          onSelectionChange={setActiveTab}
        />
      )}

      {activeTab === 'overview' && canManageEvent && (
        <>
          <Card>
            <Card.Content className="gap-2">
              <p className="text-muted">General</p>
              <Form
                id="edit-event-form"
                onSubmit={handleSave}
                className="flex flex-col gap-3"
              >
                <TextField
                  fullWidth
                  isRequired
                  value={editName}
                  onChange={setEditName}
                >
                  <Label>Name</Label>
                  <Input />
                </TextField>
                <TextField
                  fullWidth
                  value={editDescription}
                  onChange={setEditDescription}
                >
                  <Label>Description</Label>
                  <TextArea />
                </TextField>
                <div className="w-full flex flex-col gap-2">
                  <p className="text-muted">Visibility</p>
                  <SegmentedPicker
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
                    <SegmentedPicker
                      label="Ratings"
                      ariaLabel="Ratings visibility"
                      options={CONTENT_REACH_OPTIONS}
                      selectedKey={editRatingsReach}
                      onSelectionChange={setEditRatingsReach}
                      isOptionDisabled={isReachOptionDisabled}
                    />
                    <SegmentedPicker
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
            </Card.Content>
          </Card>
          <PendingButton
            type="submit"
            form="edit-event-form"
            variant="primary"
            isPending={isSaving}
            isDisabled={!editName.trim()}
            className="w-3/4 self-center"
          >
            Save changes
          </PendingButton>
        </>
      )}

      {activeTab === 'members' && canManageMembers && (
        <LazyAccordion
          allowsMultipleExpanded
          sections={[
            {
              id: 'add',
              title: 'Add members',
              content: () => (
                <>
                  <p className="text-sm text-muted mb-3">
                    All 4 players need the player role to start a game.
                  </p>
                  <Form
                    onSubmit={handleAddMember}
                    className="flex flex-col gap-4"
                  >
                    <PlayerAutocomplete
                      label="User"
                      items={players}
                      selectedKey={selectedUserId || null}
                      onSelectionChange={(id) =>
                        setSelectedUserId(id != null ? String(id) : '')
                      }
                    />
                    <CheckboxGroup
                      value={selectedRoles}
                      onChange={(v) => setSelectedRoles(v as EventRole[])}
                    >
                      <Label>Roles</Label>
                      {ALL_ROLES.map((role) => {
                        const meta = ROLE_META[role];
                        return (
                          <Checkbox key={role} value={role}>
                            <Checkbox.Content>
                              <Checkbox.Control>
                                <Checkbox.Indicator />
                              </Checkbox.Control>
                              <span className="flex items-center gap-2 leading-none">
                                {meta.icon({ className: 'w-6 h-6 shrink-0' })}
                                {meta.label}
                              </span>
                            </Checkbox.Content>
                          </Checkbox>
                        );
                      })}
                    </CheckboxGroup>
                    <PendingButton
                      type="submit"
                      variant="primary"
                      isPending={isAddingMember}
                      isDisabled={!selectedUserId || !selectedRoles.length}
                      className="w-3/4 self-center"
                    >
                      Add member
                    </PendingButton>
                  </Form>
                </>
              ),
            },
            {
              id: 'members',
              title: `Members [${members?.length ?? 0}]`,
              content: () =>
                membersLoading ? (
                  <Spinner size="sm" />
                ) : !members?.length ? (
                  <p className="text-sm text-muted">No members yet</p>
                ) : (
                  // One grid for all rows (rows use subgrid), like the v2 table:
                  // name + roles in the first column, the action in the second,
                  // so the action column lines up and names use all the room
                  // their roles leave.
                  <ul
                    aria-label="Event members"
                    className="grid grid-cols-[minmax(0,1fr)_auto] gap-x-3"
                  >
                    {members.map((member) => {
                      const isOwner = member.user_id === event.owner_user_id;

                      return (
                        <li
                          key={member.user_id}
                          className="col-span-2 grid grid-cols-subgrid items-center px-2 py-2 border-b border-separator text-foreground last:border-b-0"
                        >
                          <div className="flex min-w-0 items-center justify-between gap-2">
                            <span className="min-w-0 truncate">
                              {member.name}
                            </span>
                            <div className="flex shrink-0 items-center gap-2 text-muted">
                              {ALL_ROLES.filter((role) =>
                                member.roles.includes(role)
                              ).map((role) => {
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
                          <div className="flex justify-center whitespace-nowrap">
                            {isOwner ? (
                              <span
                                className="flex justify-center text-muted"
                                title="Owner"
                                aria-label="Owner"
                                role="img"
                              >
                                <CrownIcon className="w-6 h-6" />
                              </span>
                            ) : (
                              <Button
                                size="sm"
                                variant="danger-soft"
                                onPress={() =>
                                  handleRemoveMember(member.user_id)
                                }
                              >
                                Remove
                              </Button>
                            )}
                          </div>
                        </li>
                      );
                    })}
                  </ul>
                ),
            },
          ]}
        />
      )}
    </div>
  );
}
