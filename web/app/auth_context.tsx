import React, { createContext, useEffect, useState } from "react";
import { useNavigate } from "react-router";
import Keycloak from "keycloak-js";

interface AuthContextType {
    keycloak: Keycloak | null;
    authenticated: boolean,
}

const AuthContext = createContext<AuthContextType>({
    keycloak: null,
    authenticated: false,
});

export const AuthProvider = ({ children } : { children: React.ReactNode }) => {
    const [keycloak, setKeycloak] = useState<Keycloak | null>(null);
    const [authenticated, setAuthenticated] = useState(false);
    const navigate = useNavigate();

    useEffect(() => {
      const kc = new Keycloak({
        url: "https://keycloak.localhost",
        realm: "Chatserver",
        clientId: "WebClient", // Change if needed for different clients
      });
      console.log("Hello from auth provider");
      kc.init({ onLoad: "check-sso", checkLoginIframe: false })
        .then((auth) => {
          setAuthenticated(auth);
          setKeycloak(kc);
          console.log("AuthProvider", auth);
          console.log("kc initialized");
          if (!auth) {
            console.log("Logging in");
            kc.login();
          } else {
            // Check for the required role
            const roles = kc.tokenParsed?.realm_access?.roles || [];
            if (!roles.includes("admin")) {
              console.log("User does not have admin role, redirecting...");
              navigate("/unauthorized");
            }
            console.log("You have the right roles. Welcome")
            // Create websocket connection
            // Establish a websocket connection
            // const ws = new WebSocket('wss://chat.localhost/ws?token=' + keycloak.token);
            // ws.onopen = () => {
            //   console.log('WebSocket connection established');
            // };
            // ws.onmessage = (event) => {
            //   const data = JSON.parse(event.data);
            //   // Process incoming messages
            //   console.log(data);
            // };
            // ws.onerror = (error) => console.error('WebSocket error:', error);
            // ws.onclose = () => console.log('WebSocket connection closed');
          }
        })
        .catch((err) => console.error("Keycloak initialization failed:", err));
    }, [navigate]);
    return (
      <AuthContext.Provider value={{ keycloak: keycloak, authenticated: authenticated }}>
        {children}
      </AuthContext.Provider>
    );
}

export default AuthContext;