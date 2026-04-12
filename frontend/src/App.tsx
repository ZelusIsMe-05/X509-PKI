import { useState, useEffect } from 'react';
import { verifySession, getStoredUsername, clearTokens } from './api/auth';
import LoadingPage   from './pages/LoadingPage';
import LoginPage     from './pages/LoginPage';
import DashboardPage from './pages/DashboardPage';

// App-level state machine
type AppState = 'loading' | 'login' | 'dashboard';

/**
 * App is a thin orchestrator that owns the top-level authentication state and
 * decides which page to render. All UI logic lives in the individual page files.
 *
 *  loading   → check localStorage for a saved session
 *  login     → show LoginPage (login + register)
 *  dashboard → show DashboardPage (user info + logout)
 */ 
export default function App() {
  const [appState, setAppState]     = useState<AppState>('loading');
  const [loggedInUser, setLoggedInUser] = useState('');
  const [sessionMsg, setSessionMsg]  = useState('');

  // On mount: attempt to restore a previous session silently
  useEffect(() => {
    const restoreSession = async () => {
      const saved = getStoredUsername();
      if (saved) {
        const verified = await verifySession();
        if (verified) {
          setLoggedInUser(verified);
          setSessionMsg('✅ Previous session restored automatically');
          setAppState('dashboard');
          return;
        }
      }
      setAppState('login');
    };

    restoreSession();
  }, []);

  const handleLoginSuccess = (username: string) => {
    setLoggedInUser(username);
    setSessionMsg('');
    setAppState('dashboard');
  };

  const handleLogout = () => {
    clearTokens();
    setLoggedInUser('');
    setSessionMsg('');
    setAppState('login');
  };

  switch (appState) {
    case 'loading':
      return <LoadingPage />;
    case 'dashboard':
      return (
        <DashboardPage
          username={loggedInUser}
          message={sessionMsg}
          onLogout={handleLogout}
        />
      );
    default:
      return <LoginPage onLoginSuccess={handleLoginSuccess} />;
  }
}