import { useState } from 'react';
import './AdminDashboardPage.css';

interface AdminDashboardPageProps {
  username: string;
  onLogout: () => void;
}

export default function AdminDashboardPage({ username, onLogout }: AdminDashboardPageProps) {
  const [isExiting, setIsExiting] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(true);

  const handleLogout = () => {
    setIsExiting(true);
    setTimeout(() => {
      onLogout();
    }, 400); // Matches .fade-out CSS duration
  };

  const cards = [
    { title: 'Change Password', icon: '🔑', desc: 'Đổi mật khẩu quản trị hệ thống.' },
    { title: 'System Parameters', icon: '⚙️', desc: 'Thiết lập thông số kỹ thuật chuẩn (Thuật toán, Hash, Độ dài khoá...).' },
    { title: 'Generate Root Keys', icon: '🔐', desc: 'Phát sinh cặp khoá Public / Private key cho Root Certificate.' },
    { title: 'Generate Root Cert', icon: '📜', desc: 'Phát sinh Root Certificate cho toàn hệ thống.' },
    { title: 'Reject Requests', icon: '❌', desc: 'Từ chối việc yêu cầu cấp chứng nhận X.509.' },
    { title: 'Approve Requests', icon: '✅', desc: 'Phê duyệt và phát sinh chứng nhận X.509.' },
    { title: 'Manage Certs', icon: '🗂️', desc: 'Quản lý các chứng nhận đã cấp phát (revoke, renew).' },
    { title: 'Revocation Approvals', icon: '🛑', desc: 'Quản lý phê duyệt các yêu cầu thu hồi chứng nhận.' },
    { title: 'Update CRL', icon: '🔄', desc: 'Cập nhật danh sách thu hồi chứng nhận (CRL).' },
    { title: 'System Logs', icon: '📋', desc: 'Theo dõi nhật ký quá trình hoạt động chính của hệ thống.' },
  ];

  return (
    <div className={`admin-body ${isExiting ? 'fade-out' : ''}`}>
      {/* Sidebar */}
      <div className={`admin-sidebar ${sidebarOpen ? 'open' : 'closed'}`}>
        <div className="admin-logo">
          {sidebarOpen ? '🛡️ Admin Panel' : '🛡️'}
        </div>
        <div className="admin-nav">
          <div className="admin-nav-item active" title="Dashboard">
            <span>⊞</span>
            {sidebarOpen && <p>Dashboard</p>}
          </div>
          <div className="admin-nav-item" title="Settings">
            <span>⚙</span>
            {sidebarOpen && <p>Settings</p>}
          </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="admin-content-wrapper">
        <div className="admin-topbar">
          <button className="sidebar-toggle-btn" onClick={() => setSidebarOpen(!sidebarOpen)} title="Toggle Sidebar">
            ☰
          </button>
          <button className="admin-logout-btn" onClick={handleLogout} title="Logout">
            🚪 {sidebarOpen && 'Logout'}
          </button>
        </div>

        <div className="admin-main">
          <div className="admin-header">
            <h1>Welcome, {username}</h1>
            <p>X.509 PKI Administration Dashboard</p>
          </div>
          
          <div className="admin-grid">
            {cards.map((c, idx) => (
              <div className="admin-card" key={idx}>
                <div className="admin-card-icon">{c.icon}</div>
                <h3 className="admin-card-title">{c.title}</h3>
                <p className="admin-card-desc">{c.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
