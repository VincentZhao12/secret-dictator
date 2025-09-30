import { CardFascist, CardLiberal, type Card } from "@types";

interface PolicyCardProps {
  type: Card;
  onClick?: () => void;
}

export function PolicyCard({ type, onClick }: PolicyCardProps) {
  const baseClasses =
    "relative w-16 h-20 border-4 border-black rounded-lg flex items-center justify-center";

  const typeClasses = {
    [CardLiberal]: "bg-blue-500 shadow-[3px_3px_0px_black]",
    [CardFascist]: "bg-orange-600 shadow-[3px_3px_0px_black]",
  };
  return (
    <div className={`${baseClasses} ${typeClasses[type]}`} onClick={onClick}>
      <div className="w-full h-full flex items-center justify-center">
        {type === CardLiberal ? (
          <div className="w-4 h-4 bg-blue-700 rounded-full"></div>
        ) : (
          <div className="w-6 h-6 bg-red-800 transform rotate-45"></div>
        )}
      </div>
    </div>
  );
}
