import type { ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  className?: string;
  size?: 'sm' | 'md' | 'lg';
}

export function Card({ children, className = '', size = 'md' }: CardProps) {
  const baseClasses = 'bg-orange-200/90 rounded-xl shadow-[6px_6px_0px_black] border-4 border-black relative overflow-hidden';
  
  const sizeClasses = {
    sm: 'p-6',
    md: 'p-10 md:p-14',
    lg: 'p-12 md:p-16'
  };
  
  return (
    <div className={`${baseClasses} ${sizeClasses[size]} ${className}`}>
      {children}
    </div>
  );
}