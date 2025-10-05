import { Board, Container, ActionPanel } from "@components";
import {
  type GameState,
  type Player,
  type PlayerRole,
  type ActionMessage,
  type VoteResult,
  Nomination,
  Executive,
  Setup,
  GameOver,
  VotePending,
  VoteJa,
  VoteNein,
  TeamLiberal,
  TeamFascist,
  ActionNominate,
} from "@types";
import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  FaCrown,
  FaGavel,
  FaSkull,
  FaWifi,
  FaLink,
  FaCheck,
  FaTrophy,
  FaHome,
} from "react-icons/fa";

interface GameProps {
  state: GameState;
  currentPlayerId: string;
  onAction: (msg: ActionMessage) => void;
  gameId: string;
}

interface PlayerCardProps {
  player: Player;
  index: number;
  isPresident: boolean;
  isChancellor: boolean;
  isNominee: boolean;
  isCurrentPlayer: boolean;
  vote?: VoteResult;
  onClick: (playerId: string, playerIndex: number) => void;
}

function PlayerCard({
  player,
  index,
  isPresident,
  isChancellor,
  isNominee,
  isCurrentPlayer,
  vote,
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

  return (
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
      {isNominee && !player.is_executed && !isChancellor && (
        <div className="absolute -bottom-2 -left-2 bg-purple-300 text-black px-2 py-0.5 rounded-full border-2 border-black shadow-[2px_2px_0px_black]">
          <span className="text-[10px] font-propaganda font-bold tracking-wider">
            NOMINEE
          </span>
        </div>
      )}

      {/* Vote indicator */}
      {vote !== undefined && vote !== VotePending && (
        <div
          className={`absolute top-1/2 -translate-y-1/2 -right-3 w-10 h-10 border-3 border-black rounded-full flex items-center justify-center shadow-[3px_3px_0px_black] font-bold text-lg ${
            vote === VoteJa
              ? "bg-green-500 text-white"
              : vote === VoteNein
              ? "bg-red-500 text-white"
              : "bg-gray-500 text-white"
          }`}
        >
          {vote === VoteJa ? "✓" : vote === VoteNein ? "✗" : "?"}
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
  );
}

interface PlayerRowProps {
  players: Player[];
  presidentIndex: number;
  chancellorIndex: number;
  nomineeIndex: number;
  currentPlayerId: string;
  votes?: VoteResult[];
  onPlayerClick: (playerId: string, playerIndex: number) => void;
}

function PlayerRow({
  players,
  presidentIndex,
  chancellorIndex,
  nomineeIndex,
  currentPlayerId,
  votes,
  onPlayerClick,
}: PlayerRowProps) {
  return (
    <div className="w-full overflow-x-auto mb-8 pt-4 pb-4">
      <div className="flex justify-start space-x-4 min-w-max px-4">
        {players.map((player, index) => (
          <PlayerCard
            key={player.id}
            player={player}
            index={index}
            isPresident={index === presidentIndex}
            isChancellor={index === chancellorIndex}
            isNominee={index === nomineeIndex}
            isCurrentPlayer={player.id === currentPlayerId}
            vote={votes && votes[index]}
            onClick={onPlayerClick}
          />
        ))}
      </div>
    </div>
  );
}

export default function Game({
  state,
  currentPlayerId,
  onAction,
  gameId,
}: GameProps) {
  console.log(state);
  const navigate = useNavigate();
  const [copied, setCopied] = useState(false);

  // Find current player index from ID
  const currentPlayerIndex = state.players.findIndex(
    (player) => player.id === currentPlayerId
  );

  const currentPlayer =
    currentPlayerIndex >= 0 ? state.players[currentPlayerIndex] : null;
  const isCurrentPlayerDead = currentPlayer?.is_executed || false;

  const handleBackToHome = () => {
    navigate("/");
  };

  const handleCopyInviteLink = async () => {
    const inviteUrl = `${window.location.origin}/join?game=${gameId}`;
    try {
      await navigator.clipboard.writeText(inviteUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy invite link:", err);
    }
  };

  const handlePlayerClick = (playerId: string, playerIndex: number) => {
    switch (state.phase) {
      case Nomination:
        onAction({
          action: ActionNominate,
          target_index: playerIndex,
          base_message: {
            type: "action",
            sender_id: currentPlayerId,
          },
        });
        break;
      case Executive:
        onAction({
          action: state.pending_action ?? "",
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
    onAction(actionMessage);
  };

  return (
    <>
      {/* Game Over Overlay */}
      {state.phase === GameOver && state.winner && (
        <div className="fixed inset-0 bg-black/80 z-50 flex items-center justify-center backdrop-blur-sm">
          <div className="max-w-2xl w-full mx-4">
            <div
              className={`border-8 border-black rounded-2xl p-12 shadow-[12px_12px_0px_black] ${
                state.winner === TeamLiberal
                  ? "bg-gradient-to-br from-blue-400 to-blue-600"
                  : state.winner === TeamFascist
                  ? "bg-gradient-to-br from-red-500 to-orange-600"
                  : "bg-gray-500"
              }`}
            >
              {/* Trophy Icon */}
              <div className="flex justify-center mb-6">
                <div className="bg-yellow-400 border-6 border-black rounded-full p-8 shadow-[8px_8px_0px_black]">
                  <FaTrophy className="text-6xl text-black" />
                </div>
              </div>

              {/* Winner Text */}
              <div className="text-center space-y-4">
                <h1 className="font-propaganda text-6xl tracking-wider text-white drop-shadow-[4px_4px_0px_black]">
                  GAME OVER
                </h1>
                <div
                  className={`inline-block px-8 py-4 border-4 border-black rounded-xl shadow-[6px_6px_0px_black] ${
                    state.winner === TeamLiberal
                      ? "bg-blue-200"
                      : state.winner === TeamFascist
                      ? "bg-orange-200"
                      : "bg-gray-200"
                  }`}
                >
                  <h2 className="font-propaganda text-4xl tracking-widest text-black">
                    {state.winner === TeamLiberal
                      ? "LIBERALS WIN!"
                      : state.winner === TeamFascist
                      ? "FASCISTS WIN!"
                      : "DRAW"}
                  </h2>
                </div>

                {/* Team Roster */}
                <div className="mt-8 bg-white/90 border-4 border-black rounded-xl p-6 shadow-[4px_4px_0px_black]">
                  <h3 className="font-propaganda text-xl tracking-wider mb-4 text-black">
                    TEAM REVEAL
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {state.players.map((player, index) => (
                      <div
                        key={index}
                        className={`flex items-center justify-between p-3 border-3 border-black rounded-lg ${
                          player.role === "liberal"
                            ? "bg-blue-100"
                            : player.role === "fascist"
                            ? "bg-red-100"
                            : player.role === "hitler"
                            ? "bg-black text-white"
                            : "bg-gray-100"
                        }`}
                      >
                        <span className="font-propaganda text-sm">
                          {player.username}
                        </span>
                        <span className="font-propaganda text-xs font-bold">
                          {player.role === "liberal"
                            ? "LIBERAL"
                            : player.role === "fascist"
                            ? "FASCIST"
                            : player.role === "hitler"
                            ? "HITLER"
                            : "?"}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Back to Home Button */}
                <div className="mt-6">
                  <button
                    onClick={handleBackToHome}
                    className="w-full bg-gray-800 hover:bg-gray-900 text-white px-6 py-4 rounded-lg border-4 border-black shadow-[4px_4px_0px_black] font-propaganda font-bold tracking-wider uppercase transition-all duration-200 hover:scale-105 flex items-center justify-center space-x-2"
                  >
                    <FaHome className="text-xl" />
                    <span>BACK TO HOME</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Invite Link Button - Top Right */}
      {state.phase === Setup && (
        <div className="fixed top-6 right-6 z-40">
          <button
            onClick={handleCopyInviteLink}
            className="bg-orange-600 hover:bg-orange-700 text-black px-4 py-3 rounded-lg border-4 border-black shadow-[4px_4px_0px_black] font-propaganda font-bold tracking-wider uppercase transition-all duration-200 hover:scale-105 flex items-center space-x-2"
          >
            {copied ? (
              <>
                <FaCheck className="text-sm" />
                <span>COPIED!</span>
              </>
            ) : (
              <>
                <FaLink className="text-sm" />
                <span>INVITE LINK</span>
              </>
            )}
          </button>
        </div>
      )}

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
            nomineeIndex={state.nominee_index}
            currentPlayerId={currentPlayerId}
            votes={state.votes}
            onPlayerClick={handlePlayerClick}
          />

          {/* Main Game Area - Board and Action Panel side by side */}
          <div className="flex flex-col lg:flex-row gap-6 items-start">
            {/* Game Board */}
            <div className="flex-1">
              <Board
                board={state.board}
                deckCount={state.deck?.length || 0}
                discardCount={state.discard?.length || 0}
              />
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
