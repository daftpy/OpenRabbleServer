import { Button, Flex, Heading, Table, Text, TextField } from "@radix-ui/themes"
import { useEffect, useState } from "react"
import { useNavigate } from "react-router";
import Keycloak from "keycloak-js";

export function Dash({ channels }: { channels: string[] }) {
  const [channelList, setChannelList] = useState(channels); // Local state for channels
  const [newChannel, setNewChannel] = useState(""); // Input field state
  const navigate = useNavigate();
  useEffect(() => {
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
        navigate("/unauthorized");
      }
    }).catch((err) => {
      console.error("Keycloak initialization failed:", err);
    });
  }, [navigate]);

  console.log("channels", channels);

  // Function to handle adding a new channel
  const addChannel = async () => {
    if (!newChannel.trim()) return; // Prevent empty channel names

    try {
      const response = await fetch("https://chat.localhost/channels", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: newChannel.trim() }),
      });

      if (!response.ok) {
        throw new Error("Failed to create channel");
      }

      setChannelList([...channelList, newChannel]); // Update UI with new channel
      setNewChannel(""); // Clear input field
    } catch (error) {
      console.error("Error creating channel:", error);
    }
  };

  return (
    <main className="p-4">
      <Flex direction="column" gap="3">
        <div>
          <Heading className="font-bold text-xl">Your OnRabble Server</Heading>
        </div>
        <Text>Welcome to your dashboard.</Text>
        <div>
          <Heading className="font-bold" color="indigo" style={{ color: "var(--indigo-9)" }}>Channels</Heading>
          <Text m="0">You can add a new channel or manage your channels below.</Text>
        </div>
        <Flex direction="row" gap="4">
          <TextField.Root 
            placeholder="Add a new channel"
            value={newChannel}
            onChange={(e) => setNewChannel(e.target.value)}
          />
          <Button onClick={addChannel}>Add</Button> {/* ✅ Calls addChannel */}
        </Flex>
        <Table.Root>
          <Table.Header>
            <Table.Row>
              <Table.ColumnHeaderCell>Channel</Table.ColumnHeaderCell>
              <Table.ColumnHeaderCell></Table.ColumnHeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {channelList.map((channel, index) => (
              <Table.Row key={index}>
                <Table.RowHeaderCell justify="start">{channel}</Table.RowHeaderCell>
                <Table.Cell justify="end"><Button color="red">remove</Button></Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table.Root>
      </Flex>
    </main>
  )
}
