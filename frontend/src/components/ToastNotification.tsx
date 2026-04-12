import { useEffect } from 'react';
import './ToastNotification.css';

export interface ToastProps {
  /** 'success' | 'error' | 'info' */
  type: 'success' | 'error' | 'info';
  /** Short title displayed in bold at the top of the toast */
  title: string;
  /** Detailed message text */
  message: string;
  /** Called when the user dismisses the toast */
  onClose: () => void;
  /** Auto-dismiss after this many ms (0 = no auto-dismiss). Default: 0 */
  autoClose?: number;
}

const ICONS: Record<ToastProps['type'], string> = {
  success: '✓',
  error:   '✕',
  info:    'ℹ',
};

const TITLES: Record<ToastProps['type'], string> = {
  success: 'Success',
  error:   'Error',
  info:    'Notice',
};

export default function ToastNotification({
  type,
  title,
  message,
  onClose,
  autoClose = 0,
}: ToastProps) {

  // Optional auto-dismiss
  useEffect(() => {
    if (autoClose <= 0) return;
    const timer = setTimeout(onClose, autoClose);
    return () => clearTimeout(timer);
  }, [autoClose, onClose]);

  // Close on backdrop click
  const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (e.target === e.currentTarget) onClose();
  };

  return (
    <div className="tn-backdrop" onClick={handleBackdropClick} role="dialog" aria-modal>

      <div className={`tn-card tn-card--${type}`}>

        {/* Close × button */}
        <button className="tn-close" onClick={onClose} aria-label="Close">✕</button>

        {/* Icon circle */}
        <div className={`tn-icon-wrap tn-icon-wrap--${type}`}>
          <span className="tn-icon">{ICONS[type]}</span>
        </div>

        {/* Text */}
        <h3 className="tn-title">{title || TITLES[type]}</h3>
        <p  className="tn-message">{message}</p>

        {/* OK button */}
        <button className={`tn-ok tn-ok--${type}`} onClick={onClose}>
          OK
        </button>

      </div>
    </div>
  );
}
