import type { ReactNode } from 'react';

interface AlertProps {
  children: ReactNode;
  variant?: 'error' | 'warning' | 'info' | 'success';
  onClose?: () => void;
  className?: string;
}

export function Alert({ 
  children, 
  variant = 'error', 
  onClose, 
  className = '' 
}: AlertProps) {
  const baseClasses = 'p-4 border-4 border-black rounded-sm font-propaganda font-bold relative';
  
  const variantClasses = {
    error: 'bg-red-100 text-red-800 border-red-800',
    warning: 'bg-yellow-100 text-yellow-800 border-yellow-800', 
    info: 'bg-blue-100 text-blue-800 border-blue-800',
    success: 'bg-green-100 text-green-800 border-green-800'
  };
  
  return (
    <div className={`${baseClasses} ${variantClasses[variant]} ${className}`}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          {children}
        </div>
        {onClose && (
          <button
            onClick={onClose}
            className="ml-4 text-current hover:opacity-70 transition-opacity font-extrabold text-lg leading-none"
            aria-label="Close alert"
          >
            ×
          </button>
        )}
      </div>
    </div>
  );
}