import type { ReactNode } from 'react';
import { X } from 'lucide-react';
import '../App.css';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: ReactNode;
  onSubmit?: () => void;
  submitLabel?: string;
  isLoading?: boolean;
}

const Modal = ({ 
  isOpen, 
  onClose, 
  title, 
  children, 
  onSubmit, 
  submitLabel = 'Save',
  isLoading = false 
}: ModalProps) => {
  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <h3>{title}</h3>
            <button onClick={onClose} style={{ background: 'none', border: 'none', cursor: 'pointer' }}>
              <X size={20} />
            </button>
          </div>
        </div>
        <div className="modal-body">
          {children}
        </div>
        {onSubmit && (
          <div className="modal-footer">
            <button className="btn btn-secondary" onClick={onClose} disabled={isLoading}>
              Cancel
            </button>
            <button className="btn btn-primary" onClick={onSubmit} disabled={isLoading}>
              {isLoading ? 'Loading...' : submitLabel}
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default Modal;
