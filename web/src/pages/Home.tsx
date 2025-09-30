import { useNavigate } from "react-router-dom";
import {
  Button,
  Heading,
  Text,
  Card,
  Container,
  Divider,
  Alert,
  Board,
} from "../components";
import { useEffect, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { makePostRequest } from "@util/MakeRequest";
import type { CreateGameResponse } from "../types";

function Home() {
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const createGame = useMutation({
    mutationFn: async () => {
      return await makePostRequest<CreateGameResponse>("games/create", {});
    },
    onSuccess: (data) => {
      setErrorMessage(null); // Clear any previous errors
      navigate(`/join?game=${data.game_id}`);
    },
    onError: (error) => {
      console.error("Failed to create game:", error);
      setErrorMessage("Failed to create game. Please try again.");
    },
  });

  useEffect(() => {
    document.title = "Home - Secret Hitler";
  }, []);

  const handleJoinGame = () => {
    navigate("/join");
  };

  return (
    <Container>
      <Card className="max-w-lg w-full">
        {/* Title */}
        <div className="text-center mb-10 relative z-10">
          <Heading level={1} variant="title">
            SECRET
          </Heading>
          <Heading level={2} variant="subtitle" className="mb-4">
            HITLER
          </Heading>
          <Divider className="w-32 mx-auto" />
        </div>

        {/* Description */}
        <Text variant="description" className="text-center mb-10 relative z-10">
          A game of political intrigue, deception, and betrayal.
        </Text>

        {/* Error Message */}
        {errorMessage && (
          <div className="mb-6 relative z-10">
            <Alert variant="error" onClose={() => setErrorMessage(null)}>
              {errorMessage}
            </Alert>
          </div>
        )}

        {/* Action Buttons */}
        <div className="space-y-5 relative z-10">
          <Button
            onClick={() => createGame.mutate()}
            variant="primary"
            className="w-full"
            loading={createGame.isPending}
          >
            CREATE GAME
          </Button>

          <Button
            onClick={handleJoinGame}
            variant="secondary"
            className="w-full"
          >
            JOIN GAME
          </Button>
        </div>

        {/* Footer */}
        <div className="mt-10 pt-6 border-t-4 border-black relative z-10">
          <Text variant="footer" className="text-center">
            For 5–10 players • Ages 13+
          </Text>
        </div>
      </Card>
    </Container>
  );
}

export default Home;
