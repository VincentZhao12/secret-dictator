import "./App.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import JoinGame from "./pages/JoinGame";
import Game from "./pages/Game";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/join" element={<JoinGame />} />
          <Route
            path="/game"
            element={
              <Game
                state={{
                  players: [
                    { id: "1", username: "Alice", role: "liberal" },
                    { id: "2", username: "Bob", role: "fascist" },
                    { id: "3", username: "Charlie", role: "hitler" },
                    { id: "4", username: "Diana", role: "liberal" },
                    { id: "5", username: "Eve", role: "hidden" },
                    { id: "5", username: "Eve", role: "hidden" },
                    { id: "5", username: "Eve", role: "hidden" },
                    { id: "5", username: "Eve", role: "hidden" },
                    { id: "5", username: "Eve", role: "hidden" },
                    { id: "5", username: "Eve", role: "hidden" },
                  ],
                  deck: [],
                  discard: [],
                  board: {
                    liberal_policies: 2,
                    fascist_policies: 1,
                    liberal_slots: 5,
                    fascist_slots: 6,
                    executive_actions: {},
                    election_tracker: { failed_elections: 1, max_failures: 3 },
                  },
                  president_index: 0,
                  chancellor_index: 1,
                  prev_president_index: -1,
                  prev_chancellor_index: -1,
                  nominee_index: -1,
                  phase: "setup",
                  pending_action: "",
                  peeked_cards: ["liberal", "fascist", "liberal"],
                  peeker_index: 0,
                  executed_players: [],
                }}
              />
            }
          />
        </Routes>
      </Router>
    </QueryClientProvider>
  );
}

export default App;
