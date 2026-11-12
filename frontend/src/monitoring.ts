import HawkCatcher from '@hawk.so/javascript';

// hawk sends errors to Hawk in production builds that have a token; in
// development they stay in the console.
const hawk =
  !import.meta.env.DEV && import.meta.env.VITE_HAWK_TOKEN
    ? new HawkCatcher({
        token: import.meta.env.VITE_HAWK_TOKEN,
        consoleTracking: true,
        context: {},
        beforeSend: (event) => {
          const titleIgnorePatterns = ['Script error.'];

          if (titleIgnorePatterns.includes(event.title.trim())) return false;

          return event;
        },
      })
    : null;

// reportError sends an error the app has caught itself (e.g. by the error
// boundary): Hawk's global handlers see only uncaught ones.
export function reportError(error: unknown) {
  const err = error instanceof Error ? error : new Error(String(error));
  hawk?.captureError(err);
}
