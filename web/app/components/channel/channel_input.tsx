import { Button, Flex, TextField } from "@radix-ui/themes"
import React, { useState } from "react"
import type { Channel } from "./channel_list";
import { PlusIcon } from "@radix-ui/react-icons";

interface props {
  channelList: Channel[];
  setChannelList: React.Dispatch<React.SetStateAction<Channel[]>>;
}

export default function ChannelInput({ channelList, setChannelList } : props) {
  const [newChannel, setNewChannel] = useState(""); // Input field state
  const [newDescription, setNewDescription] = useState("");

  const hostname = import.meta.env.VITE_HOSTNAME;

  // Function to handle adding a new channel
  const addChannel = async () => {
    if (!newChannel.trim()) return; 

    try {
      const response = await fetch(`https://chat.${hostname}/channels`, {
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
          <Button onClick={addChannel} style={{ boxShadow: "var(--shadow-3)" }} color="jade"><PlusIcon/>Add</Button> {/* ✅ Calls addChannel */}
        </Flex>
    </>
  )
}