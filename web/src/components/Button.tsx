import type { ButtonHTMLAttributes, ReactNode } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: ReactNode;
  variant?: 'primary' | 'secondary';
  size?: 'sm' | 'md' | 'lg';
}

export function Button({ 
  children, 
  variant = 'primary', 
  size = 'md', 
  className = '', 
  ...props 
}: ButtonProps) {
  const baseClasses = 'font-extrabold uppercase tracking-wider transition-all duration-200 transform hover:scale-[1.03] border-4 border-black font-propaganda';
  
  const variantClasses = {
    primary: 'bg-orange-600 hover:bg-orange-700 text-black shadow-[4px_4px_0px_black]',
    secondary: 'bg-cream hover:bg-orange-300 text-black shadow-[4px_4px_0px_black]'
  };
  
  const sizeClasses = {
    sm: 'py-2 px-4 text-sm rounded-sm',
    md: 'py-4 px-6 text-base rounded-sm',
    lg: 'py-5 px-8 text-lg rounded-sm'
  };
  
  return (
    <button
      className={`${baseClasses} ${variantClasses[variant]} ${sizeClasses[size]} ${className}`}
      {...props}
    >
      {children}
    </button>
  );
}