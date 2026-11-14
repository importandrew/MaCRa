import { useTranslation } from 'react-i18next';
import { Button, Modal } from '@heroui/react';
import { useState } from 'react';
import type { EventMember, EventRole } from '../../models/EventRole';
import { setMemberRoles } from '../../providers/event/eventApi';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';
import RoleCheckboxes from './RoleCheckboxes';

// MemberRolesModal switches a member's roles (tap a member in Manage →
// Members). Every member keeps at least one role (removing is a separate
// action); the owner may drop all of theirs and stays the owner. The server
// refuses dropping the player role of someone in a live game.
export default function MemberRolesModal({
  eventId,
  member,
  isOwner,
  onClose,
  onSaved,
}: {
  eventId: number;
  // member: the member edited; null closes the modal.
  member: EventMember | null;
  isOwner: boolean;
  onClose: () => void;
  onSaved: () => Promise<unknown>;
}) {
  const { t } = useTranslation();
  // draft is the edited roles of the member it was started for; another
  // member starts from their own roles.
  const [draft, setDraft] = useState<{
    userId: number;
    roles: EventRole[];
  } | null>(null);
  const roles =
    member && draft?.userId === member.user_id ? draft.roles : member?.roles;
  const { isPending, run } = useSubmitOnce();

  const isUnchanged =
    !member ||
    !roles ||
    (roles.length === member.roles.length &&
      roles.every((role) => member.roles.includes(role)));
  const isValid = Boolean(roles && (roles.length > 0 || isOwner));

  const save = () =>
    run(async () => {
      if (!member || !roles) return;
      if (await setMemberRoles(eventId, member.user_id, roles)) {
        showToast({ title: t('members.rolesSaved'), variant: 'success' });
        await onSaved();
        setDraft(null);
        onClose();
      }
    });

  return (
    <Modal.Backdrop
      isOpen={member != null}
      onOpenChange={(open) => {
        if (!open) {
          setDraft(null);
          onClose();
        }
      }}
    >
      <Modal.Container placement="center">
        <Modal.Dialog>
          <Modal.CloseTrigger aria-label={t('common.close')} />
          <Modal.Header>
            <Modal.Heading className="min-w-0 truncate">
              {t('members.rolesOf', { name: member?.name ?? '' })}
            </Modal.Heading>
          </Modal.Header>
          <Modal.Body className="flex flex-col gap-3">
            {member && roles && (
              <RoleCheckboxes
                value={roles}
                onChange={(next) =>
                  setDraft({ userId: member.user_id, roles: next })
                }
              />
            )}
            <p className="text-sm text-muted">
              {isOwner ? t('members.ownerRolesHint') : t('members.rolesHint')}
            </p>
          </Modal.Body>
          <Modal.Footer>
            <Button slot="close" variant="tertiary">
              {t('common.cancel')}
            </Button>
            <PendingButton
              variant="primary"
              isPending={isPending}
              isDisabled={isUnchanged || !isValid}
              onPress={save}
            >
              {t('common.save')}
            </PendingButton>
          </Modal.Footer>
        </Modal.Dialog>
      </Modal.Container>
    </Modal.Backdrop>
  );
}
