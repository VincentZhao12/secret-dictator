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
} from "../components";
import "../App.css";

function JoinGame() {
  const navigate = useNavigate();
  const [gameId, setGameId] = useState("");
  const [username, setUsername] = useState("");
  const [errors, setErrors] = useState<{ gameId?: string; username?: string }>(
    {}
  );

  useEffect(() => {
    document.title = "Join Game - Secret Hitler";
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    const gameParam = urlParams.get("game");
    if (gameParam) {
      setGameId(gameParam);
    }
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // Basic validation
    const newErrors: { gameId?: string; username?: string } = {};

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
      console.log("Joining game:", {
        gameId: gameId.trim(),
        username: username.trim(),
      });
      // TODO: Implement actual game joining logic
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
            <Button type="submit" variant="primary" className="w-full">
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
