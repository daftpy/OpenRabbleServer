import React, { createContext, useEffect, useState } from "react";
import type { ServerMessage } from "./types";

interface WebSocketContextType {
  socket: WebSocket | null;
}

const WebSocketContext = createContext<WebSocketContextType>({
  socket: null
});

export const WebSocketProvider = ({ token, children } : { token: string | undefined, children: React.ReactNode }) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const hostname = import.meta.env.VITE_HOSTNAME;
  const [messages, setMessages] = useState<ServerMessage[]>([]);

  const handleMessage = (message: ServerMessage) => {
    // add the new message to our state:
    setMessages((prevMessages) => [...prevMessages, message]);
    // additional logic for system messages, etc.
    console.log("Message added");
  };

  useEffect(() => {
    if (token == undefined) {
      return;
    }

    const ws = new WebSocket(`wss://chat.${hostname}/ws?token=${token}`);

    ws.onopen = () => {
        console.log("WebSocket connection established");
    }

    ws.onclose = () => {
        console.log("WebSocket connection closed");
    };

    ws.onmessage = (ev: MessageEvent<any>) => {
      console.log("Message Received:", ev.data);
      try {
        const parsedMessage = JSON.parse(ev.data) as ServerMessage;
        handleMessage(parsedMessage);
      } catch (err) {
        console.error("Error parsing message:", err);
      }
    };

    setSocket(ws);

    return () => {
      ws.close();
      setSocket(null);
    }
  }, [token]);

  return (
    <WebSocketContext.Provider value={{socket}}>
      { children }
    </WebSocketContext.Provider>
  )
}