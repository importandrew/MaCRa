import { useTranslation } from 'react-i18next';
import {
  Button,
  Card,
  Checkbox,
  CheckboxGroup,
  Form,
  Label,
  Spinner,
  useOverlayState,
} from '@heroui/react';
import { FormEvent, type ReactNode, useEffect, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { useNavigate } from 'react-router-dom';
import LazyAccordion from '../LazyAccordion';
import SegmentedPicker from '../SegmentedPicker';
import { Player } from '../../models/Player';
import type { EventMember, EventRole } from '../../models/EventRole';
import {
  addEventMember,
  removeEventMember,
  updateEvent,
  useFetchEventMembers,
} from '../../providers/event/eventApi';
import { useEventPage } from '../../providers/event/eventLayoutContext';
import { CrownIcon } from '../Icons/CrownIcon';
import { EyeIcon } from '../Icons/EyeIcon';
import { MahjongTileIcon } from '../Icons/MahjongTileIcon';
import { ScalesIcon } from '../Icons/ScalesIcon';
import { GetRequest, authFetch } from '../Request';
import { showToast } from '../Toast';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import PlayerAutocomplete from '../PlayerAutocomplete';
import PendingButton from '../PendingButton';
import Confirm from '../Confirm';
import EventSettingsFields from './EventSettingsFields';
import { type EventSettings, isValidSettings } from './eventSettings';

const ALL_ROLES: EventRole[] = ['player', 'referee', 'observer'];

const ROLE_ICONS: Record<
  EventRole,
  (props: { className?: string }) => ReactNode
> = {
  player: (props) => <MahjongTileIcon {...props} />,
  referee: (props) => <ScalesIcon {...props} />,
  observer: (props) => <EyeIcon {...props} />,
};

export default function EventManage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { eventId, event, hasCapability, refetch, canView } = useEventPage(
    'manage_event',
    'manage_members'
  );

  const canManageEvent = hasCapability('manage_event');
  const canManageMembers = hasCapability('manage_members');
  const hasMultipleManagingCapabilities = canManageEvent && canManageMembers;

  const {
    data: members,
    isLoading: membersLoading,
    refetch: refetchMembers,
  } = useFetchEventMembers(canManageMembers ? eventId : null);

  const [settings, setSettings] = useState<EventSettings | null>(null);
  // A finished event: only its visibility can change, members are final.
  const isFinished = Boolean(event?.finished_at);
  const { isPending: isSaving, run: runSave } = useSubmitOnce();

  const [players, setPlayers] = useState<Player[]>([]);
  const [selectedUserId, setSelectedUserId] = useState<string>('');
  const [selectedRoles, setSelectedRoles] = useState<EventRole[]>([]);
  const { isPending: isAddingMember, run: runAddMember } = useSubmitOnce();
  const { run: runRemoveMember } = useSubmitOnce();
  const removeConfirm = useOverlayState();
  const [memberToRemove, setMemberToRemove] = useState<EventMember | null>(
    null
  );
  const [activeTab, setActiveTab] = useState<'overview' | 'members'>(
    canManageEvent ? 'overview' : 'members'
  );

  useEffect(() => {
    if (event) {
      setSettings({
        name: event.name,
        description: event.description || '',
        visibility: event.visibility,
        ratingsReach: event.ratings_reach,
        gamesReach: event.games_reach,
      });
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

  const handleSave = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!settings) return;
    void runSave(async () => {
      // A finished event takes only its visibility. An empty description
      // is sent too: it clears the description.
      const ok = await updateEvent(
        eventId,
        isFinished
          ? { visibility: settings.visibility }
          : {
              name: settings.name.trim(),
              description: settings.description.trim(),
              visibility: settings.visibility,
              ratings_reach: settings.ratingsReach,
              games_reach: settings.gamesReach,
            }
      );

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

  const askRemoveMember = (member: EventMember) => {
    setMemberToRemove(member);
    removeConfirm.open();
  };

  const handleRemoveMember = () =>
    runRemoveMember(async () => {
      removeConfirm.close();
      if (!memberToRemove) return;
      const ok = await removeEventMember(eventId, memberToRemove.user_id);
      if (ok) {
        showToast({ title: t('members.removed'), variant: 'success' });
        refetchMembers();
        refetch();
      }
    });

  if (!canView) {
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
                {isFinished && (
                  <p className="text-sm text-muted">
                    {t('events.finishedReadOnly')}
                  </p>
                )}
                {settings && (
                  <EventSettingsFields
                    value={settings}
                    onChange={setSettings}
                    isFinished={isFinished}
                  />
                )}
              </Form>
            </Card.Content>
          </Card>
          <PendingButton
            type="submit"
            form="edit-event-form"
            variant="primary"
            isPending={isSaving}
            isDisabled={!settings || !isValidSettings(settings)}
            className="w-3/4 self-center"
          >
            {t('events.saveChanges')}
          </PendingButton>
        </>
      )}

      {activeTab === 'members' && canManageMembers && isFinished && (
        <p className="text-sm text-muted">{t('members.finishedReadOnly')}</p>
      )}

      {activeTab === 'members' && canManageMembers && (
        <LazyAccordion
          allowsMultipleExpanded
          sections={[
            ...(isFinished
              ? []
              : [
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
                            isDisabled={
                              !selectedUserId || !selectedRoles.length
                            }
                            className="w-3/4 self-center"
                          >
                            {t('members.add')}
                          </PendingButton>
                        </Form>
                      </>
                    ),
                  },
                ]),
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
                            {isFinished ? null : isOwner ? (
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
                                onPress={() => askRemoveMember(member)}
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
      <Confirm
        isOpen={removeConfirm.isOpen}
        onClose={removeConfirm.close}
        onConfirm={handleRemoveMember}
        headerText={t('members.removeConfirmTitle')}
        prompt={t('members.removeConfirmPrompt', {
          name: memberToRemove?.name ?? '',
        })}
        cancelText={t('common.cancel')}
        confirmText={t('members.remove')}
        confirmVariant="danger"
      />
    </div>
  );
}
