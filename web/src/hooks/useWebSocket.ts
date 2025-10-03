import { useEffect, useRef, useState, useCallback } from "react";

type Options = {
  onOpen?: (event: Event) => void;
  onClose?: (event: CloseEvent) => void;
  onError?: (event: Event) => void;
  onMessage?: (event: MessageEvent) => void;
  reconnect?: boolean; // auto-reconnect on close?
  reconnectInterval?: number; // ms
  deps?: any[];
};

export function useWebSocket(url: string, options: Options = {}) {
  const {
    onOpen,
    onClose,
    onError,
    onMessage,
    reconnect = true,
    reconnectInterval = 2000,
    deps = [],
  } = options;

  const ws = useRef<WebSocket | null>(null);
  const reconnectTimeout = useRef<number | null>(null);

  // Store the latest callbacks in refs to avoid re-creating the connection on each render
  const onOpenRef = useRef<typeof onOpen>(onOpen);
  const onCloseRef = useRef<typeof onClose>(onClose);
  const onErrorRef = useRef<typeof onError>(onError);
  const onMessageRef = useRef<typeof onMessage>(onMessage);

  useEffect(() => {
    onOpenRef.current = onOpen;
  }, [onOpen]);

  useEffect(() => {
    onCloseRef.current = onClose;
  }, [onClose]);

  useEffect(() => {
    onErrorRef.current = onError;
  }, [onError]);

  useEffect(() => {
    onMessageRef.current = onMessage;
  }, [onMessage]);

  const [isConnected, setIsConnected] = useState(false);

  const connect = useCallback(() => {
    console.log(`attempting to connect to ${url}`);
    ws.current = new WebSocket(url);
    console.log("web socket created", ws.current);

    ws.current.onopen = (event) => {
      console.log("open", event);
      setIsConnected(true);
      onOpenRef.current?.(event);
    };

    ws.current.onclose = (event) => {
      console.log(event);
      setIsConnected(false);
      onCloseRef.current?.(event);
      if (reconnect) {
        reconnectTimeout.current = setTimeout(connect, reconnectInterval);
      }
    };

    ws.current.onerror = (event) => {
      console.log(event);
      onErrorRef.current?.(event);
    };

    ws.current.onmessage = (event) => {
      onMessageRef.current?.(event);
    };
  }, [url, reconnect, reconnectInterval, ...deps]);

  useEffect(() => {
    connect();
    return () => {
      ws.current?.close();
      if (reconnectTimeout.current) {
        clearTimeout(reconnectTimeout.current);
        reconnectTimeout.current = null;
      }
    };
  }, [connect]);

  const sendMessage = useCallback((message: string | object) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(
        typeof message === "string" ? message : JSON.stringify(message)
      );
    }
  }, []);

  return { sendMessage, isConnected };
}
