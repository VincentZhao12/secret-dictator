import { useGameState } from "@hooks";
import Game from "./Game";
import { useToast } from "../contexts/ToastContext";
import { Heading } from "../components/Heading";
import { Button } from "../components/Button";
import { FaWifi, FaExclamationTriangle, FaRedo } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { ConnectionErrorTypePlayerInvalid } from "@types";

export default function Play() {
  // const { gameId, playerId } = useGameParams();
  const { showError } = useToast();
  const navigate = useNavigate();

  const queryString = window.location.search;
  const urlParams = new URLSearchParams(queryString);
  const gameId = urlParams.get("game");
  const playerId = urlParams.get("player");

  function getPlayerId(): string {
    if (playerId) {
      return playerId;
    }
    return localStorage.getItem(`player_id_for_${gameId}`) ?? "";
  }
  const resolvedPlayerId = getPlayerId();

  const {
    gameState,
    isConnected,
    isConnecting,
    connectionError,
    connectionErrorType,
    lastError,
    sendMessage,
  } = useGameState(gameId ?? "", resolvedPlayerId, (error) => {
    console.log(error);
    showError(error.message);
  });

  // Connecting screen
  if (isConnecting || (!gameState && !connectionError)) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-orange-700 via-orange-600 to-orange-500 font-propaganda">
        <div className="text-center">
          <div className="mb-6 inline-flex items-center justify-center w-20 h-20 rounded-full border-4 border-black shadow-[6px_6px_0px_black] bg-orange-200">
            <FaWifi className="text-black text-3xl" />
          </div>
          <Heading level={1} variant="title" className="mb-2">
            CONNECTING
          </Heading>
          <p className="font-propaganda tracking-wider text-black/80">
            Joining game...
          </p>
        </div>
      </div>
    );
  }

  // Redirect to join screen if player ID is invalid
  if (connectionErrorType === ConnectionErrorTypePlayerInvalid) {
    navigate(`/join?game=${gameId}`);
    return null;
  }

  // Connection error screen
  if (connectionError || (!isConnected && lastError)) {
    const message = connectionError || "Connection lost";

    const handleRetry = () => window.location.reload();
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-orange-700 via-orange-600 to-orange-500 font-propaganda">
        <div className="max-w-md w-full mx-4 text-center">
          <div className="mb-6 inline-flex items-center justify-center w-20 h-20 rounded-full border-4 border-black shadow-[6px_6px_0px_black] bg-red-300">
            <FaExclamationTriangle className="text-black text-3xl" />
          </div>
          <Heading level={1} variant="title" className="mb-2">
            CONNECTION ERROR
          </Heading>
          <p className="font-propaganda tracking-wider text-black/80 mb-6">
            {message}
          </p>
          <div className="flex items-center justify-center gap-3">
            <Button size="md" onClick={handleRetry}>
              <span className="inline-flex items-center gap-2">
                <FaRedo className="text-sm" />
                Retry
              </span>
            </Button>
            <a href="/" className="inline-block">
              <Button size="md" variant="secondary">
                Back Home
              </Button>
            </a>
          </div>
        </div>
      </div>
    );
  }

  return (
    <>
      <Game
        state={gameState!}
        currentPlayerId={resolvedPlayerId}
        onAction={sendMessage}
        gameId={gameId ?? ""}
      />
    </>
  );
}
