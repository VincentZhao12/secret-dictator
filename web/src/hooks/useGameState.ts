import {
  MessageTypeActionError,
  MessageTypeGameState,
  type ActionErrorMessage,
  type ActionMessage,
  type GameState,
  type GameStateMessage,
  type Message,
  type MessageType,
} from "@types";
import { useWebSocket } from "./useWebSocket";
import { useState } from "react";

export function useGameState(
  gameId: string,
  playerId: string,
  onError: (error: Error) => void
) {
  const url = `${
    import.meta.env.VITE_BASE_URL_WS
  }/api/v1/play?game=${gameId}&player=${playerId}`;
  const [gameState, setGameState] = useState<GameState | null>(null);

  const { sendMessage: sendMessageRaw, isConnected } = useWebSocket(url, {
    onMessage: (messageEvent) => {
      console.log(messageEvent.data);
      const data: Message = JSON.parse(messageEvent.data);

      console.log(data);

      if (!data || !data.base_message.type) {
        // TODO: log malformed messages
        onError(new Error("Malformed message"));
        return;
      }

      const messageType: MessageType = data.base_message.type;
      switch (messageType) {
        case MessageTypeActionError:
          const errMessage: ActionErrorMessage = data;
          onError(new Error(errMessage.reason));
          break;
        case MessageTypeGameState:
          const gameStateMessage: GameStateMessage = data;
          setGameState(gameStateMessage.game_state);
          break;
        default:
          onError(new Error("Unexpected message type"));
          break;
      }
    },
    // onError: onError,
    deps: [gameId, playerId],
  });

  function sendMessage(message: ActionMessage) {
    sendMessageRaw(JSON.stringify(message));
  }

  return { gameState, sendMessage, isConnected };
}
