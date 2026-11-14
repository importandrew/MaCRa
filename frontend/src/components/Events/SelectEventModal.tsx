import { useTranslation } from 'react-i18next';
import { Modal } from '@heroui/react';
import { useEffect, useMemo, useState } from 'react';
import type { Event } from '../../models/Event';
import { useNavigate } from 'react-router-dom';
import { useFetchVisibleEvents } from '../../providers/event/eventApi';
import { useEvent } from '../../providers/event/eventConstants';
import LoadingSpinner from '../LoadingTemplate/LoadingSpinner';
import EventSearchField from './EventSearchField';
import FinishedEvents from './FinishedEvents';
import { matchesQuery, splitFinished } from './eventListFilters';

// SEARCH_FROM: the picker shows its search field from this many events on.
const SEARCH_FROM = 6;

type SelectEventModalProps = {
  isOpen: boolean;
  onClose: () => void;
};

export default function SelectEventModal({
  isOpen,
  onClose,
}: SelectEventModalProps) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { selectedEventId, selectEvent } = useEvent();
  const { data, isLoading, isError, refetch } = useFetchVisibleEvents();

  useEffect(() => {
    if (isOpen) {
      refetch();
    }
  }, [isOpen, refetch]);

  const participatingEvents = useMemo(
    () => data?.filter((e) => e.is_participating) ?? [],
    [data]
  );
  // Active events on top (by their last game, as the server orders them),
  // finished ones collapsed under them: still selectable, for the last game
  // and history.
  const [query, setQuery] = useState('');
  const isSearching = query.trim() !== '';
  const found = participatingEvents.filter((e) => matchesQuery(e, query));
  const { active, finished } = splitFinished(found);

  const handleSelect = async (id: number) => {
    await selectEvent(id);
    onClose();
    navigate('/');
  };

  const renderEvent = (event: Event) => {
    const isSelected = selectedEventId === event.id;
    const isFinished = Boolean(event.finished_at);
    return (
      <button
        key={event.id}
        type="button"
        onClick={() => handleSelect(event.id)}
        className={`w-full text-left rounded-lg border px-4 py-3 transition-colors hover:bg-default ${
          isSelected ? 'border-accent bg-accent-soft' : 'border-border'
        }`}
      >
        <div
          className={`font-medium break-words ${isFinished ? 'text-muted' : 'text-foreground'}`}
        >
          {event.name}
        </div>
        {event.description && (
          <p className="text-sm text-muted mt-1 line-clamp-3 break-words">
            {event.description}
          </p>
        )}
      </button>
    );
  };

  return (
    <Modal.Backdrop
      isOpen={isOpen}
      onOpenChange={(open) => {
        if (!open) onClose();
      }}
    >
      <Modal.Container placement="center" size="lg" scroll="inside">
        <Modal.Dialog>
          <Modal.CloseTrigger aria-label={t('common.close')} />
          <Modal.Header>
            <Modal.Heading>{t('events.selectEvent')}</Modal.Heading>
          </Modal.Header>
          <Modal.Body className="flex flex-col gap-3 pb-6">
            {isLoading && (
              <div className="flex justify-center py-8">
                <LoadingSpinner label={t('events.loadingEvents')} size="md" />
              </div>
            )}
            {!isLoading && isError && (
              <p className="text-center text-danger py-6">
                {t('events.loadFailed')}
              </p>
            )}
            {!isLoading && !isError && participatingEvents.length === 0 && (
              <p className="text-center text-muted py-6">
                {t('events.notParticipating')}
              </p>
            )}
            {!isLoading &&
              !isError &&
              participatingEvents.length >= SEARCH_FROM && (
                <EventSearchField value={query} onChange={setQuery} />
              )}
            {!isLoading && !isError && isSearching && found.length === 0 && (
              <p className="text-center text-muted py-6">
                {t('events.nothingFound', { query: query.trim() })}
              </p>
            )}
            {!isLoading && !isError && (
              <>
                {active.map((event) => renderEvent(event))}
                <FinishedEvents
                  count={finished.length}
                  isSearching={isSearching}
                >
                  {finished.map((event) => renderEvent(event))}
                </FinishedEvents>
              </>
            )}
          </Modal.Body>
        </Modal.Dialog>
      </Modal.Container>
    </Modal.Backdrop>
  );
}
