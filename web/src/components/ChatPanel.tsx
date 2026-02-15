import { useEffect, useRef, useState } from "react";
import type { ChatEntry } from "@types";

interface ChatPanelProps {
  chatHistory: ChatEntry[];
  currentPlayerId: string;
  onSend: (text: string) => void;
}

const senderTextColorClasses = [
  "text-blue-700",
  "text-violet-700",
  "text-fuchsia-700",
  "text-cyan-700",
  "text-lime-700",
  "text-red-700",
  "text-emerald-700",
  "text-orange-700",
];

const senderBorderColorClasses = [
  "border-l-blue-600",
  "border-l-violet-600",
  "border-l-fuchsia-600",
  "border-l-cyan-600",
  "border-l-lime-600",
  "border-l-red-600",
  "border-l-emerald-600",
  "border-l-orange-600",
];

function getColorIndexForSender(senderKey: string): number {
  let hash = 0;
  for (let i = 0; i < senderKey.length; i++) {
    hash = (hash * 31 + senderKey.charCodeAt(i)) >>> 0;
  }
  return hash % senderTextColorClasses.length;
}

export function ChatPanel({
  chatHistory,
  currentPlayerId,
  onSend,
}: ChatPanelProps) {
  const [draft, setDraft] = useState("");
  const listRef = useRef<HTMLDivElement | null>(null);
  const lastSeenMessageKeyRef = useRef<string>("");

  const lastMessage = chatHistory[chatHistory.length - 1];
  const lastMessageKey = lastMessage
    ? `${lastMessage.sent_at_unix}-${lastMessage.sender_id}-${lastMessage.text}`
    : "";

  useEffect(() => {
    const listEl = listRef.current;
    if (!listEl || !lastMessageKey) {
      return;
    }

    if (lastSeenMessageKeyRef.current === lastMessageKey) {
      return;
    }

    const shouldAnimate = lastSeenMessageKeyRef.current !== "";
    listEl.scrollTo({
      top: listEl.scrollHeight,
      behavior: shouldAnimate ? "smooth" : "auto",
    });
    lastSeenMessageKeyRef.current = lastMessageKey;
  }, [lastMessageKey]);

  const handleSend = () => {
    const text = draft.trim();
    if (!text) {
      return;
    }

    onSend(text);
    setDraft("");
  };

  return (
    <div className="bg-orange-100/90 border-4 border-black rounded-xl shadow-[6px_6px_0px_black] p-4">
      <h3 className="font-propaganda text-lg tracking-wider mb-3">CHAT</h3>

      <div
        ref={listRef}
        className="h-56 overflow-y-auto bg-white/90 border-3 border-black rounded-lg p-3 space-y-2"
      >
        {chatHistory.length === 0 ? (
          <p className="text-sm font-propaganda text-gray-600 tracking-wide">
            No messages yet.
          </p>
        ) : (
          chatHistory.map((entry, index) => {
            const senderKey = entry.sender_id || entry.sender_name || "unknown";
            const colorIndex = getColorIndexForSender(senderKey);

            return (
              <div
                key={`${entry.sent_at_unix}-${index}`}
                className={`text-sm bg-gray-50 border border-gray-200 border-l-4 rounded-md px-2 py-1 ${senderBorderColorClasses[colorIndex]}`}
              >
                <span
                  className={`font-propaganda font-bold tracking-wider ${senderTextColorClasses[colorIndex]}`}
                >
                  {entry.sender_id !== "" && entry.sender_id === currentPlayerId
                    ? "YOU"
                    : entry.sender_name.toUpperCase()}
                  {"  "}
                </span>{" "}
                <span className="font-propaganda tracking-wide text-gray-800">
                  {entry.text}
                </span>
              </div>
            );
          })
        )}
      </div>

      <div className="mt-3 flex min-w-0 gap-2">
        <input
          type="text"
          value={draft}
          onChange={(event) => setDraft(event.target.value)}
          onKeyDown={(event) => {
            if (event.key === "Enter") {
              event.preventDefault();
              handleSend();
            }
          }}
          maxLength={500}
          placeholder="Say something..."
          className="min-w-0 flex-1 bg-white border-3 border-black rounded-lg px-3 py-2 font-propaganda tracking-wide text-black focus:outline-none"
        />
        <button
          type="button"
          onClick={handleSend}
          className="shrink-0 bg-orange-600 hover:bg-orange-700 text-black px-4 py-2 rounded-lg border-3 border-black shadow-[3px_3px_0px_black] font-propaganda font-bold tracking-wider transition-transform duration-150 hover:scale-105"
        >
          SEND
        </button>
      </div>
    </div>
  );
}
