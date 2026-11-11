import react from '@vitejs/plugin-react';
import { defineConfig, loadEnv } from 'vite';
import { compression } from 'vite-plugin-compression2';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig(({ mode }) => {
  // The backend the dev server forwards /api and /auth to, as nginx does on
  // prod, so the app and the API share one origin (no CORS, one port to open).
  // MACRA_API_PROXY (environment or .env.local) overrides it, e.g. for a
  // second backend on another port.
  const apiTarget =
    loadEnv(mode, '.', 'MACRA_').MACRA_API_PROXY || 'http://127.0.0.1:3001';

  return {
    server: {
      port: 3000,
      proxy: {
        '/api': apiTarget,
        '/auth': apiTarget,
      },
    },
    plugins: [react(), compression(), tailwindcss()],
    build: {
      chunkSizeWarningLimit: 1024,
      rolldownOptions: {
        output: {
          codeSplitting: {
            groups: [
              { name: 'heroui', test: /node_modules[\\/]@heroui/, priority: 2 },
              { name: 'vendor', test: /node_modules/, priority: 1 },
            ],
          },
        },
      },
    },
  };
});
