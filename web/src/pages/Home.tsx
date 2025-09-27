import { Button, Heading, Text, Card, Container, Divider } from "../components";

function Home() {
  const handleCreateGame = () => {
    console.log("Creating new game...");
  };

  const handleJoinGame = () => {
    console.log("Joining game...");
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

        {/* Action Buttons */}
        <div className="space-y-5 relative z-10">
          <Button
            onClick={handleCreateGame}
            variant="primary"
            className="w-full"
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
