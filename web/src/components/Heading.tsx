import { createElement } from 'react';
import type { ReactNode } from 'react';

interface HeadingProps {
  children: ReactNode;
  level?: 1 | 2 | 3 | 4 | 5 | 6;
  variant?: 'title' | 'subtitle' | 'section';
  className?: string;
}

export function Heading({ 
  children, 
  level = 1, 
  variant = 'title', 
  className = '' 
}: HeadingProps) {
  const tagName = `h${level}`;
  
  const baseClasses = 'font-extrabold tracking-widest font-propaganda';
  
  const variantClasses = {
    title: 'text-5xl md:text-6xl text-black drop-shadow-[2px_2px_0px_orange]',
    subtitle: 'text-4xl md:text-5xl text-orange-700 drop-shadow-[2px_2px_0px_black]',
    section: 'text-2xl md:text-3xl text-black drop-shadow-[1px_1px_0px_orange]'
  };
  
  return createElement(
    tagName,
    { className: `${baseClasses} ${variantClasses[variant]} ${className}` },
    children
  );
}