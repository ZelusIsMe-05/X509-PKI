import { useState } from 'react';
import './LoginPage.css';
import { registerUser, loginUser } from '../api/auth';
import ToastNotification, { type ToastProps } from '../components/ToastNotification';

// 50 animated dot spans for the background ring
const SPANS = Array.from({ length: 50 }, (_, i) => i + 1);

type Mode = 'login' | 'register';

interface LoginPageProps {
  /** Called with the username when authentication succeeds. */
  onLoginSuccess: (username: string) => void;
}

interface ToastState {
  type: ToastProps['type'];
  title: string;
  message: string;
  onCloseCallback?: () => void;
}

const EyeIcon = ({ size = 20 }: { size?: number }) => (
  <svg xmlns="http://www.w3.org/2000/svg" width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
    <circle cx="12" cy="12" r="3" />
  </svg>
);

const EyeOffIcon = ({ size = 20 }: { size?: number }) => (
  <svg xmlns="http://www.w3.org/2000/svg" width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
    <path d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68" />
    <path d="M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61" />
    <line x1="2" x2="22" y1="2" y2="22" />
  </svg>
);

export default function LoginPage({ onLoginSuccess }: LoginPageProps) {
  const [mode, setMode]         = useState<Mode>('login');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState(''); // Added confirm password
  const [showPassword, setShowPassword] = useState(false); // Eye toggle state
  const [showConfirmPassword, setShowConfirmPassword] = useState(false); // Eye toggle for confirm
  const [loading, setLoading]   = useState(false);
  const [toast, setToast]       = useState<ToastState | null>(null);

  // Thêm state này để chặn trình duyệt tự động điền lúc mới tải trang
  const [isReadOnly, setIsReadOnly] = useState(true);
  const [isExiting, setIsExiting] = useState(false);

  // ── helpers ──────────────────────────────────────────────────
  const showToast = (type: ToastProps['type'], title: string, message: string, onCloseCallback?: () => void) =>
    setToast({ type, title, message, onCloseCallback });

  const closeToast = () => {
    if (toast?.onCloseCallback) {
      toast.onCloseCallback();
    }
    setToast(null);
  };

  const switchMode = (m: Mode) => {
    setMode(m);
    setToast(null);
    setUsername('');
    setPassword('');
    setConfirmPassword('');
    setShowPassword(false);
    setShowConfirmPassword(false);
    setIsReadOnly(true); // Khôi phục trạng thái chặn autofill khi chuyển mode
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') mode === 'login' ? handleLogin() : handleRegister();
  };

  // ── actions ───────────────────────────────────────────────────
  const handleLogin = async () => {
    if (!username || !password) {
      showToast('error', 'Missing Fields', 'Please enter both a username and password.');
      return;
    }
    setLoading(true);
    try {
      const data = await loginUser(username, password);
      setIsExiting(true);
      setTimeout(() => {
        onLoginSuccess(data.username);
      }, 400);
    } catch (err: any) {
      const raw: string = err.message ?? '';
      const friendly = raw.includes('invalid username or password')
        ? 'The username or password you entered is incorrect.'
        : raw.includes('failed to')
        ? 'A server error occurred. Please try again.'
        : raw || 'Unable to connect to the server.';
      showToast('error', 'Login Failed', friendly);
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async () => {
    if (!username || !password || !confirmPassword) {
      showToast('error', 'Missing Fields', 'Please fill in all fields (username, password, confirm password).');
      return;
    }
    if (password !== confirmPassword) {
      showToast('error', 'Password Mismatch', 'Your passwords do not match. Please try again.');
      return;
    }
    if (password.length < 6) {
      showToast('error', 'Weak Password', 'Password must be at least 6 characters long.');
      return;
    }
    setLoading(true);
    try {
      await registerUser(username, password);
      showToast('success', 'Account Created', 'Registration successful! You can now log in.', () => {
        switchMode('login');
      });
    } catch (err: any) {
      const raw: string = err.message ?? '';
      const friendly = raw.includes('already exists')
        ? `The username "${username}" is already taken. Please choose another.`
        : raw || 'Unable to connect to the server.';
      showToast('error', 'Registration Failed', friendly);
    } finally {
      setLoading(false);
    }
  };

  // ── render ────────────────────────────────────────────────────
  return (
    <div className={`lp-body ${isExiting ? 'fade-out' : ''}`}>

      {/* Toast notification modal */}
      {toast && (
        <ToastNotification
          type={toast.type}
          title={toast.title}
          message={toast.message}
          onClose={closeToast}
        />
      )}

      <div className={mode === 'login' ? 'lp-container' : ''} style={mode === 'register' ? { display: 'flex', justifyContent: 'center', alignItems: 'center' } : {}}>

        {/* Animated spinning-dot ring only visible in login mode */}
        {mode === 'login' && SPANS.map((i) => (
          <span key={i} style={{ '--i': i } as React.CSSProperties} />
        ))}

        {/* Auth card */}
        <div className={mode === 'login' ? 'lp-login-box' : 'lp-login-box-relative'}>
          {/* Đổi div thành form để trình duyệt nhận diện chuẩn lưu mật khẩu */}
          <form className="lp-form" onSubmit={(e) => e.preventDefault()}>

            <h2 className="lp-title">{mode === 'login' ? 'Login' : 'Register'}</h2>

            {/* Username */}
            <div className="lp-input-box">
              <input
                id="lp-username"
                type="text"
                placeholder=" "
                value={username}
                readOnly={isReadOnly} // Bật readOnly lúc đầu
                onFocus={() => setIsReadOnly(false)} // Tắt readOnly khi click vào
                onChange={(e) => setUsername(e.target.value)}
                onKeyDown={handleKeyDown}
                autoComplete="username"
              />
              <label htmlFor="lp-username">Username</label>
            </div>

            {/* Password */}
            <div className="lp-input-box">
              <input
                id="lp-password"
                type={showPassword ? "text" : "password"}
                className={showPassword ? "lp-password-input" : ""}
                placeholder=" "
                value={password}
                readOnly={isReadOnly} // Bật readOnly lúc đầu
                onFocus={() => setIsReadOnly(false)} // Tắt readOnly khi click vào
                onChange={(e) => setPassword(e.target.value)}
                onKeyDown={handleKeyDown}
                autoComplete={mode === 'login' ? 'current-password' : 'new-password'}
              />
              <label htmlFor="lp-password">Password</label>
              <button
                type="button"
                className="lp-password-toggle"
                onClick={() => setShowPassword(!showPassword)}
                tabIndex={-1}
                title={showPassword ? "Hide password" : "Show password"}
              >
                {showPassword ? <EyeOffIcon size={20} /> : <EyeIcon size={20} />}
              </button>
            </div>

            {/* Confirm Password (only visible during registration) */}
            {mode === 'register' && (
              <div className="lp-input-box">
                <input
                  id="lp-confirm-password"
                  type={showConfirmPassword ? "text" : "password"}
                  className={showConfirmPassword ? "lp-password-input" : ""}
                  placeholder=" "
                  value={confirmPassword}
                  readOnly={isReadOnly}
                  onFocus={() => setIsReadOnly(false)}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  onKeyDown={handleKeyDown}
                  autoComplete="new-password"
                />
                <label htmlFor="lp-confirm-password">Confirm Password</label>
                <button
                  type="button"
                  className="lp-password-toggle"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  tabIndex={-1}
                  title={showConfirmPassword ? "Hide password" : "Show password"}
                >
                  {showConfirmPassword ? <EyeOffIcon size={20} /> : <EyeIcon size={20} />}
                </button>
              </div>
            )}

            {/* Forgot password (login mode only) */}
            {mode === 'login' && (
              <div className="lp-forgot-pass">
                <a href="#" onClick={(e) => e.preventDefault()}>Forgot your password?</a>
              </div>
            )}

            {/* Submit */}
            <button
              type="button"
              className="lp-btn"
              onClick={mode === 'login' ? handleLogin : handleRegister}
              disabled={loading}
            >
              {loading ? 'Please wait…' : (mode === 'login' ? 'Login' : 'Register')}
            </button>

            {/* Mode switch */}
            <div className="lp-link-row">
              {mode === 'login' ? (
                <>Don't have an account?&nbsp;
                  <button type="button" onClick={() => switchMode('register')}>Signup</button>
                </>
              ) : (
                <>Already have an account?&nbsp;
                  <button type="button" onClick={() => switchMode('login')}>Login</button>
                </>
              )}
            </div>

          </form>
        </div>
      </div>
    </div>
  );
}