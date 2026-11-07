import { useCallback, useRef, useState } from 'react';

// useSubmitOnce runs an action that sends a mutating request, ignoring new
// calls while one is in flight. isPending alone isn't enough: its re-render
// lands too late for a double tap (two POST /new 90 ms apart reached
// production), so a ref blocks the second call at once. Pass isPending to the
// button for the pending look.
export function useSubmitOnce() {
  const busyRef = useRef(false);
  const [isPending, setIsPending] = useState(false);

  const run = useCallback(
    async <T>(action: () => Promise<T>): Promise<T | undefined> => {
      if (busyRef.current) return undefined;

      busyRef.current = true;
      setIsPending(true);
      try {
        return await action();
      } finally {
        busyRef.current = false;
        setIsPending(false);
      }
    },
    []
  );

  return { isPending, run };
}
