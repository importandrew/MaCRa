import React from 'react';
import ReactDOM, { Container } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import HawkCatcher from '@hawk.so/javascript';
import App from './App.tsx';
import './index.css';
import { AuthProvider } from './providers/auth/AuthContext.tsx';
import { EventProvider } from './providers/event/EventContext.tsx';

if (process.env.NODE_ENV === 'development') {
  // DEBUG mode in Console
} else {
  new HawkCatcher({
    token:
      'REDACTED_HAWK_TOKEN',
    consoleTracking: true,
    context: {},
    beforeSend: (event) => {
      const titleIgnorePatterns = ['Script error.'];

      if (titleIgnorePatterns.includes(event.title.trim())) return false;

      return event;
    },
  });
}

ReactDOM.createRoot(document.getElementById('root') as Container).render(
  <React.StrictMode>
    <AuthProvider>
      <EventProvider>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </EventProvider>
    </AuthProvider>
  </React.StrictMode>
);
