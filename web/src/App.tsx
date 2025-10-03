import "./App.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import JoinGame from "./pages/JoinGame";
import Game from "./pages/Game";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import Play from "./pages/Play";
import { ToastProvider } from "./contexts/ToastContext";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ToastProvider>
        <Router>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/join" element={<JoinGame />} />
            <Route path="/play" element={<Play />} />
            <Route
              path="/game"
              element={
                <Game
                  state={{
                    players: [
                      {
                        id: "1",
                        username: "Alice",
                        role: "liberal",
                        is_executed: false,
                        is_connected: true,
                      },
                      {
                        id: "2",
                        username: "Bob",
                        role: "fascist",
                        is_executed: false,
                        is_connected: true,
                      },
                      {
                        id: "3",
                        username: "Charlie",
                        role: "hitler",
                        is_executed: false,
                        is_connected: true,
                      },
                      {
                        id: "4",
                        username: "Diana",
                        role: "liberal",
                        is_executed: false,
                        is_connected: true,
                      },
                      {
                        id: "5",
                        username: "Eve",
                        role: "hidden",
                        is_executed: true,
                        is_connected: true,
                      },
                      {
                        id: "5",
                        username: "Eve",
                        role: "hidden",
                        is_executed: false,
                        is_connected: true,
                      },
                      {
                        id: "5",
                        username: "Eve",
                        role: "hidden",
                        is_executed: false,
                        is_connected: true,
                      },
                      {
                        id: "5",
                        username: "Eve",
                        role: "hidden",
                        is_executed: false,
                        is_connected: true,
                      },
                      {
                        id: "5",
                        username: "Eve",
                        role: "hidden",
                        is_executed: false,
                        is_connected: true,
                      },
                      {
                        id: "5",
                        username: "Eve",
                        role: "hidden",
                        is_executed: false,
                        is_connected: true,
                      },
                    ],
                    deck: [],
                    discard: [],
                    board: {
                      liberal_policies: 2,
                      fascist_policies: 1,
                      liberal_slots: 5,
                      fascist_slots: 6,
                      executive_actions: {},
                      election_tracker: {
                        failed_elections: 1,
                        max_failures: 3,
                      },
                    },
                    president_index: 0,
                    chancellor_index: 1,
                    prev_president_index: -1,
                    prev_chancellor_index: -1,
                    nominee_index: -1,
                    phase: "election",
                    pending_action: "",
                    peeked_cards: ["liberal", "fascist", "liberal"],
                    peeker_index: 0,
                    resume_order_index: -1,
                    winner: "unassigned",
                    host_id: "1",
                  }}
                  onAction={(actionMessage) => console.log(actionMessage)}
                  currentPlayerId="1"
                  gameId="test-game-id"
                />
              }
            />
          </Routes>
        </Router>
      </ToastProvider>
    </QueryClientProvider>
  );
}

export default App;
