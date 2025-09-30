import { Board, Container, ActionPanel } from "../components";
import type { GameState, Player, PlayerRole, Action } from "../types";
import { useMemo } from "react";
import { FaCrown, FaGavel } from "react-icons/fa";

interface GameProps {
  state: GameState;
}

interface PlayerCardProps {
  player: Player;
  index: number;
  isPresident: boolean;
  isChancellor: boolean;
  onClick: (playerId: string, playerIndex: number) => void;
}

function PlayerCard({
  player,
  index,
  isPresident,
  isChancellor,
  onClick,
}: PlayerCardProps) {
  const getRoleColor = (role: PlayerRole) => {
    switch (role) {
      case "liberal":
        return "bg-blue-500 text-white";
      case "fascist":
        return "bg-red-600 text-white";
      case "hitler":
        return "bg-black text-red-500";
      case "hidden":
      default:
        return "bg-gray-600 text-white";
    }
  };

  const getRoleDisplay = (role: PlayerRole) => {
    switch (role) {
      case "liberal":
        return "L";
      case "fascist":
        return "F";
      case "hitler":
        return "H";
      case "hidden":
      default:
        return "?";
    }
  };

  return useMemo(
    () => (
      <div
        onClick={() => onClick(player.id, index)}
        className="relative bg-orange-200/90 border-4 border-black rounded-xl p-3 shadow-[4px_4px_0px_black] cursor-pointer hover:scale-105 transition-transform duration-200 min-w-[120px]"
      >
        {/* Role indicator */}
        <div
          className={`absolute -top-2 -left-2 w-8 h-8 rounded-full border-2 border-black flex items-center justify-center font-propaganda text-sm font-bold ${getRoleColor(
            player.role
          )} shadow-[2px_2px_0px_black]`}
        >
          {getRoleDisplay(player.role)}
        </div>

        {/* President/Chancellor indicators */}
        {isPresident && (
          <div className="absolute -top-2 -right-2 w-8 h-8 bg-yellow-500 border-2 border-black rounded-full flex items-center justify-center shadow-[2px_2px_0px_black]">
            <FaCrown className="text-black text-sm" />
          </div>
        )}
        {isChancellor && (
          <div className="absolute -bottom-2 -right-2 w-8 h-8 bg-purple-600 border-2 border-black rounded-full flex items-center justify-center shadow-[2px_2px_0px_black]">
            <FaGavel className="text-white text-sm" />
          </div>
        )}

        {/* Username */}
        <div className="pt-2 pb-1 text-center">
          <span className="font-propaganda text-black text-sm tracking-wider">
            {player.username.toUpperCase()}
          </span>
        </div>
      </div>
    ),
    [player, index, isPresident, isChancellor, onClick]
  );
}

interface PlayerRowProps {
  players: Player[];
  presidentIndex: number;
  chancellorIndex: number;
  onPlayerClick: (playerId: string, playerIndex: number) => void;
}

function PlayerRow({
  players,
  presidentIndex,
  chancellorIndex,
  onPlayerClick,
}: PlayerRowProps) {
  return useMemo(
    () => (
      <div className="w-full overflow-x-auto mb-8 pt-4 pb-4">
        <div className="flex justify-start space-x-4 min-w-max px-4">
          {players.map((player, index) => (
            <PlayerCard
              key={player.id}
              player={player}
              index={index}
              isPresident={index === presidentIndex}
              isChancellor={index === chancellorIndex}
              onClick={onPlayerClick}
            />
          ))}
        </div>
      </div>
    ),
    [players, presidentIndex, chancellorIndex, onPlayerClick]
  );
}

export default function Game({ state }: GameProps) {
  // In a real app, this would come from authentication/game context
  const currentPlayerIndex = 0; // For demo purposes

  const handlePlayerClick = (playerId: string, playerIndex: number) => {
    console.log(`Clicked player: ${playerId} at index ${playerIndex}`);
    // TODO: Handle player interactions based on game phase
    // This could trigger investigation, execution, special election, etc.
  };

  const handleAction = (action: Action, data?: any) => {
    console.log(`Action: ${action}`, data);
    // TODO: Send action to game server
    // This would handle all game actions like voting, legislation, veto, etc.
  };

  return (
    <Container>
      <div className="w-full space-y-6">
        {/* Players Row */}
        <PlayerRow
          players={state.players}
          presidentIndex={state.president_index}
          chancellorIndex={state.chancellor_index}
          onPlayerClick={handlePlayerClick}
        />

        {/* Main Game Area - Board and Action Panel side by side */}
        <div className="flex flex-col lg:flex-row gap-6 items-start">
          {/* Game Board */}
          <div className="flex-1">
            <Board board={state.board} />
          </div>

          {/* Action Panel */}
          <div className="w-full lg:w-96 flex-shrink-0">
            <ActionPanel
              gameState={state}
              currentPlayerIndex={currentPlayerIndex}
              onAction={handleAction}
            />
          </div>
        </div>
      </div>
    </Container>
  );
}
