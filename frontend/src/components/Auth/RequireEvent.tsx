import { Spinner } from '@heroui/react';
import { useEffect } from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useEvent } from '../../providers/event/eventConstants';
import { showToast } from '../Toast';

export default function RequireEvent() {
  const { selectedEventId, isLoading } = useEvent();

  useEffect(() => {
    if (!isLoading && !selectedEventId) {
      showToast({
        title: 'Select an event first',
        color: 'warning',
      });
    }
  }, [isLoading, selectedEventId]);

  // The id is known from localStorage before the event loads, so render the
  // page right away instead of swapping it for a spinner while the event
  // loads: the page shows its own spinner for its data. The key remounts the
  // page when another event is selected, so it doesn't keep the old event's
  // data. If the event fails to load, EventProvider clears the id and the
  // redirect below runs.
  if (selectedEventId) {
    return <Outlet key={selectedEventId} />;
  }

  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-screen my-auto justify-center items-center">
        <Spinner label="Loading event..." size="lg" />
      </div>
    );
  }

  return <Navigate to="/events" replace />;
}
