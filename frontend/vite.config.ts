import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';
import { compression } from 'vite-plugin-compression2';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig(() => {
  return {
    server: {
      port: 3000,
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
