import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Button,
  Heading,
  Text,
  Card,
  Container,
  Input,
  Divider,
  Alert,
} from "../components";
import "../App.css";
import { useMutation } from "@tanstack/react-query";
import { makePostRequest } from "@util/MakeRequest";
import { type JoinGameResponse, type JoinGameRequest } from "@types";

function JoinGame() {
  const navigate = useNavigate();
  const [gameId, setGameId] = useState("");
  const [username, setUsername] = useState("");
  const [errors, setErrors] = useState<{
    join?: string;
    gameId?: string;
    username?: string;
  }>({});

  const joinGame = useMutation({
    mutationFn: async (request: JoinGameRequest) => {
      return makePostRequest<JoinGameResponse>("games/join", request);
    },
    onError: (error) => {
      const newErrors: { join?: string; gameId?: string; username?: string } = {
        join: error.message,
      };

      setErrors(newErrors);
    },
    onSuccess: (resp) => {
      // TODO: Implement transition to play after play exists
      localStorage.setItem(`player_id_for_${resp.game_id}`, resp.player_id);
      localStorage.setItem("last_username", username);
      navigate(`/play?game=${resp.game_id}&player=${resp.player_id}`);
      console.log(`joined game ${resp.player_id}`);
    },
  });

  useEffect(() => {
    document.title = "Join Game - Secret Hitler";
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    const gameParam = urlParams.get("game");
    if (gameParam) {
      setGameId(gameParam);
    }

    // Load last username from localStorage
    const lastUsername = localStorage.getItem("last_username");
    if (lastUsername) {
      setUsername(lastUsername);
    }
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Basic validation
    const newErrors: { join?: string; gameId?: string; username?: string } = {};

    if (!gameId.trim()) {
      newErrors.gameId = "Game ID is required";
    }

    if (!username.trim()) {
      newErrors.username = "Username is required";
    } else if (username.trim().length < 2) {
      newErrors.username = "Username must be at least 2 characters";
    }

    setErrors(newErrors);

    // If no errors, proceed with joining
    if (Object.keys(newErrors).length === 0) {
      joinGame.mutate({
        username,
        game_id: gameId,
      });
    }
  };

  const handleBackToHome = () => {
    navigate("/");
  };

  return (
    <Container>
      <Card className="max-w-md w-full">
        {/* Title */}
        <div className="text-center mb-8">
          <Heading level={1} variant="section" className="mb-2">
            JOIN GAME
          </Heading>
          <Divider className="w-24 mx-auto" />
        </div>

        {/* Description */}
        <Text variant="description" className="text-center mb-8">
          Enter the game ID and your username to join an existing game.
        </Text>

        {errors.join && (
          <div className="mb-6 relative z-10">
            <Alert variant="error" onClose={() => setErrors({})}>
              {errors.join}
            </Alert>
          </div>
        )}

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-6">
          <Input
            label="Game ID"
            type="text"
            placeholder="Enter game ID"
            value={gameId}
            onChange={(e) => setGameId(e.target.value)}
            error={errors.gameId}
            required
          />

          <Input
            label="Username"
            type="text"
            placeholder="Enter your username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            error={errors.username}
            required
            maxLength={20}
          />

          <div className="space-y-3">
            <Button
              type="submit"
              variant="primary"
              className="w-full"
              loading={joinGame.isPending}
            >
              JOIN GAME
            </Button>

            <Button
              type="button"
              variant="secondary"
              className="w-full"
              onClick={handleBackToHome}
            >
              BACK TO HOME
            </Button>
          </div>
        </form>

        {/* Footer */}
        <div className="mt-8 pt-6 border-t-4 border-black">
          <Text variant="footer" className="text-center">
            Secret Hitler • Online Multiplayer
          </Text>
        </div>
      </Card>
    </Container>
  );
}

export default JoinGame;
