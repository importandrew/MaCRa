import React from 'react';
import ReactDOM, { Container } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App.tsx';
import './index.css';
import './i18n';
import './monitoring';
import { AuthProvider } from './providers/auth/AuthContext.tsx';
import { EventProvider } from './providers/event/EventContext.tsx';

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
