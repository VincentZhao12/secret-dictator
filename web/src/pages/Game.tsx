import { Board, Container, ActionPanel } from "@components";
import {
  type GameState,
  type Player,
  type PlayerRole,
  type ActionMessage,
  Election,
  Nomination,
  Executive,
} from "@types";
import { useMemo } from "react";
import { FaCrown, FaGavel, FaSkull, FaWifi } from "react-icons/fa";

interface GameProps {
  state: GameState;
  currentPlayerId: string;
  onAction: (msg: ActionMessage) => void;
}

interface PlayerCardProps {
  player: Player;
  index: number;
  isPresident: boolean;
  isChancellor: boolean;
  isCurrentPlayer: boolean;
  onClick: (playerId: string, playerIndex: number) => void;
}

function PlayerCard({
  player,
  index,
  isPresident,
  isChancellor,
  isCurrentPlayer,
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
        onClick={() => !player.is_executed && onClick(player.id, index)}
        className={`relative border-4 rounded-xl p-3 shadow-[4px_4px_0px_black] transition-transform duration-200 min-w-[120px] ${
          player.is_executed
            ? "bg-gray-400/60 cursor-default opacity-75 border-black"
            : player.is_connected
            ? "bg-orange-200/90 cursor-pointer hover:scale-105 border-black"
            : "bg-orange-200/60 cursor-pointer hover:scale-105 border-red-500"
        }`}
      >
        {/* Disconnected indicator */}
        {!player.is_connected && !player.is_executed && (
          <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-red-500 text-white px-2 py-0.5 rounded-full border-2 border-black shadow-[2px_2px_0px_black] flex items-center space-x-1">
            <FaWifi className="text-xs rotate-45" />
            <span className="text-xs font-propaganda font-bold">OFFLINE</span>
          </div>
        )}

        {/* Dead player overlay */}
        {player.is_executed && (
          <div className="absolute inset-0 bg-black/20 rounded-lg flex items-center justify-center">
            <FaSkull className="text-red-600 text-2xl" />
          </div>
        )}

        {/* Role indicator */}
        <div
          className={`absolute -top-2 -left-2 w-8 h-8 rounded-full border-2 border-black flex items-center justify-center font-propaganda text-sm font-bold ${getRoleColor(
            player.role
          )} shadow-[2px_2px_0px_black] ${
            player.is_executed ? "opacity-50" : ""
          }`}
        >
          {getRoleDisplay(player.role)}
        </div>

        {/* President/Chancellor indicators */}
        {isPresident && !player.is_executed && (
          <div className="absolute -top-2 -right-2 w-8 h-8 bg-yellow-500 border-2 border-black rounded-full flex items-center justify-center shadow-[2px_2px_0px_black]">
            <FaCrown className="text-black text-sm" />
          </div>
        )}
        {isChancellor && !player.is_executed && (
          <div className="absolute -bottom-2 -right-2 w-8 h-8 bg-purple-600 border-2 border-black rounded-full flex items-center justify-center shadow-[2px_2px_0px_black]">
            <FaGavel className="text-white text-sm" />
          </div>
        )}

        {/* Username */}
        <div className="pt-2 pb-1 text-center">
          <span
            className={`font-propaganda text-sm tracking-wider ${
              player.is_executed
                ? "text-gray-600 line-through"
                : !player.is_connected
                ? "text-gray-500"
                : "text-black"
            }`}
          >
            {player.username.toUpperCase()}
          </span>
          {isCurrentPlayer && (
            <div className="mt-1">
              <span className="inline-block bg-green-500 text-white px-2 py-0.5 rounded-full text-xs font-propaganda font-bold border-2 border-black shadow-[2px_2px_0px_black]">
                YOU
              </span>
            </div>
          )}
        </div>
      </div>
    ),
    [player, index, isPresident, isChancellor, isCurrentPlayer, onClick]
  );
}

interface PlayerRowProps {
  players: Player[];
  presidentIndex: number;
  chancellorIndex: number;
  currentPlayerId: string;
  onPlayerClick: (playerId: string, playerIndex: number) => void;
}

function PlayerRow({
  players,
  presidentIndex,
  chancellorIndex,
  currentPlayerId,
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
              isCurrentPlayer={player.id === currentPlayerId}
              onClick={onPlayerClick}
            />
          ))}
        </div>
      </div>
    ),
    [players, presidentIndex, chancellorIndex, currentPlayerId, onPlayerClick]
  );
}

export default function Game({ state, currentPlayerId, onAction }: GameProps) {
  // Find current player index from ID
  const currentPlayerIndex = state.players.findIndex(
    (player) => player.id === currentPlayerId
  );

  const currentPlayer =
    currentPlayerIndex >= 0 ? state.players[currentPlayerIndex] : null;
  const isCurrentPlayerDead = currentPlayer?.is_executed || false;

  const handlePlayerClick = (playerId: string, playerIndex: number) => {
    console.log(`Clicked player: ${playerId} at index ${playerIndex}`);

    switch (state.phase) {
      case Nomination:
        onAction({
          action: "nominate",
          target_index: playerIndex,
          base_message: {
            type: "action",
            sender_id: currentPlayerId,
          },
        });
        break;
      case Executive:
        onAction({
          action: "execute",
          target_index: playerIndex,
          base_message: {
            type: "action",
            sender_id: currentPlayerId,
          },
        });
        break;
    }
  };

  const handleAction = (actionMessage: ActionMessage) => {
    console.log(`Action: ${actionMessage}`);
    onAction(actionMessage);
  };

  return (
    <>
      {/* Full screen overlay for dead players */}
      {isCurrentPlayerDead && (
        <div className="fixed inset-0 bg-black/30 z-50 pointer-events-none">
          {/* Execution alert in bottom right */}
          <div className="absolute bottom-6 right-6">
            <div className="bg-black/90 text-white px-6 py-4 rounded-xl border-4 border-red-600 shadow-[6px_6px_0px_red-800]">
              <div className="flex items-center space-x-3">
                <FaSkull className="text-red-500 text-2xl" />
                <div>
                  <h2 className="font-propaganda text-xl tracking-wider text-red-400">
                    YOU HAVE BEEN EXECUTED
                  </h2>
                  <p className="font-propaganda text-sm text-gray-300 mt-1">
                    You can observe but cannot participate
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      <Container>
        <div className="w-full space-y-6">
          {/* Players Row */}
          <PlayerRow
            players={state.players}
            presidentIndex={state.president_index}
            chancellorIndex={state.chancellor_index}
            currentPlayerId={currentPlayerId}
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
                currentPlayerId={currentPlayerId}
                onAction={handleAction}
              />
            </div>
          </div>
        </div>
      </Container>
    </>
  );
}
