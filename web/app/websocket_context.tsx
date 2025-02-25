import React, { createContext, useEffect, useState } from "react";
import type { ServerMessage } from "./messages";
import { emitter } from "./root";

interface WebSocketContextType {
  socket: WebSocket | null;
  messages: ServerMessage[];
  systemMessages: ServerMessage[];
}

export const WebSocketContext = createContext<WebSocketContextType>({
  socket: null,
  messages: [],
  systemMessages: [],
});

export const WebSocketProvider = ({ token, children } : { token: string | undefined, children: React.ReactNode }) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const hostname = import.meta.env.VITE_HOSTNAME;
  const [messages, setMessages] = useState<ServerMessage[]>([]);
  const [systemMessages, setSystemMessages] = useState<ServerMessage[]>([]);

  const handleMessage = (message: ServerMessage) => {
    // add the new message to our state:
    if (message.type === "connected_users" || message.type === "user_status" || message.type === "active_channels") {
      setSystemMessages((prevMessages) => [...prevMessages, message]);
      console.log("System message added:", message.type);
    } else {
      setMessages((prevMessages) => [...prevMessages, message]);
      console.log("Message added");
    }
    // additional logic for system messages, etc.

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
        // handleMessage(parsedMessage);
        emitter.emit(parsedMessage.type, parsedMessage);
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
    <WebSocketContext.Provider value={{socket: socket, messages: messages, systemMessages: systemMessages}}>
      { children }
    </WebSocketContext.Provider>
  )
}
