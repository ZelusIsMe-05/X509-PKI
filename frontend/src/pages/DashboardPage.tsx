import './DashboardPage.css';

interface DashboardPageProps {
  /** The username of the currently authenticated user. */
  username: string;
  /** Optional status message to display (e.g. "session restored"). */
  message?: string;
  /** Called when the user clicks the Logout button. */
  onLogout: () => void;
}

export default function DashboardPage({ username, message, onLogout }: DashboardPageProps) {
  const accessToken = localStorage.getItem('auth_access_token') ?? '';

  return (
    <div className="dp-body">

      {/* Decorative corner glow */}
      <div className="dp-glow dp-glow--tl" aria-hidden />
      <div className="dp-glow dp-glow--br" aria-hidden />

      <div className="dp-card">

        {/* ── Header ─────────────────────────────────── */}
        <div className="dp-header">
          <span className="dp-badge">JWT Auth</span>
          <h1 className="dp-title">X.509 PKI Dashboard</h1>
        </div>

        {/* ── User identity banner ────────────────────── */}
        <div className="dp-user-banner">
          <div className="dp-avatar" aria-label={`Avatar for ${username}`}>
            {username.charAt(0).toUpperCase()}
          </div>
          <div className="dp-user-info">
            <span className="dp-user-label">Signed in as</span>
            <span className="dp-username">{username}</span>
          </div>
          <div className="dp-status-dot" title="Session active" />
        </div>

        {/* ── Access token preview ────────────────────── */}
        <div className="dp-token-card">
          <div className="dp-token-header">
            <span className="dp-token-icon">🎫</span>
            <span className="dp-token-label">Access Token (JWT)</span>
          </div>
          <code className="dp-token-value">
            {accessToken ? `${accessToken.slice(0, 64)}…` : '—'}
          </code>
        </div>

        {/* ── Session status message ──────────────────── */}
        {message && (
          <div className="dp-message">
            <span>{message}</span>
          </div>
        )}

        {/* ── Logout ─────────────────────────────────── */}
        <button
          id="logout-btn"
          className="dp-logout-btn"
          onClick={onLogout}
        >
          🚪 Logout
        </button>

      </div>
    </div>
  );
}
