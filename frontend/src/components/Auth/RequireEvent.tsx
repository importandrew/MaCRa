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

  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-screen my-auto justify-center items-center">
        <Spinner label="Loading event..." size="lg" />
      </div>
    );
  }

  if (!selectedEventId) {
    return <Navigate to="/events" replace />;
  }

  return <Outlet />;
}
