import { useMemo, useState } from "react";
import { Button } from "./Button";
import { PolicyCard } from "./PolicyCard";
import { FaThumbsUp, FaThumbsDown, FaPlay } from "react-icons/fa";
import type { GameState, ActionMessage } from "../types";
import {
  ActionVote,
  ActionLegislate,
  ActionEndTurn,
  ActionStartGame,
  Election,
  Legislation1,
  Legislation2,
  Executive,
  Nomination,
  ActionInvestigate,
  ActionExecution,
  ActionSpecialElection,
  ActionPolicyPeek,
  Setup,
  GameOver,
  Paused,
  MessageTypeAction,
  VotePending,
  VoteJa,
  VoteNein,
} from "../types";

interface ActionPanelProps {
  gameState: GameState;
  currentPlayerId: string;
  onAction: (action: ActionMessage) => void;
}

export function ActionPanel({
  gameState,
  currentPlayerId,
  onAction,
}: ActionPanelProps) {
  const [selectedCard, setSelectedCard] = useState<number>(-1);

  // Find current player index from ID
  const currentPlayerIndex = gameState.players.findIndex(
    (player) => player.id === currentPlayerId
  );

  const currentPlayer =
    currentPlayerIndex >= 0 ? gameState.players[currentPlayerIndex] : null;
  const isCurrentPlayerDead = currentPlayer?.is_executed || false;

  const isCurrentPlayerTurn = (requiredIndex: number) => {
    return currentPlayerIndex === requiredIndex;
  };

  const getPhaseTitle = () => {
    switch (gameState.phase) {
      case Election:
        return "ELECTION - Vote on the Government";
      case Legislation1:
        return "LEGISLATION - President selects 2 policies";
      case Legislation2:
        return "LEGISLATION - Chancellor selects 1 policy";
      case Executive:
        return "EXECUTIVE ACTION";
      case Nomination:
        return "NOMINATION - President nominates a chancellor";
      case Setup:
        return "SETUP - Game host sets up the game";
      case GameOver:
        return `GAME OVER - Game over ${gameState.winner}`;
      case Paused:
        return "PAUSED - Game paused";
      default:
        return "Waiting for next phase...";
    }
  };

  const canVote = () => {
    return gameState.phase === Election;
  };

  const canLegislate = () => {
    return (
      (gameState.phase === Legislation1 &&
        isCurrentPlayerTurn(gameState.president_index)) ||
      (gameState.phase === Legislation2 &&
        isCurrentPlayerTurn(gameState.chancellor_index))
    );
  };

  const handleVote = (vote: boolean) => {
    onAction({
      action: ActionVote,
      vote: vote,
      base_message: {
        sender_id: currentPlayerId,
        type: MessageTypeAction,
      },
    });
  };

  const handleCardSelection = (cardIndex: number) => {
    if (selectedCard === cardIndex) {
      setSelectedCard(-1);
    } else {
      setSelectedCard(cardIndex);
    }
  };

  const handleLegislate = () => {
    if (selectedCard !== -1) {
      onAction({
        action: ActionLegislate,
        target_index: selectedCard,
        base_message: {
          sender_id: currentPlayerId,
          type: MessageTypeAction,
        },
      });
      setSelectedCard(-1);
    }
  };

  const getPeekedCards = () => {
    return gameState.peeked_cards || [];
  };

  const getRequiredSelections = () => {
    if (gameState.phase === Legislation1) return 1; // President discards 1, keeps 2
    if (gameState.phase === Legislation2) return 1; // Chancellor discards 1, keeps 1
    return 0;
  };

  const getExecutiveMessage = () => {
    switch (gameState.pending_action) {
      case ActionInvestigate:
        return "Click on a player to investigate";
      case ActionExecution:
        return "Click on a player to execute";
      case ActionSpecialElection:
        return "Click on a player to elect as president";
      case ActionPolicyPeek:
        return "End turn once you are done peeking";
      default:
        return "";
    }
  };

  const handleStartGame = () => {
    onAction({
      action: ActionStartGame,
      base_message: {
        sender_id: currentPlayerId,
        type: MessageTypeAction,
      },
    });
  };

  const renderVotingInterface = () => {
    if (!canVote()) return null;

    // Check if current player has already voted
    const currentVote = gameState.votes && gameState.votes[currentPlayerIndex];
    const hasVoted = currentVote !== undefined && currentVote !== VotePending;

    if (hasVoted) {
      return (
        <div className="text-center mb-4">
          <div className="bg-green-100 border-4 border-black rounded-lg px-6 py-4 shadow-[4px_4px_0px_black] inline-block">
            <div className="flex items-center space-x-3">
              {currentVote === VoteJa ? (
                <FaThumbsUp className="text-green-600 text-xl" />
              ) : currentVote === VoteNein ? (
                <FaThumbsDown className="text-red-600 text-xl" />
              ) : null}
              <div className="text-left">
                <p className="text-black font-propaganda text-sm font-bold">
                  VOTE SUBMITTED
                </p>
                <p className="text-gray-700 font-propaganda text-xs">
                  Waiting for other players...
                </p>
              </div>
            </div>
          </div>
        </div>
      );
    }

    return (
      <div className="flex justify-center space-x-4 mb-4">
        <Button
          onClick={() => handleVote(true)}
          variant="primary"
          className="flex items-center space-x-2"
        >
          <FaThumbsUp />
          <span>JA</span>
        </Button>
        <Button
          onClick={() => handleVote(false)}
          variant="secondary"
          className="flex items-center space-x-2"
        >
          <FaThumbsDown />
          <span>NEIN</span>
        </Button>
      </div>
    );
  };

  const renderPolicySelection = () => {
    if (!canLegislate() || getPeekedCards().length === 0) return null;

    return (
      <div className="space-y-4">
        <div className="text-center">
          <p className="text-black font-propaganda text-sm">
            Select {getRequiredSelections()} card(s) to discard:
          </p>
        </div>
        <div className="flex justify-center space-x-3">
          {getPeekedCards().map((card, index) => {
            const isSelected = selectedCard === index;
            return (
              <div
                key={index}
                className={`cursor-pointer transition-transform duration-150 ${
                  isSelected ? "scale-110 z-10" : "hover:scale-105"
                }`}
              >
                <PolicyCard
                  type={card}
                  onClick={() => handleCardSelection(index)}
                />
                {isSelected && (
                  <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
                    <span className="bg-yellow-400/80 text-black font-bold px-2 py-1 rounded shadow-lg text-xs border-2 border-black">
                      SELECTED
                    </span>
                  </div>
                )}
              </div>
            );
          })}
        </div>
        <div className="flex justify-center">
          <Button
            onClick={handleLegislate}
            variant="primary"
            disabled={selectedCard === -1}
          >
            CONFIRM SELECTION
          </Button>
        </div>
      </div>
    );
  };

  const renderExecutiveAction = () => {
    if (gameState.phase !== Executive) return null;

    const isCurrentPlayerPresident =
      gameState.president_index === currentPlayerIndex;

    if (isCurrentPlayerPresident) {
      return (
        <div className="text-center">
          <div className="text-black font-propaganda text-sm flex flex-col items-center justify-center">
            <span>{getExecutiveMessage()}</span>
            {gameState.pending_action === ActionPolicyPeek && (
              <div className="flex justify-center space-x-3 mt-2 relative min-h-[5.5rem]">
                {getPeekedCards().map((card, index) => (
                  <div
                    key={index}
                    className="transition-transform duration-150 scale-110 z-10 relative static"
                  >
                    <PolicyCard type={card} />
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      );
    }

    return (
      <div className="text-center">
        <p className="text-black font-propaganda text-sm flex items-center justify-center space-x-2">
          <span>Waiting for president to perform executive action</span>
        </p>
      </div>
    );
  };

  const renderEndTurnButton = () => {
    const shouldShowEndTurn =
      gameState.phase === Executive &&
      (gameState.pending_action === ActionInvestigate ||
        gameState.pending_action === ActionPolicyPeek) &&
      gameState.president_index === currentPlayerIndex;

    if (!shouldShowEndTurn) return null;

    return (
      <div className="flex justify-center">
        <Button
          onClick={() =>
            onAction({
              action: ActionEndTurn,
              base_message: {
                sender_id: currentPlayerId,
                type: MessageTypeAction,
              },
            })
          }
          variant="secondary"
        >
          END TURN
        </Button>
      </div>
    );
  };

  const renderStartGameButton = () => {
    const isHost = currentPlayerId === gameState.host_id;
    const isSetupPhase = gameState.phase === Setup;
    const hasEnoughPlayers = gameState.players.length >= 5;

    if (!isHost || !isSetupPhase) return null;

    return (
      <div className="flex flex-col items-center space-y-4">
        <button
          onClick={handleStartGame}
          disabled={!hasEnoughPlayers}
          className={`
            relative overflow-hidden
            font-propaganda text-2xl tracking-wider
            px-12 py-6
            border-4 border-black rounded-xl
            shadow-[6px_6px_0px_black]
            transition-all duration-200
            flex items-center justify-center space-x-4
            ${
              hasEnoughPlayers
                ? "bg-green-500 hover:bg-green-600 hover:scale-105 active:scale-95 active:shadow-[3px_3px_0px_black] cursor-pointer text-white"
                : "bg-gray-400 cursor-not-allowed opacity-60 text-gray-600"
            }
          `}
        >
          <FaPlay className="text-2xl" />
          <span>START GAME</span>
        </button>
        {!hasEnoughPlayers && (
          <div className="bg-red-200/90 border-4 border-black rounded-lg px-6 py-3 shadow-[4px_4px_0px_black]">
            <p className="text-red-800 font-propaganda text-sm tracking-wide">
              ⚠️ Need at least 5 players to start
            </p>
          </div>
        )}
        {hasEnoughPlayers && (
          <p className="text-green-700 font-propaganda text-sm tracking-wide">
            ✓ Ready to begin!
          </p>
        )}
      </div>
    );
  };

  return useMemo(
    () => (
      <div
        className={`border-4 border-black rounded-xl p-6 shadow-[6px_6px_0px_black] mb-6 ${
          isCurrentPlayerDead ? "bg-gray-400/60 opacity-75" : "bg-orange-200/90"
        }`}
      >
        {/* Phase Title */}
        <div className="text-center mb-4">
          <h3
            className={`font-propaganda text-xl tracking-wider ${
              isCurrentPlayerDead ? "text-gray-600" : "text-black"
            }`}
          >
            {isCurrentPlayerDead ? "YOU ARE DEAD" : getPhaseTitle()}
          </h3>
          {isCurrentPlayerDead && (
            <p className="text-gray-600 font-propaganda text-sm mt-2">
              You cannot participate in the game
            </p>
          )}
        </div>

        {/* Only show interactive content if player is alive */}
        {!isCurrentPlayerDead && (
          <>
            {/* Start Game Button (Setup Phase) */}
            {renderStartGameButton()}

            {/* Voting Interface */}
            {renderVotingInterface()}

            {/* Policy Selection Interface */}
            {renderPolicySelection()}

            {/* Executive Action Interface */}
            {renderExecutiveAction()}

            {/* End Turn Button */}
            {renderEndTurnButton()}
          </>
        )}
      </div>
    ),
    [gameState, currentPlayerId, selectedCard, onAction, isCurrentPlayerDead]
  );
}
