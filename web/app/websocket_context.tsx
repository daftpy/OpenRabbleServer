import React, { createContext, useEffect, useState } from "react";

interface WebSocketContextType {
  socket: WebSocket | null;
}

const WebSocketContext = createContext<WebSocketContextType>({
  socket: null
});

export const WebSocketProvider = ({ token, children } : { token: string | undefined, children: React.ReactNode }) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);

  useEffect(() => {
    if (token == undefined) {
      return;
    }

    const ws = new WebSocket(`wss://chat.localhost/ws?token=${token}`);

    ws.onopen = () => {
        console.log("WebSocket connection established");
    }

    ws.onclose = () => {
        console.log("WebSocket connection closed");
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