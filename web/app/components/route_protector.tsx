import type React from "react"
import { useContext } from "react";
import AuthContext from "~/contexts/auth_context";

interface props {
  children: React.ReactNode;
}

export default function RouteProtector({ children } : props) {
  const { keycloak, authenticated } = useContext(AuthContext);
  if (!authenticated || !keycloak) {
    return <div>Loading...</div>;
  }
  return (
    <>
      { children }
    </>
  )
}