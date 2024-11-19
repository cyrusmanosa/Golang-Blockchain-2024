import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import fs from 'fs';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    https: {
      key: fs.readFileSync("./ca-key.pem"),
      cert: fs.readFileSync("./ca-cert.pem"),
    },
  },
});
