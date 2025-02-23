import type React from "react"
import { useEffect } from "react";
import Keycloak from "keycloak-js";
import { useNavigate } from "react-router";

interface props {
  children: React.ReactNode;
}

export default function RouteProtector({ children } : props) {
  const navigate = useNavigate();
  useEffect(() => {
    console.log("Route Protect log");
    console.log("Initializing Keycloak...");

    const keycloak = new Keycloak({
      url: "https://keycloak.localhost",
      realm: "Chatserver",
      clientId: "WebClient",
    });

    keycloak.init({ onLoad: "check-sso", checkLoginIframe: false }).then((authenticated: boolean) => {
      console.log("Keycloak initialized, authenticated:", authenticated);
      const roles = keycloak.tokenParsed?.realm_access?.roles || [];
      if (!authenticated) {
        console.log("User not authenticated, redirecting to login...");
        keycloak.login(); // Redirect to Keycloak login
      } else if (roles.find((role) => role == "admin") == undefined) {
        console.log(keycloak.tokenParsed);
        navigate("/unauthorized");
      }
      // Establish a websocket connection
      const ws = new WebSocket('wss://chat.localhost/ws?token=' + keycloak.token);
      ws.onopen = () => {
        console.log('WebSocket connection established');
      };
      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        // Process incoming messages
        console.log(data);
      };
      ws.onerror = (error) => console.error('WebSocket error:', error);
      ws.onclose = () => console.log('WebSocket connection closed');
    }).catch((err) => {
      console.error("Keycloak initialization failed:", err);
    });
  }, [navigate])
  return (
    <>
      { children }
    </>
  )
}