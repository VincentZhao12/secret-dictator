import { makePostRequest } from "./MakeRequest";
import type {
  CreateGameResponse,
  JoinGameRequest,
  JoinGameResponse,
  AddBotRequest,
  AddBotResponse,
  RemoveBotRequest,
  RemovePlayerRequest,
} from "../types";

export const createGame = async (): Promise<CreateGameResponse> => {
  return makePostRequest<CreateGameResponse>("games/create", {});
};

export const joinGame = async (
  request: JoinGameRequest
): Promise<JoinGameResponse> => {
  return makePostRequest<JoinGameResponse>("games/join", request);
};

export const addBot = async (
  request: AddBotRequest
): Promise<AddBotResponse> => {
  return makePostRequest<AddBotResponse>("games/bots/add", request);
};

export const removeBot = async (request: RemoveBotRequest): Promise<void> => {
  return makePostRequest<void>("games/bots/remove", request);
};

export const removePlayer = async (
  request: RemovePlayerRequest
): Promise<void> => {
  return makePostRequest<void>("games/players/remove", request);
};
