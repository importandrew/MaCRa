import { useTranslation } from 'react-i18next';
import { Button, Modal } from '@heroui/react';
import { useState } from 'react';
import type { Substitute, TourSeat, TourTable } from '../../models/Tour';
import { useSubmitOnce } from '../../hooks/useSubmitOnce';
import { penalize, seatSitter } from '../../providers/event/tournamentApi';
import PendingButton from '../PendingButton';
import SegmentedPicker from '../SegmentedPicker';
import { showToast } from '../Toast';

// LATE_PENALTIES: the official fines for being late (MCR 2021: up to 10
// minutes −10, 11–15 minutes −20); 0 is none.
const LATE_PENALTIES = ['0', '10', '20'] as const;
type LatePenalty = (typeof LATE_PENALTIES)[number];

// SeatModal is a referee's change at a seat of a tournament table (area 8,
// "Substitutes"): seat a substitute there from the next hand, choosing what
// the seat's player gets for the tour (the seat's result as usual, or 0 / 0),
// or seat the player back, with a fine for being late if the game is live.
export default function SeatModal({
  eventId,
  target,
  substitutes,
  onClose,
  onDone,
}: {
  eventId: number;
  // target: the table and seat to change; null while the modal is closed.
  target: { table: TourTable; seat: TourSeat } | null;
  substitutes: Substitute[];
  onClose: () => void;
  onDone: () => Promise<unknown>;
}) {
  const { t } = useTranslation();
  const { isPending, run } = useSubmitOnce();
  // sitter: who is to sit there (a substitute's ID, or the seat's player's
  // to seat them back); null until picked.
  const [sitter, setSitter] = useState<number | null>(null);
  const [forfeit, setForfeit] = useState(false);
  const [penalty, setPenalty] = useState<LatePenalty>('0');

  const close = () => {
    setSitter(null);
    setForfeit(false);
    setPenalty('0');
    onClose();
  };

  const seat = target?.seat;
  const table = target?.table;
  const isBack = seat != null && sitter === seat.user_id;
  const options = [
    // The player, when someone else sits there now.
    ...(seat?.sitter
      ? [
          {
            value: String(seat.user_id),
            label: t('seat.back', { name: seat.name }),
          },
        ]
      : []),
    ...substitutes
      .filter((sub) => sub.id !== seat?.sitter?.id)
      .map((sub) => ({
        value: String(sub.id),
        label:
          sub.table > 0
            ? t('seat.busy', { name: sub.name, table: sub.table })
            : sub.name,
        isDisabled: sub.table > 0,
      })),
  ];

  const save = () =>
    run(async () => {
      if (!seat || !table || sitter == null) return;
      const ok = await seatSitter(eventId, table.game_id, {
        seat: seat.user_id,
        sitter,
        forfeit: !isBack && forfeit,
        log_size: table.log_size,
      });
      if (!ok) return;
      if (isBack && penalty !== '0' && table.status === 'live') {
        await penalize(
          eventId,
          table.game_id,
          seat.user_id,
          Number(penalty),
          table.log_size
        );
      }
      showToast({
        title: isBack
          ? t('seat.backDone', { name: seat.name })
          : t('seat.seated', { table: table.table }),
        variant: 'success',
      });
      await onDone();
      close();
    });

  return (
    <Modal.Backdrop
      isOpen={target != null}
      onOpenChange={(open) => {
        if (!open) close();
      }}
    >
      <Modal.Container placement="center" scroll="inside">
        <Modal.Dialog>
          <Modal.CloseTrigger aria-label={t('common.close')} />
          <Modal.Header>
            <Modal.Heading>
              {t('seat.title', { table: table?.table ?? 0, name: seat?.name })}
            </Modal.Heading>
          </Modal.Header>
          <Modal.Body className="flex flex-col gap-3">
            {seat?.sitter && (
              <p className="text-sm">
                {t('seat.sitsNow', { name: seat.sitter.name })}
              </p>
            )}
            {options.length === 0 ? (
              <p className="text-sm text-muted">{t('seat.noSubstitutes')}</p>
            ) : (
              <SegmentedPicker
                label={t('seat.who')}
                orientation="vertical"
                options={options}
                selectedKey={sitter != null ? String(sitter) : null}
                onSelectionChange={(id) => setSitter(Number(id))}
              />
            )}
            {sitter != null && !isBack && (
              <>
                <SegmentedPicker
                  label={t('seat.outcome', { name: seat?.name })}
                  options={[
                    { value: 'usual', label: t('seat.usual') },
                    { value: 'forfeit', label: t('seat.forfeit') },
                  ]}
                  selectedKey={forfeit ? 'forfeit' : 'usual'}
                  onSelectionChange={(value) => setForfeit(value === 'forfeit')}
                />
                <p className="text-sm text-muted">
                  {forfeit ? t('seat.forfeitHint') : t('seat.usualHint')}
                </p>
              </>
            )}
            {isBack && table?.status === 'live' && (
              <SegmentedPicker
                label={t('seat.latePenalty')}
                options={LATE_PENALTIES.map((value) => ({
                  value,
                  label: value === '0' ? t('seat.noPenalty') : `−${value}`,
                }))}
                selectedKey={penalty}
                onSelectionChange={setPenalty}
              />
            )}
            <p className="text-sm text-muted">{t('seat.fromNextHand')}</p>
          </Modal.Body>
          <Modal.Footer>
            <Button slot="close" variant="tertiary">
              {t('common.cancel')}
            </Button>
            <PendingButton
              variant="primary"
              isPending={isPending}
              isDisabled={sitter == null}
              onPress={save}
            >
              {isBack ? t('seat.backButton') : t('seat.seat')}
            </PendingButton>
          </Modal.Footer>
        </Modal.Dialog>
      </Modal.Container>
    </Modal.Backdrop>
  );
}
