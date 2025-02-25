/*
    The home page of the admin dashboard. Landing here will display the current channels
    of the chatserver as well as connected users.
*/
import { useEffect, useState } from "react"
import { useNavigate } from "react-router";
import ChannelInput from "~/components/channel_input";
import ChannelList from "~/components/channel_list";
import AuthContext from "~/auth_context";
import { useContext } from "react";

import { Button, Container, Flex, Heading, Text } from "@radix-ui/themes";
import type { Channel } from "~/components/channel_list";


/*
  TODO: Added a test navigate button here to move between pages. It works properly
  and does not accidentally trigger a refresh of the auth provider. Perfect!
*/

export function HomePage({ channels }: { channels: Channel[] }) {
  const [channelList, setChannelList] = useState<Channel[]>(channels);
  const { keycloak, authenticated } = useContext(AuthContext);
  const [ws, setWs] = useState<WebSocket | null>(null);
  const navigate = useNavigate();

  // Establish WebSocket connection when authenticated and keycloak are available.
  // useEffect(() => {
  //   if (authenticated && keycloak) {
  //     const token = keycloak.token;
  //     const socket = new WebSocket(`wss://chat.localhost/ws?token=${token}`);
      
  //     socket.onopen = () => {
  //       console.log("WebSocket connection established");
  //     };

  //     socket.onmessage = (event) => {
  //       const data = JSON.parse(event.data);
  //       console.log("Message received:", data);
  //       // Process incoming messages, e.g. update channel list or user status
  //     };

  //     socket.onerror = (error) => {
  //       console.error("WebSocket error:", error);
  //     };

  //     socket.onclose = () => {
  //       console.log("WebSocket connection closed");
  //     };

  //     setWs(socket);

  //     // Cleanup function to close the socket when component unmounts or dependencies change
  //     return () => {
  //       socket.close();
  //     };
  //   }
  // }, [authenticated, keycloak]);

    return (
      <main className="p-4">
        <Container>
          <Flex direction="column" gap="3">
            <div>
              <Button onClick={() => navigate("/about")}>About</Button>
              <Heading weight={"bold"} className="text-xl">Your OnRabble Server</Heading>
            </div>
            <Text>Welcome to your dashboard.</Text>
            <div>
              <Heading weight={"bold"} style={{ color: "var(--indigo-9)" }}>Channels</Heading>
              <Text m="0">You can add a new channel or manage your channels below.</Text>
            </div>
            <ChannelInput channelList={channelList} setChannelList={setChannelList} />
            <ChannelList channels={channelList} />
          </Flex>
        </Container>
      </main>   
    )
}
