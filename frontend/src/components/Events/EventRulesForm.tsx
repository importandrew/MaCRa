import { useTranslation } from 'react-i18next';
import { Card, Form } from '@heroui/react';
import { type FormEvent, useState } from 'react';
import {
  type GameRules,
  isValidRules,
  sameRules,
} from '../../models/GameRules';
import type { Event } from '../../models/Event';
import { updateEventRules } from '../../providers/event/eventApi';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import PendingButton from '../PendingButton';
import { showToast } from '../Toast';
import GameRulesFields from './GameRulesFields';

// EventRulesForm is Manage → Rules: the event's game rules. They can change
// only until the first game (the server refuses them after it) and not in a
// finished event; then the form only shows them.
export default function EventRulesForm({
  event,
  onSaved,
}: {
  event: Event;
  onSaved: () => void;
}) {
  const { t } = useTranslation();
  // draft holds the owner's edits; until there are any (and again after a
  // save), the form shows the event's rules.
  const [draft, setDraft] = useState<GameRules | null>(null);
  const rules = draft ?? event.rules;
  const { isPending, run } = useSubmitOnce();

  const isFinished = Boolean(event.finished_at);
  const isReadOnly = isFinished || Boolean(event.rules_frozen);

  const handleSave = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    void run(async () => {
      if (await updateEventRules(event.id, rules)) {
        showToast({ title: t('rules.saved'), variant: 'success' });
        setDraft(null);
        onSaved();
      }
    });
  };

  return (
    <>
      <Card>
        <Card.Content className="gap-4">
          <p className="text-sm text-muted">
            {isFinished
              ? t('rules.finishedReadOnly')
              : event.rules_frozen
                ? t('rules.frozen')
                : t('rules.editableHint')}
          </p>
          <Form
            id="event-rules-form"
            onSubmit={handleSave}
            className="flex flex-col gap-4"
          >
            <GameRulesFields
              value={rules}
              onChange={setDraft}
              isReadOnly={isReadOnly}
            />
          </Form>
        </Card.Content>
      </Card>
      {!isReadOnly && (
        <PendingButton
          type="submit"
          form="event-rules-form"
          variant="primary"
          isPending={isPending}
          isDisabled={!isValidRules(rules) || sameRules(rules, event.rules)}
          className="w-3/4 self-center"
        >
          {t('rules.save')}
        </PendingButton>
      )}
    </>
  );
}
