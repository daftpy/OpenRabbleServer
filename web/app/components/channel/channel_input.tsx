import { Button, Flex, TextField } from "@radix-ui/themes"
import { useState } from "react"
import { PlusIcon } from "@radix-ui/react-icons";
import { useFetcher } from "react-router";

export default function ChannelInput() {
  const [newChannel, setNewChannel] = useState(""); // Input field state
  const [newDescription, setNewDescription] = useState("");
  const channelFetcher = useFetcher();

  // Function to handle adding a new channel
  const addChannel = async () => {
    if (!newChannel.trim()) return;

    channelFetcher.submit(
      {
        name: newChannel,
        description: newDescription,
        intent: "add"
      },
      {
        method: "POST",
        action: "/channels"
      }
    );

    // Clear the inputs
    setNewChannel("");
    setNewDescription("");
  };

  return (
    <>
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
          <Button onClick={addChannel} style={{ boxShadow: "var(--shadow-3)" }} color="jade"><PlusIcon/>Add</Button>
        </Flex>
    </>
  )
}