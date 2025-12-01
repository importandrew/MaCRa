import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from './authConstants';
import { AuthCreateContext } from './authTypes';

const NavigateSetter = () => {
  const navigate = useNavigate();
  const { navigateRef }: AuthCreateContext = useAuth();

  useEffect(() => {
    navigateRef.current = navigate;
  }, [navigate, navigateRef]);

  return null;
};

export default NavigateSetter;
