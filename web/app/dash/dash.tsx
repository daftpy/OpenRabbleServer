import { Button, Flex, Heading, Table, Text, TextField } from "@radix-ui/themes"
import { useEffect, useState } from "react"
import { useNavigate } from "react-router";
import Keycloak from "keycloak-js";

interface Channel {
  name: string;
  description: string | null;
}

export function Dash({ channels }: { channels: Channel[] }) {
  const [channelList, setChannelList] = useState<Channel[]>(channels); // Local state for channels
  const [newChannel, setNewChannel] = useState(""); // Input field state
  const [newDescription, setNewDescription] = useState("");
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
    if (!newChannel.trim()) return; 

    try {
      const response = await fetch("https://chat.localhost/channels", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: newChannel.trim(), description: newDescription.trim() }),
      });

      if (!response.ok) {
        throw new Error("Failed to create channel");
      }

      const newChannelObj: Channel = {
        name: newChannel.trim(),
        description: newDescription.trim() || null,
      };

      setChannelList([...channelList, newChannelObj]); 
      setNewChannel(""); 
      setNewDescription("");
    } catch (error) {
      console.error("Error creating channel:", error);
    }
  };

  return (
    <main className="p-4">
      <Flex direction="column" gap="3">
        <div>
          <Heading weight={"bold"} className="text-xl">Your OnRabble Server</Heading>
        </div>
        <Text>Welcome to your dashboard.</Text>
        <div>
          <Heading weight={"bold"} style={{ color: "var(--indigo-9)" }}>Channels</Heading>
          <Text m="0">You can add a new channel or manage your channels below.</Text>
        </div>
        <Flex direction="row" gap="4">
          <TextField.Root 
            placeholder="Channel Name"
            value={newChannel}
            onChange={(e) => setNewChannel(e.target.value)}
          />
          <TextField.Root
            placeholder="description"
            value={newDescription}
            onChange={(e) => setNewDescription(e.target.value)}
            className="flex-grow"
          >
          </TextField.Root>
          <Button onClick={addChannel} style={{ boxShadow: "var(--shadow-3)" }}>Add</Button> {/* ✅ Calls addChannel */}
        </Flex>
        <Table.Root>
          <Table.Header>
            <Table.Row>
              <Table.ColumnHeaderCell width="125px">Channel</Table.ColumnHeaderCell>
              <Table.ColumnHeaderCell width={"auto"}>Description</Table.ColumnHeaderCell>
              <Table.ColumnHeaderCell></Table.ColumnHeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {channelList.map((channel, index) => (
              <Table.Row key={index}>
                <Table.RowHeaderCell justify="start">{channel.name}</Table.RowHeaderCell>
                <Table.Cell justify={"start"}>{ channel.description ? <>{channel.description}</> : <>...</>}</Table.Cell>
                <Table.Cell justify="end"><Button color="red" size={"1"}  style={{ boxShadow: "var(--shadow-1)" }}>x</Button></Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table.Root>
      </Flex>
    </main>
  )
}
