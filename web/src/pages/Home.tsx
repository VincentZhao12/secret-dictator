import { useNavigate } from "react-router-dom";
import {
  Button,
  Heading,
  Text,
  Card,
  Container,
  Divider,
  Alert,
} from "../components";
import { useEffect, useState } from "react";
import { FaChevronDown, FaChevronUp, FaInfoCircle } from "react-icons/fa";
import { useMutation } from "@tanstack/react-query";
import { createGame as createGameApi } from "@util/GameApi";

function Home() {
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [showDisclosures, setShowDisclosures] = useState(false);

  const createGame = useMutation({
    mutationFn: createGameApi,
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
            DICTATOR
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
          <Text variant="footer" className="text-center mb-4">
            For 5–10 players • Ages 13+
          </Text>

          {/* Disclosures */}
          <div className="text-center">
            <button
              onClick={() => setShowDisclosures(!showDisclosures)}
              className="inline-flex items-center gap-2 text-xs text-black/60 hover:text-black/80 transition-colors duration-200 font-propaganda"
            >
              <FaInfoCircle className="text-sm" />
              <span>Legal & Attribution</span>
              {showDisclosures ? (
                <FaChevronUp className="text-xs" />
              ) : (
                <FaChevronDown className="text-xs" />
              )}
            </button>

            {showDisclosures && (
              <div className="mt-4 p-4 bg-black/5 border-2 border-black/10 rounded-lg text-left">
                <div className="text-xs text-black/70 space-y-2 font-propaganda leading-relaxed">
                  <p>
                    Based on{" "}
                    <a
                      href="https://www.secrethitler.com/"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="underline hover:text-black/90"
                    >
                      Secret Hitler
                    </a>
                    .
                  </p>
                  <p>
                    This project is a non-commercial fan implementation and is
                    not affiliated with the original creators.
                  </p>
                  <p>
                    Secret Hitler is licensed under{" "}
                    <a
                      href="https://creativecommons.org/licenses/by-nc-sa/4.0/"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="underline hover:text-black/90"
                    >
                      Creative Commons BY–NC–SA 4.0
                    </a>
                    .
                  </p>
                  <p>
                    You must give credit to the original creators, may not use
                    this for commercial purposes, and any adaptations must be
                    licensed under the same terms.
                  </p>
                  <p>
                    This is not an exhaustive list of requirements; please
                    ensure you comply with the original license.
                  </p>
                </div>
              </div>
            )}
          </div>
        </div>
      </Card>
    </Container>
  );
}

export default Home;
