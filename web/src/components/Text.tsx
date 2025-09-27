import type { ReactNode } from 'react';

interface TextProps {
  children: ReactNode;
  variant?: 'body' | 'description' | 'caption' | 'footer';
  className?: string;
  as?: 'p' | 'span' | 'div';
}

export function Text({ 
  children, 
  variant = 'body', 
  className = '', 
  as = 'p' 
}: TextProps) {
  const Component = as;
  
  const variantClasses = {
    body: 'text-black/90 leading-relaxed',
    description: 'text-black/90 italic leading-relaxed',
    caption: 'text-black/80 text-sm',
    footer: 'text-black/80 text-xs uppercase tracking-widest'
  };
  
  return (
    <Component className={`${variantClasses[variant]} ${className}`}>
      {children}
    </Component>
  );
}