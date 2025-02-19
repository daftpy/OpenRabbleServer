import { Button, Flex, Heading, Text, TextField } from "@radix-ui/themes"
import type { Route } from "../+types/root"
import { useEffect } from "react"
import Keycloak from "keycloak-js";

export function Dash() {
  useEffect(() => {
    console.log("Initializing Keycloak...");
    const keycloak = new Keycloak({
      url: "https://keycloak.localhost",
      realm: "Chatserver",
      clientId: "WebClient",
    });

    keycloak.init({ onLoad: "check-sso", checkLoginIframe: false }).then((authenticated: boolean) => {
      console.log("Keycloak initialized, authenticated:", authenticated);

      if (!authenticated) {
        console.log("User not authenticated, redirecting to login...");
        keycloak.login(); // Redirect to Keycloak login
      }
    }).catch((err) => {
      console.error("Keycloak initialization failed:", err);
    });
  }, []);
  return (
    <main className="p-4">
      <Flex direction="column" gap="3">
        <div>
          <Heading className="font-bold text-xl">Your OnRabble Server</Heading>
        </div>
        <Text>Welcome to your dashboard.</Text>
        <div>
          <Heading className="font-bold" color="indigo" style={{ color: "var(--indigo-9)"}}>Channels</Heading>
          <Text m="0">You can add a new channel or manage your channels below.</Text>
        </div>
        <Flex direction="row" gap="4">
          <TextField.Root placeholder="Add a new channel">
          </TextField.Root>
          <Button>Add</Button>
        </Flex>
      </Flex>
    </main>
  )
}