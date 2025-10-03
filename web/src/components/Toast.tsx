import { useEffect } from "react";
import { FaCheckCircle, FaExclamationCircle, FaTimes } from "react-icons/fa";

export type ToastType = "success" | "error" | "info";

export interface Toast {
  id: string;
  message: string;
  type: ToastType;
  duration?: number;
}

interface ToastProps {
  toast: Toast;
  onClose: (id: string) => void;
}

export function Toast({ toast, onClose }: ToastProps) {
  useEffect(() => {
    const duration = toast.duration ?? 2000;
    const timer = setTimeout(() => {
      onClose(toast.id);
    }, duration);

    return () => clearTimeout(timer);
  }, [toast.id, toast.duration, onClose]);

  const getToastStyles = () => {
    switch (toast.type) {
      case "success":
        return "bg-green-500 border-green-700";
      case "error":
        return "bg-red-500 border-red-700";
      case "info":
        return "bg-blue-500 border-blue-700";
      default:
        return "bg-gray-500 border-gray-700";
    }
  };

  const getIcon = () => {
    switch (toast.type) {
      case "success":
        return <FaCheckCircle className="text-white text-xl" />;
      case "error":
        return <FaExclamationCircle className="text-white text-xl" />;
      case "info":
        return <FaExclamationCircle className="text-white text-xl" />;
      default:
        return null;
    }
  };

  return (
    <div
      className={`${getToastStyles()} text-white px-5 py-4 rounded-lg border-4 border-black shadow-[4px_4px_0px_black] font-propaganda tracking-wider transition-all duration-300 transform hover:scale-105 min-w-[300px] max-w-[500px]`}
    >
      <div className="flex items-center justify-between space-x-3">
        <div className="flex items-center space-x-3 flex-1">
          {getIcon()}
          <p className="font-bold text-sm break-words">{toast.message}</p>
        </div>
        <button
          onClick={() => onClose(toast.id)}
          className="text-white hover:text-gray-200 transition-colors flex-shrink-0"
          aria-label="Close toast"
        >
          <FaTimes className="text-lg" />
        </button>
      </div>
    </div>
  );
}
