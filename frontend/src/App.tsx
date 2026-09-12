import { HeroUIProvider, ToastProvider } from '@heroui/react';
import { Navigate, Route, Routes, useNavigate } from 'react-router-dom';
import './App.css';
import RequireAuth from './components/Auth/RequireAuth';
import RequireEvent from './components/Auth/RequireEvent';
import CurrentGame from './components/CurrentGame/CurrentGame';
import Draft from './components/Draft/Draft';
import EventCreate from './components/Events/EventCreate';
import EventLayout from './components/Events/EventLayout';
import EventManage from './components/Events/EventManage';
import EventOverview from './components/Events/EventOverview';
import EventRating from './components/Events/EventRating';
import EventsList from './components/Events/EventsList';
import FinishedGames from './components/Events/FinishedGames';
import LiveGames from './components/Events/LiveGames';
import GameResult from './components/GameResult/GameResult';
import History from './components/History/History';
import LoginPage from './components/Login/Login';
import Landing from './components/Main/Landing';
import ProfilePage from './components/Profile/Profile';
import ViewGame from './components/ViewGame/ViewGame';
import AssistantLayout from './layouts/AssistantLayout';
import EventPagesLayout from './layouts/EventPagesLayout';
import RootLayout from './layouts/RootLayout';
import NavigateSetter from './providers/auth/NavigateSetter';

function App() {
  const navigate = useNavigate();

  return (
    <HeroUIProvider navigate={navigate}>
      <ToastProvider
        placement="bottom-left"
        toastProps={{
          timeout: 3000,
          shouldShowTimeoutProgress: true,
        }}
        regionProps={{
          className: 'z-200',
        }}
      />
      <NavigateSetter />
      <Routes>
        <Route element={<RootLayout />}>
          <Route element={<AssistantLayout />}>
            <Route path="/" element={<Landing />} />
            <Route path="/login" element={<LoginPage />} />
            <Route
              path="/signup"
              element={<Navigate to="/login?tab=signup" replace />}
            />

            <Route element={<RequireAuth />}>
              <Route path="/profile" element={<ProfilePage />} />

              <Route element={<RequireEvent />}>
                <Route path="/now" element={<CurrentGame />} />
                <Route path="/last" element={<GameResult />} />
                <Route path="/new" element={<Draft />} />
                <Route path="/history" element={<History />} />
              </Route>
            </Route>

            <Route
              path="/event/:eventId/game/:gameId"
              element={<ViewGame />}
            />
          </Route>

          <Route element={<EventPagesLayout />}>
            <Route path="/events">
              <Route index element={<EventsList />} />
              <Route element={<RequireAuth />}>
                <Route path="new" element={<EventCreate />} />
              </Route>
            </Route>

            <Route path="/event/:eventId" element={<EventLayout />}>
              <Route index element={<EventOverview />} />
              <Route path="manage" element={<EventManage />} />
              <Route path="rating" element={<EventRating />} />
              <Route path="games/live" element={<LiveGames />} />
              <Route path="games/finished" element={<FinishedGames />} />
            </Route>
          </Route>
        </Route>
      </Routes>
    </HeroUIProvider>
  );
}

export default App;
