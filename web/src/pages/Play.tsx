import { useGameState } from "@hooks";
import Game from "./Game";

export default function Play() {
  // const { gameId, playerId } = useGameParams();

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

  const { gameState, isConnected, sendMessage } = useGameState(
    gameId ?? "",
    getPlayerId(),
    (error) => console.log(error)
  );

  if (!gameState || !isConnected) {
    return <div>connecting</div>;
  }

  return (
    <>
      <Game
        state={gameState}
        currentPlayerId={playerId ?? ""}
        onAction={sendMessage}
      />
    </>
  );
}
