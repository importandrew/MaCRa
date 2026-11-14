import { useTranslation } from 'react-i18next';
import { Button, Card, Form, Spinner, useOverlayState } from '@heroui/react';
import { FormEvent, useState } from 'react';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { useNavigate } from 'react-router-dom';
import LazyAccordion from '../LazyAccordion';
import SegmentedPicker from '../SegmentedPicker';
import type { EventMember, EventRole } from '../../models/EventRole';
import {
  addEventMember,
  removeEventMember,
  updateEvent,
  useFetchEventMembers,
} from '../../providers/event/eventApi';
import { useEventPage } from '../../providers/event/eventLayoutContext';
import { CrownIcon } from '../Icons/CrownIcon';
import { showToast } from '../Toast';
import LoadingTemplate from '../LoadingTemplate/LoadingTemplate';
import PlayerAutocomplete from '../PlayerAutocomplete';
import { usePlayerSearch } from '../../providers/player/playerApi';
import PendingButton from '../PendingButton';
import Confirm from '../Confirm';
import EventSettingsFields from './EventSettingsFields';
import EventRulesForm from './EventRulesForm';
import EventSettingsCard from './EventSettingsCard';
import MemberRolesModal from './MemberRolesModal';
import InviteLink from './InviteLink';
import RoleCheckboxes from './RoleCheckboxes';
import { ALL_ROLES, ROLE_ICONS } from './roles';
import { type EventSettings, isValidSettings } from './eventSettings';

export default function EventManage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { eventId, event, hasCapability, refetch, canView } = useEventPage(
    'manage_event',
    'manage_members'
  );

  const canManageEvent = hasCapability('manage_event');
  const canManageMembers = hasCapability('manage_members');
  // The sections: overview and rules with manage_event, members with
  // manage_members; the picker shows when there is more than one.
  const hasMultipleSections = canManageEvent;

  const {
    data: members,
    isLoading: membersLoading,
    refetch: refetchMembers,
  } = useFetchEventMembers(canManageMembers ? eventId : null);

  // draft holds the owner's edits; until there are any (and again after a
  // save), the form shows the event's own settings.
  const [draft, setDraft] = useState<EventSettings | null>(null);
  const settings: EventSettings | null =
    draft ??
    (event && {
      name: event.name,
      description: event.description || '',
      visibility: event.visibility,
      ratingsReach: event.ratings_reach,
      gamesReach: event.games_reach,
    });
  // A finished event: only its visibility can change, members are final.
  const isFinished = Boolean(event?.finished_at);
  const { isPending: isSaving, run: runSave } = useSubmitOnce();

  // The add-member picker searches the server as the user types.
  const { players, search: searchPlayers } = usePlayerSearch();
  const [selectedUserId, setSelectedUserId] = useState<string>('');
  const [selectedRoles, setSelectedRoles] = useState<EventRole[]>([]);
  const { isPending: isAddingMember, run: runAddMember } = useSubmitOnce();
  const { run: runRemoveMember } = useSubmitOnce();
  const removeConfirm = useOverlayState();
  const [memberToRemove, setMemberToRemove] = useState<EventMember | null>(
    null
  );
  // memberToEdit: the member whose roles are being switched (a tap on them).
  const [memberToEdit, setMemberToEdit] = useState<EventMember | null>(null);
  const [activeTab, setActiveTab] = useState<
    'overview' | 'settings' | 'members'
  >(canManageEvent ? 'overview' : 'members');

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
        // The edits stay on screen until the event has reloaded with them:
        // clearing them first would show the old values for a moment.
        await refetch();
        setDraft(null);
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

      {hasMultipleSections && (
        <SegmentedPicker
          ariaLabel={t('events.manageSections')}
          options={[
            { value: 'overview', label: t('events.overview') },
            { value: 'settings', label: t('eventSettings.title') },
            ...(canManageMembers
              ? [{ value: 'members' as const, label: t('members.title') }]
              : []),
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
                    onChange={setDraft}
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

      {activeTab === 'settings' && canManageEvent && (
        <>
          <EventRulesForm event={event} onSaved={refetch} />
          <EventSettingsCard event={event} onSaved={refetch} />
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
                    id: 'invite',
                    title: t('invite.section'),
                    content: () => <InviteLink eventId={eventId} />,
                  },
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
                            onSearch={searchPlayers}
                            description={t('members.searchHint')}
                            selectedKey={selectedUserId || null}
                            onSelectionChange={(id) =>
                              setSelectedUserId(id != null ? String(id) : '')
                            }
                          />
                          <RoleCheckboxes
                            value={selectedRoles}
                            onChange={setSelectedRoles}
                          />
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
                          {/* A tap opens the role switcher; a finished
                              event's members can't change. */}
                          <button
                            type="button"
                            disabled={isFinished}
                            aria-label={t('members.editRoles', {
                              name: member.name,
                            })}
                            onClick={() => setMemberToEdit(member)}
                            className="flex min-w-0 items-center justify-between gap-2 rounded-lg text-left enabled:cursor-pointer enabled:hover:bg-default"
                          >
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
                          </button>
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
      <MemberRolesModal
        eventId={eventId}
        member={memberToEdit}
        isOwner={memberToEdit?.user_id === event.owner_user_id}
        onClose={() => setMemberToEdit(null)}
        // The overview's player count and the capabilities may change too.
        onSaved={() => Promise.all([refetchMembers(), refetch()])}
      />
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
