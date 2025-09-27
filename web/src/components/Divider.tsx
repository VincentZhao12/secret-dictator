interface DividerProps {
  className?: string;
  orientation?: 'horizontal' | 'vertical';
}

export function Divider({ className = '', orientation = 'horizontal' }: DividerProps) {
  const baseClasses = 'bg-black';
  const orientationClasses = {
    horizontal: 'w-full h-1',
    vertical: 'h-full w-1'
  };
  
  return (
    <div className={`${baseClasses} ${orientationClasses[orientation]} ${className}`} />
  );
}