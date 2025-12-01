import { HeroUIProvider } from '@heroui/react';
import { Route, Routes, useNavigate } from 'react-router-dom';
import { Slide, ToastContainer } from 'react-toastify';
import Draft from './components/Draft/Draft';
import End from './components/End/End';
import LoginPage from './components/Login/Login';
import MainPage from './components/Main/Main';
import Now from './components/Now/Now';
import SignUpPage from './components/SignUp/SignUp';
import NavigateSetter from './providers/auth/NavigateSetter';
import 'react-toastify/dist/ReactToastify.css';
import './App.css';

function App() {
  const navigate = useNavigate();

  return (
    <HeroUIProvider navigate={navigate}>
      <NavigateSetter />
      <Routes>
        <Route path="/" element={<MainPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/signup" element={<SignUpPage />} />
        <Route path="/now" element={<Now />} />
        <Route path="/last" element={<End />} />
        <Route path="/new" element={<Draft />} />
      </Routes>
      <ToastContainer
        className="text-small"
        position="bottom-left"
        theme="colored"
        autoClose={5000}
        newestOnTop
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        pauseOnHover
        transition={Slide}
      />
    </HeroUIProvider>
  );
}

export default App;
