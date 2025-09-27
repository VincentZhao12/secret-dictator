import type { ReactNode } from 'react';

interface BadgeProps {
  children: ReactNode;
  variant?: 'primary' | 'secondary' | 'warning' | 'danger';
  size?: 'sm' | 'md';
  className?: string;
}

export function Badge({ 
  children, 
  variant = 'primary', 
  size = 'md', 
  className = '' 
}: BadgeProps) {
  const baseClasses = 'inline-flex items-center font-extrabold uppercase tracking-wider border-2 border-black font-propaganda';
  
  const variantClasses = {
    primary: 'bg-orange-600 text-black',
    secondary: 'bg-cream text-black',
    warning: 'bg-yellow-400 text-black',
    danger: 'bg-red-500 text-white'
  };
  
  const sizeClasses = {
    sm: 'px-2 py-1 text-xs rounded',
    md: 'px-3 py-1 text-sm rounded'
  };
  
  return (
    <span className={`${baseClasses} ${variantClasses[variant]} ${sizeClasses[size]} ${className}`}>
      {children}
    </span>
  );
}