import type { InputHTMLAttributes } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export function Input({ 
  label, 
  error, 
  className = '', 
  ...props 
}: InputProps) {
  const baseClasses = 'w-full px-4 py-3 border-4 border-black rounded-sm font-propaganda font-bold text-black bg-cream focus:bg-white focus:outline-none focus:ring-0 transition-colors placeholder:text-black/60';
  
  return (
    <div className="space-y-2">
      {label && (
        <label className="block text-black font-extrabold uppercase tracking-wider font-propaganda text-sm">
          {label}
        </label>
      )}
      <input
        className={`${baseClasses} ${error ? 'border-red-500' : ''} ${className}`}
        {...props}
      />
      {error && (
        <p className="text-red-500 text-sm font-bold font-propaganda">
          {error}
        </p>
      )}
    </div>
  );
}