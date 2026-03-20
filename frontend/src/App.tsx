import { HeroUIProvider, ToastProvider } from '@heroui/react';
import { Route, Routes, useNavigate } from 'react-router-dom';
import './App.css';
import Draft from './components/Draft/Draft';
import End from './components/End/End';
import History from './components/History/History';
import Live from './components/Live/Live';
import LoginPage from './components/Login/Login';
import MainPage from './components/Main/Main';
import Now from './components/Now/Now';
import ProfilePage from './components/Profile/Profile';
import SignUpPage from './components/SignUp/SignUp';
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
        <Route path="/" element={<MainPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/signup" element={<SignUpPage />} />
        <Route path="/profile" element={<ProfilePage />} />
        <Route path="/now" element={<Now />} />
        <Route path="/last" element={<End />} />
        <Route path="/new" element={<Draft />} />
        <Route path="/history" element={<History />} />
        <Route path="/game/:id/history" element={<End />} />
        <Route path="/game/:id/live" element={<Live />} />
      </Routes>
    </HeroUIProvider>
  );
}

export default App;
