import "./App.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import JoinGame from "./pages/JoinGame";
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
          </Routes>
        </Router>
      </ToastProvider>
    </QueryClientProvider>
  );
}

export default App;
