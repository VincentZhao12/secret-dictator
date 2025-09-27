import type { ReactNode } from 'react';

interface ContainerProps {
  children: ReactNode;
  className?: string;
}

export function Container({ children, className = '' }: ContainerProps) {
  return (
    <div className={`min-h-screen bg-gradient-to-br from-orange-700 via-orange-600 to-orange-500 flex items-center justify-center p-6 font-propaganda ${className}`}>
      {children}
    </div>
  );
}