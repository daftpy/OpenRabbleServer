import React, { createContext, useContext, useEffect, useState } from "react";
import type { ServerMessage } from "../messages";
import { emitter } from "~/root";
import type { ChatMessageType } from "~/components/message/live_view";

interface WebSocketContextType {
  socket: WebSocket | null;
  messages: ChatMessageType[];
}

// How many messages the WebSocketProvider holds. Cyclical bufferr
const MAX_MESSAGES = 100;

export const WebSocketContext = createContext<WebSocketContextType>({
  socket: null,
  messages: [],
});

export const WebSocketProvider = ({ token, children } : { token: string | undefined, children: React.ReactNode }) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const hostname = import.meta.env.VITE_HOSTNAME;
  const [messages, setMessages] = useState<ChatMessageType[]>([]);

  const addMessage = (msg: ChatMessageType) => {
    setMessages((prev) => {
      const newMessages = [...prev, msg]; 
      return newMessages.length > MAX_MESSAGES ? newMessages.slice(-MAX_MESSAGES) : newMessages;
    });
  };

  useEffect(() => {
    // If not token, deny
    if (token == undefined) {
      return;
    }

    // Create the websocket connection using the token
    const ws = new WebSocket(`wss://chat.${hostname}/ws?token=${token}`);
    // On open
    ws.onopen = () => {
        console.log("WebSocket connection established");
    }
    // On close
    ws.onclose = () => {
        console.log("WebSocket connection closed");
    };

    // On meessage
    ws.onmessage = (ev: MessageEvent<any>) => {
      console.log("Message Received:", ev.data);
      try {
        const message = JSON.parse(ev.data) as ServerMessage;
        
        // TODO: consider removal
        // Emit the message for other parts of the application to react to
        emitter.emit(message.type, message);
        console.log("Emitted:", message.type);

        // Add the message to the circular buffer
        if (message.type === "chat_message") {
          addMessage({id: message.payload.id, username: message.payload.username, channel: message.payload.channel, message: message.payload.message, authored_at: message.payload.authored_at });
        } else if (message.type === "bulk_chat_messages") {
          setMessages((prev) => {
            const newMessages = [...prev, ...message.payload.messages];
            return newMessages.length > MAX_MESSAGES ? newMessages.slice(-MAX_MESSAGES) : newMessages;
          });
        }
      } catch (err) {
        console.error("Error parsing message:", err);
      }
    };

    // Set the socket
    setSocket(ws);
    return () => {
      ws.close();
      setSocket(null);
    }
  }, [token]);

  return (
    <WebSocketContext.Provider value={{socket: socket, messages: messages}}>
      { children }
    </WebSocketContext.Provider>
  )
}

export function useWebSocket() {
  const context = useContext(WebSocketContext);
  if (!context) throw new Error("useWebSocket must be used within WebSocketProvider");
  return context;
}
