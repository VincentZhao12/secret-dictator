import type { ButtonHTMLAttributes, ReactNode } from "react";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: ReactNode;
  variant?: "primary" | "secondary";
  size?: "sm" | "md" | "lg";
  loading?: boolean;
}

export function Button({
  children,
  variant = "primary",
  size = "md",
  className = "",
  loading = false,
  disabled,
  ...props
}: ButtonProps) {
  const baseClasses =
    "font-extrabold uppercase tracking-wider transition-all duration-200 transform border-4 border-black font-propaganda cursor-pointer relative";
  
  const isDisabled = disabled || loading;
  const hoverClasses = isDisabled ? "" : "hover:scale-[1.03]";
  const disabledClasses = isDisabled ? "opacity-70 cursor-not-allowed" : "";

  const variantClasses = {
    primary:
      "bg-orange-600 hover:bg-orange-700 text-black shadow-[4px_4px_0px_black]",
    secondary:
      "bg-cream hover:bg-orange-300 text-black shadow-[4px_4px_0px_black]",
  };

  const sizeClasses = {
    sm: "py-2 px-4 text-sm rounded-sm",
    md: "py-4 px-6 text-base rounded-sm",
    lg: "py-5 px-8 text-lg rounded-sm",
  };

  return (
    <button
      className={`${baseClasses} ${hoverClasses} ${disabledClasses} ${variantClasses[variant]} ${sizeClasses[size]} ${className}`}
      disabled={isDisabled}
      {...props}
    >
      {loading && (
        <div className="absolute inset-0 flex items-center justify-center">
          <div className="w-5 h-5 border-2 border-black border-t-transparent rounded-full animate-spin"></div>
        </div>
      )}
      <span className={loading ? "invisible" : "visible"}>
        {children}
      </span>
    </button>
  );
}
