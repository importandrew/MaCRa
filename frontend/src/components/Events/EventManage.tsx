import { useTranslation } from 'react-i18next';
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
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { useNavigate } from 'react-router-dom';
import LazyAccordion from '../LazyAccordion';
import SegmentedPicker from '../SegmentedPicker';
import { Player } from '../../models/Player';
import { EventRole } from '../../models/EventRole';
import {
  EventVisibility,
  ContentReach,
  CONTENT_REACHES,
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

const ROLE_ICONS: Record<
  EventRole,
  (props: { className?: string }) => ReactNode
> = {
  player: (props) => <MahjongTileIcon {...props} />,
  referee: (props) => <ScalesIcon {...props} />,
  observer: (props) => <EyeIcon {...props} />,
};

const VISIBILITIES: EventVisibility[] = ['private', 'public'];

export default function EventManage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const reachOptions = CONTENT_REACHES.map((value) => ({
    value,
    label: t(`reach.${value}`),
  }));
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
  const { isPending: isSaving, run: runSave } = useSubmitOnce();

  const [players, setPlayers] = useState<Player[]>([]);
  const [selectedUserId, setSelectedUserId] = useState<string>('');
  const [selectedRoles, setSelectedRoles] = useState<EventRole[]>([]);
  const { isPending: isAddingMember, run: runAddMember } = useSubmitOnce();
  const { run: runRemoveMember } = useSubmitOnce();
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

  const handleSave = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    void runSave(async () => {
      const ok = await updateEvent(eventId, {
        name: editName.trim(),
        description: editDescription.trim() || undefined,
        visibility: editVisibility,
        ratings_reach: editRatingsReach,
        games_reach: editGamesReach,
      });

      if (ok) {
        showToast({ title: t('events.updated'), variant: 'success' });
        refetch();
      }
    });
  };

  const handleAddMember = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!selectedUserId || !selectedRoles.length) return;

    void runAddMember(async () => {
      const ok = await addEventMember(
        eventId,
        parseInt(selectedUserId, 10),
        selectedRoles
      );

      if (ok) {
        showToast({ title: t('members.added'), variant: 'success' });
        setSelectedUserId('');
        setSelectedRoles([]);
        refetchMembers();
        // The overview's player count comes with the event.
        refetch();
      }
    });
  };

  const handleRemoveMember = (userId: number) =>
    runRemoveMember(async () => {
      const ok = await removeEventMember(eventId, userId);
      if (ok) {
        showToast({ title: t('members.removed'), variant: 'success' });
        refetchMembers();
        refetch();
      }
    });

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
        title={t('events.manageTitle', { name: event.name })}
      >
        {t('events.manageTitle', { name: event.name })}
      </h1>

      {hasMultipleManagingCapabilities && (
        <SegmentedPicker
          ariaLabel={t('events.manageSections')}
          options={[
            { value: 'overview', label: t('events.overview') },
            { value: 'members', label: t('members.title') },
          ]}
          selectedKey={activeTab}
          onSelectionChange={setActiveTab}
        />
      )}

      {activeTab === 'overview' && canManageEvent && (
        <>
          <Card>
            <Card.Content className="gap-2">
              <p className="text-muted">{t('events.general')}</p>
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
                  <Label>{t('events.name')}</Label>
                  <Input />
                </TextField>
                <TextField
                  fullWidth
                  value={editDescription}
                  onChange={setEditDescription}
                >
                  <Label>{t('events.description')}</Label>
                  <TextArea />
                </TextField>
                <div className="w-full flex flex-col gap-2">
                  <p className="text-muted">{t('events.visibility')}</p>
                  <SegmentedPicker
                    label={t('events.general')}
                    options={VISIBILITIES.map((value) => ({
                      value,
                      label: t(`visibility.${value}`),
                    }))}
                    selectedKey={editVisibility}
                    onSelectionChange={(value) => {
                      setEditVisibility(value);
                      applyVisibilityPreset(value);
                    }}
                  />
                  <div className="flex flex-col gap-5">
                    <SegmentedPicker
                      label={t('events.ratingsReach')}
                      options={reachOptions}
                      selectedKey={editRatingsReach}
                      onSelectionChange={setEditRatingsReach}
                      isOptionDisabled={isReachOptionDisabled}
                    />
                    <SegmentedPicker
                      label={t('events.gamesReach')}
                      options={reachOptions}
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
            {t('events.saveChanges')}
          </PendingButton>
        </>
      )}

      {activeTab === 'members' && canManageMembers && (
        <LazyAccordion
          allowsMultipleExpanded
          sections={[
            {
              id: 'add',
              title: t('members.addTitle'),
              content: () => (
                <>
                  <p className="text-sm text-muted mb-3">
                    {t('members.playerRoleHint')}
                  </p>
                  <Form
                    onSubmit={handleAddMember}
                    className="flex flex-col gap-4"
                  >
                    <PlayerAutocomplete
                      label={t('members.user')}
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
                      <Label>{t('members.roles')}</Label>
                      {ALL_ROLES.map((role) => (
                        <Checkbox key={role} value={role}>
                          <Checkbox.Content>
                            <Checkbox.Control>
                              <Checkbox.Indicator />
                            </Checkbox.Control>
                            <span className="flex items-center gap-2 leading-none">
                              {ROLE_ICONS[role]({
                                className: 'w-6 h-6 shrink-0',
                              })}
                              {t(`roles.${role}`)}
                            </span>
                          </Checkbox.Content>
                        </Checkbox>
                      ))}
                    </CheckboxGroup>
                    <PendingButton
                      type="submit"
                      variant="primary"
                      isPending={isAddingMember}
                      isDisabled={!selectedUserId || !selectedRoles.length}
                      className="w-3/4 self-center"
                    >
                      {t('members.add')}
                    </PendingButton>
                  </Form>
                </>
              ),
            },
            {
              id: 'members',
              title: t('members.titleWithCount', {
                n: members?.length ?? 0,
              }),
              content: () =>
                membersLoading ? (
                  <Spinner size="sm" />
                ) : !members?.length ? (
                  <p className="text-sm text-muted">{t('members.none')}</p>
                ) : (
                  // One grid for all rows (rows use subgrid), like the v2 table:
                  // name + roles in the first column, the action in the second,
                  // so the action column lines up and names use all the room
                  // their roles leave.
                  <ul
                    aria-label={t('members.title')}
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
                              ).map((role) => (
                                <span
                                  key={role}
                                  className="inline-flex shrink-0"
                                  title={t(`roles.${role}`)}
                                  aria-label={t(`roles.${role}`)}
                                  role="img"
                                >
                                  {ROLE_ICONS[role]({ className: 'w-6 h-6' })}
                                </span>
                              ))}
                            </div>
                          </div>
                          <div className="flex justify-center whitespace-nowrap">
                            {isOwner ? (
                              <span
                                className="flex justify-center text-muted"
                                title={t('roles.owner')}
                                aria-label={t('roles.owner')}
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
                                {t('members.remove')}
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
