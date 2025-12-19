import React from 'react';
import ReactDOM, { Container } from 'react-dom/client';
import App from './App.tsx';
import './index.css';
import { BrowserRouter } from 'react-router-dom';
import Layout from './Layout.tsx';
import { AuthProvider } from './providers/auth/AuthContext.tsx';
import HawkCatcher from '@hawk.so/javascript';

new HawkCatcher({
  token:
    'REDACTED_HAWK_TOKEN',
  consoleTracking: true,
});

ReactDOM.createRoot(document.getElementById('root') as Container).render(
  <React.StrictMode>
    <AuthProvider>
      <Layout>
        <BrowserRouter>
          <main className="flex grow light text-foreground bg-background">
            <App />
          </main>
        </BrowserRouter>
      </Layout>
    </AuthProvider>
  </React.StrictMode>
);
