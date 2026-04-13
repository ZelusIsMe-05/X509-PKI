import { useState } from 'react';
import './DashboardPage.css';

interface DashboardPageProps {
  username: string;
  message?: string;
  onLogout: () => void;
}

export default function DashboardPage({ username, message, onLogout }: DashboardPageProps) {
  const [isExiting, setIsExiting] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(true);

  const handleLogout = () => {
    setIsExiting(true);
    setTimeout(() => {
      onLogout();
    }, 400); // Matches .fade-out CSS transition
  };

  const cards = [
    { title: 'Change Password', icon: '🔑', desc: 'Đổi mật khẩu hệ thống.' },
    { title: 'Generate Keys', icon: '🔐', desc: 'Phát sinh Các cặp khoá Public key / Private key cho cá nhân.' },
    { title: 'Request X.509 Cert (CSR)', icon: '📜', desc: 'Yêu cầu cấp phát Chứng nhận X.509 cho Website từ 1 cặp khoá của cá nhân.' },
    { title: 'My Certificates', icon: '🗂️', desc: 'Xem danh sách yêu cầu cấp phát chứng nhận và các chứng nhận đã được cấp.' },
    { title: 'Revoke Certificate', icon: '🛑', desc: 'Yêu cầu Thu hồi chứng nhận X.509 đã được cấp phát cho cá nhân.' },
    { title: 'Global CRL', icon: '🌍', desc: 'Tra cứu danh sách thu hồi chứng nhận của toàn hệ thống.' },
    { title: 'Upload external Cert', icon: '📤', desc: 'Upload các chứng nhận khác để theo dõi và xem thông tin.' },
  ];

  return (
    <div className={`client-body ${isExiting ? 'fade-out' : ''}`}>
      {/* Sidebar */}
      <div className={`client-sidebar ${sidebarOpen ? 'open' : 'closed'}`}>
        <div className="client-logo">
          {sidebarOpen ? '🌐 Client Panel' : '🌐'}
        </div>
        <div className="client-nav">
          <div className="client-nav-item active" title="Dashboard">
            <span>⊞</span>
            {sidebarOpen && <p>Dashboard</p>}
          </div>
          <div className="client-nav-item" title="Profile">
            <span>👤</span>
            {sidebarOpen && <p>Profile</p>}
          </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="client-content-wrapper">
        <div className="client-topbar">
          <button className="sidebar-toggle-btn" onClick={() => setSidebarOpen(!sidebarOpen)} title="Toggle Sidebar">
            ☰
          </button>
          <button className="client-logout-btn" onClick={handleLogout} title="Logout">
            🚪 {sidebarOpen && 'Logout'}
          </button>
        </div>

        <div className="client-main">
          <div className="client-header">
            <h1>Welcome, {username}</h1>
            <p>X.509 PKI Client Dashboard</p>
            {message && <div className="client-message">{message}</div>}
          </div>
          
          <div className="client-grid">
            {cards.map((c, idx) => (
              <div className="client-card" key={idx}>
                <div className="client-card-icon">{c.icon}</div>
                <h3 className="client-card-title">{c.title}</h3>
                <p className="client-card-desc">{c.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
