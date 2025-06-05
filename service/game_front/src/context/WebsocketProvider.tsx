import { useEffect, useRef, useState, createContext, useContext, JSX } from "react";
import { getCookie } from "../utils/utils";
import { WS_URI } from "../api/api";
import { useAuth } from "./AuthProvider";
import { WSEvent } from "../models/events";

interface WSContextType {
    connected: boolean;
    wsRef: React.RefObject<WebSocket | null>;
    addMessageListener: (listener: (data: any) => void) => void;
    removeMessageListener: (listener: (data: any) => void) => void;
}

const WSContext = createContext<WSContextType | undefined>(undefined);

export const WebSocketProvider = ({ children }: { children: JSX.Element }) => {
  const { user } = useAuth();
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectRef = useRef<number | null>(null);
  const [connected, setConnected] = useState(false);
  const listenersRef = useRef<((data: any) => void)[]>([]);

  const connectWS = () => {
    const token = getCookie("token");
    if (!user || !token) return;

    const ws = new WebSocket(WS_URI);
    wsRef.current = ws;

    ws.onopen = () => {
        console.log("WebSocket открыт");
        ws.send(JSON.stringify({ token }));
        setConnected(true);
    };

    ws.onmessage = (ev) => {
        try {
            if (ev.data == "") return; // Пинг от сервера игнорируем
            const data: WSEvent = JSON.parse(ev.data);
            console.log("Получено сообщение:", data);
            listenersRef.current.forEach((listener) => listener(data));
        } catch (err) {
            console.error("Ошибка парсинга:", err);
        }
    };

    ws.onclose = () => {
      console.warn("WebSocket закрыт");
      setConnected(false);
      listenersRef.current = [];
      if (reconnectRef.current === null) {
            reconnectRef.current = setTimeout(() => {
            reconnectRef.current = null;
            connectWS();
        }, 5000);
      }
    };

    ws.onerror = (err) => {
      console.error("WebSocket ошибка:", err);
      ws.close();
    };
  };

  useEffect(() => {
    if (user) connectWS();
    return () => {
        if (wsRef.current) {
            wsRef.current.close();
            wsRef.current = null;
        }
        if (reconnectRef.current) {
            clearTimeout(reconnectRef.current);
            reconnectRef.current = null;
        }
    };
  }, [user]);

  const addMessageListener = (listener: (data: any) => void) => {
    listenersRef.current.push(listener);
  };

  const removeMessageListener = (listener: (data: any) => void) => {
    listenersRef.current = listenersRef.current.filter((l) => l !== listener);
  };

  return (
    <WSContext.Provider value={{
        connected: connected,
        wsRef: wsRef,
        addMessageListener: addMessageListener,
        removeMessageListener: removeMessageListener
    }}>
      {children}
    </WSContext.Provider>
  );
};

export const useWS = () => {
    const context = useContext(WSContext);
    if (!context) {
        throw new Error("useWS must be used within an WebsocketProvider");
    }
    return context;
};
