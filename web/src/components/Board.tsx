import {
  ActionExecution,
  ActionInvestigate,
  ActionPolicyPeek,
  CardFascist,
  CardLiberal,
  type Action,
  type Board as BoardType,
  type Card,
} from "../types";
import { FaMagnifyingGlass } from "react-icons/fa6";
import { RiStackFill } from "react-icons/ri";
import { GiSilverBullet } from "react-icons/gi";
import { useMemo } from "react";
import { PolicyCard } from "./PolicyCard";

interface BoardProps {
  board: BoardType;
  className?: string;
}

function ActionIcon({ action }: { action: Action }) {
  switch (action) {
    case ActionInvestigate:
      return (
        <div className="w-4 h-4 bg-black rounded-full">
          <FaMagnifyingGlass />
        </div>
      );
    case ActionPolicyPeek:
      return (
        <div className="w-4 h-4 bg-black rounded-full">
          <RiStackFill />
        </div>
      );
    case ActionExecution:
      return (
        <div className="w-4 h-4 bg-black rounded-full">
          <GiSilverBullet />
        </div>
      );
    default:
      return <div className="w-4 h-4 bg-black rounded-full"></div>;
  }
}

interface PolicySlotProps {
  isActive: boolean;
  type: Card;
  executiveAction?: Action;
}

function PolicySlot({ isActive, type, executiveAction }: PolicySlotProps) {
  const baseClasses =
    "relative w-16 h-20 border-4 border-black rounded-lg flex items-center justify-center";

  const typeClasses = {
    [CardLiberal]: isActive
      ? "bg-blue-500 shadow-[3px_3px_0px_black]"
      : "bg-blue-200/50 shadow-[2px_2px_0px_black]",
    [CardFascist]: isActive
      ? "bg-orange-600 shadow-[3px_3px_0px_black]"
      : "bg-red-200/50 shadow-[2px_2px_0px_black]",
  };

  return useMemo(
    () =>
      isActive ? (
        <PolicyCard type={type} />
      ) : (
        <div className={`${baseClasses} ${typeClasses[type]}`}>
          {isActive && <PolicyCard type={type} />}
          {/* Placeholder for executive action icon */}
          {executiveAction && !isActive && (
            <div className="absolute top-1 right-1 w-4 h-4 bg-black rounded-full opacity-30">
              <ActionIcon action={executiveAction} />
            </div>
          )}
        </div>
      ),
    [isActive, type, executiveAction]
  );
}

interface ElectionTrackerProps {
  tracker: BoardType["election_tracker"];
}

function ElectionTracker({ tracker }: ElectionTrackerProps) {
  return useMemo(
    () => (
      <div className="">
        <div className="flex justify-center space-x-2">
          {Array.from({ length: tracker.max_failures }, (_, i) => (
            <div
              key={i}
              className={`w-8 h-8 border-3 border-black rounded-full flex items-center justify-center ${
                i < tracker.failed_elections
                  ? "bg-red-500 shadow-[2px_2px_0px_black]"
                  : "bg-gray-400/50 shadow-[1px_1px_0px_black]"
              }`}
            >
              {i < tracker.failed_elections && (
                <span className="text-white font-bold text-sm">✗</span>
              )}
            </div>
          ))}
        </div>
        <div className="text-center mt-2">
          <span className="text-cream text-xs font-propaganda tracking-wider">
            {tracker.failed_elections}/{tracker.max_failures} FAILED ELECTIONS
          </span>
        </div>
      </div>
    ),
    [tracker]
  );
}

interface TrackProps {
  type: Card;
  slots: number;
  policies: number;
  executiveActions: { [key: number]: Action };
}

const Track = function ({
  type,
  slots,
  policies,
  executiveActions,
}: TrackProps) {
  return useMemo(
    () => (
      <div className="flex justify-center space-x-3 mb-4">
        {Array.from({ length: slots }, (_, i) => (
          <PolicySlot
            key={i}
            isActive={i < policies}
            type={type}
            executiveAction={executiveActions[i]}
          />
        ))}
      </div>
    ),
    [type, slots, policies, executiveActions]
  );
};

export function Board({ board, className = "" }: BoardProps) {
  return (
    <div className={`w-full max-w-4xl mx-auto space-y-6 ${className}`}>
      {/* Fascist Track */}
      <div className="bg-gradient-to-r from-orange-500 to-orange-600 border-4 border-black rounded-xl p-6 shadow-[6px_6px_0px_black]">
        <div className="text-center mb-6">
          <h2 className="text-cream font-propaganda text-2xl md:text-3xl tracking-widest drop-shadow-[2px_2px_0px_black]">
            FASCIST
          </h2>
        </div>

        <Track
          type={CardFascist}
          slots={board.fascist_slots}
          policies={board.fascist_policies}
          executiveActions={board.executive_actions}
        />

        {/* Track decorations */}
        <div className="flex justify-between items-center px-4">
          <div className="w-8 h-8 bg-black rounded-full border-2 border-cream shadow-[2px_2px_0px_cream]"></div>
          <div className="flex-1 h-1 bg-black mx-4"></div>
          <div className="w-8 h-8 bg-black rounded-full border-2 border-cream shadow-[2px_2px_0px_cream]"></div>
        </div>
      </div>

      {/* Liberal Track */}
      <div className="bg-gradient-to-r from-blue-400 to-blue-300 border-4 border-black rounded-xl p-6 shadow-[6px_6px_0px_black]">
        <div className="text-center mb-6">
          <h2 className="text-cream font-propaganda text-2xl md:text-3xl tracking-widest drop-shadow-[2px_2px_0px_black]">
            LIBERAL
          </h2>
        </div>

        <Track
          type={CardLiberal}
          slots={board.liberal_slots}
          policies={board.liberal_policies}
          executiveActions={board.executive_actions}
        />

        {/* Track decorations */}
        <div className="flex justify-between items-center px-4">
          <div className="w-8 h-8 bg-black rounded-full border-2 border-cream shadow-[2px_2px_0px_cream]"></div>
          <div className="flex-1 h-1 bg-black mx-4"></div>
          <div className="w-8 h-8 bg-black rounded-full border-2 border-cream shadow-[2px_2px_0px_cream]"></div>
        </div>

        {/* Election Tracker */}
        <div className="flex justify-center">
          <ElectionTracker tracker={board.election_tracker} />
        </div>
      </div>
    </div>
  );
}
