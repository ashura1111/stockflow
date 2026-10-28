import { useEffect } from 'react';
import '../App.css';

interface ToastProps {
  message: string;
  type: 'success' | 'error' | 'info';
  onClose: () => void;
  duration?: number;
}

const Toast = ({ message, type, onClose, duration = 3000 }: ToastProps) => {
  useEffect(() => {
    const timer = setTimeout(onClose, duration);
    return () => clearTimeout(timer);
  }, [onClose, duration]);

  const bgColor = {
    success: '#dcfce7',
    error: '#fee2e2',
    info: '#dbeafe',
  }[type];

  const textColor = {
    success: '#16a34a',
    error: '#dc2626',
    info: '#2563eb',
  }[type];

  return (
    <div
      style={{
        position: 'fixed',
        bottom: 20,
        right: 20,
        backgroundColor: bgColor,
        color: textColor,
        padding: '12px 16px',
        borderRadius: '6px',
        boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
        zIndex: 2000,
        maxWidth: '300px',
      }}
    >
      {message}
    </div>
  );
};

export default Toast;
