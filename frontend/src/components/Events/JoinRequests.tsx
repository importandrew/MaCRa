import { useTranslation } from 'react-i18next';
import { Button, Spinner, Switch, useOverlayState } from '@heroui/react';
import { useState } from 'react';
import type { Event } from '../../models/Event';
import type { JoinRequest, JoiningSettings } from '../../models/JoinRequest';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import {
  answerJoinRequest,
  updateJoining,
  useFetchJoinRequests,
} from '../../providers/event/joinRequestApi';
import { getDateTimeFormatted } from '../../utils/date-tools';
import Confirm from '../Confirm';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';
import ValueStepper from '../ValueStepper';

// DEFAULT_PLAYER_LIMIT is where the limit starts when the owner turns it on
// (eight tables).
const DEFAULT_PLAYER_LIMIT = 32;

function SettingSwitch({
  isSelected,
  onChange,
  isDisabled,
  children,
}: {
  isSelected: boolean;
  onChange: (value: boolean) => void;
  isDisabled?: boolean;
  children: string;
}) {
  return (
    <Switch
      size="sm"
      isSelected={isSelected}
      onChange={onChange}
      isDisabled={isDisabled}
    >
      <Switch.Content>
        <Switch.Control>
          <Switch.Thumb />
        </Switch.Control>
        {children}
      </Switch.Content>
    </Switch>
  );
}

// JoiningForm sets whether the event takes requests to join and its player
// limit, saved together.
function JoiningForm({
  event,
  onSaved,
}: {
  event: Event;
  onSaved: () => Promise<unknown>;
}) {
  const { t } = useTranslation();
  const saved: JoiningSettings = {
    accept_requests: event.accept_requests,
    player_limit: event.player_limit,
  };
  // draft holds the owner's edits; until there are any (and again after a
  // save), the form shows the event's own values.
  const [draft, setDraft] = useState<JoiningSettings | null>(null);
  const value = draft ?? saved;
  const { isPending, run } = useSubmitOnce();
  const isPublic = event.visibility === 'public';
  const isChanged =
    value.accept_requests !== saved.accept_requests ||
    value.player_limit !== saved.player_limit;
  const isValid = value.player_limit == null || value.player_limit >= 1;

  const save = () =>
    run(async () => {
      if (await updateJoining(event.id, value)) {
        showToast({ title: t('joinRequests.saved'), variant: 'success' });
        await onSaved();
        setDraft(null);
      }
    });

  return (
    <div className="flex flex-col gap-3">
      <SettingSwitch
        isSelected={value.accept_requests}
        onChange={(on) => setDraft({ ...value, accept_requests: on })}
      >
        {t('joinRequests.accept')}
      </SettingSwitch>
      <p className="text-sm text-muted">
        {isPublic
          ? t('joinRequests.acceptHint')
          : t('joinRequests.privateHint')}
      </p>
      <SettingSwitch
        isSelected={value.player_limit != null}
        onChange={(on) =>
          setDraft({
            ...value,
            player_limit: on ? DEFAULT_PLAYER_LIMIT : null,
          })
        }
      >
        {t('joinRequests.limit')}
      </SettingSwitch>
      {value.player_limit != null && (
        <ValueStepper
          label={t('joinRequests.limitLabel')}
          value={value.player_limit}
          min={1}
          onChange={(limit) => setDraft({ ...value, player_limit: limit })}
        />
      )}
      <p className="text-sm text-muted">{t('joinRequests.limitHint')}</p>
      <PendingButton
        variant="primary"
        isPending={isPending}
        isDisabled={!isChanged || !isValid}
        onPress={save}
        className="w-3/4 self-center"
      >
        {t('joinRequests.save')}
      </PendingButton>
    </div>
  );
}

// JoinRequests is the "Requests to join" section of Manage → Members: the
// switch and the player limit, then the pending requests, oldest first, each
// accepted (the user becomes a player) or declined after a confirmation (the
// user can't ask again).
export default function JoinRequests({
  event,
  onChanged,
}: {
  event: Event;
  // onChanged reloads what an answer changes: the members, the event's
  // player count and the pending count.
  onChanged: () => Promise<unknown>;
}) {
  const { t } = useTranslation();
  const { requests, isLoading, isError, refetch } = useFetchJoinRequests(
    event.id
  );
  const { run } = useSubmitOnce();
  const [answering, setAnswering] = useState<number | null>(null);
  const declineConfirm = useOverlayState();
  const [toDecline, setToDecline] = useState<JoinRequest | null>(null);

  const answer = (request: JoinRequest, reply: 'accept' | 'decline') =>
    run(async () => {
      setAnswering(request.user_id);
      const ok = await answerJoinRequest(event.id, request.user_id, reply);
      if (ok) {
        showToast({
          title: t(
            reply === 'accept'
              ? 'joinRequests.accepted'
              : 'joinRequests.declinedToast',
            { name: request.name }
          ),
          variant: 'success',
        });
      }
      // Also after an error: the request may be gone (withdrawn).
      await Promise.all([refetch(), onChanged()]);
      setAnswering(null);
    });

  return (
    <div className="flex flex-col gap-4">
      <JoiningForm event={event} onSaved={onChanged} />
      <h3 className="font-medium">{t('joinRequests.pendingTitle')}</h3>
      {isLoading ? (
        <Spinner size="sm" />
      ) : isError ? (
        <p className="text-sm text-danger">{t('joinRequests.loadFailed')}</p>
      ) : !requests?.length ? (
        <p className="text-sm text-muted">{t('joinRequests.none')}</p>
      ) : (
        <ul aria-label={t('joinRequests.pendingTitle')}>
          {requests.map((request) => (
            <li
              key={request.user_id}
              className="flex flex-wrap items-center justify-between gap-2 border-b border-separator px-2 py-2 text-foreground last:border-b-0"
            >
              <div className="flex min-w-0 flex-col">
                <span className="truncate">{request.name}</span>
                <span className="text-xs text-muted">
                  {getDateTimeFormatted(request.created_at)}
                </span>
              </div>
              <div className="flex gap-2">
                <PendingButton
                  size="sm"
                  variant="primary"
                  isPending={answering === request.user_id}
                  onPress={() => answer(request, 'accept')}
                >
                  {t('joinRequests.acceptOne')}
                </PendingButton>
                <Button
                  size="sm"
                  variant="danger-soft"
                  isDisabled={answering === request.user_id}
                  onPress={() => {
                    setToDecline(request);
                    declineConfirm.open();
                  }}
                >
                  {t('joinRequests.decline')}
                </Button>
              </div>
            </li>
          ))}
        </ul>
      )}
      <Confirm
        isOpen={declineConfirm.isOpen}
        onClose={declineConfirm.close}
        onConfirm={() => {
          declineConfirm.close();
          if (toDecline) void answer(toDecline, 'decline');
        }}
        headerText={t('joinRequests.declineTitle')}
        prompt={t('joinRequests.declinePrompt', {
          name: toDecline?.name ?? '',
        })}
        cancelText={t('common.cancel')}
        confirmText={t('joinRequests.decline')}
        confirmVariant="danger"
      />
    </div>
  );
}
