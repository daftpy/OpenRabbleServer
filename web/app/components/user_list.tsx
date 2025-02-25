import { Flex, Heading, Text } from "@radix-ui/themes"
import { useContext, useEffect, useState, type JSX } from "react"
import type { ServerMessage } from "~/messages";
import { emitter } from "~/root";
import { WebSocketContext } from "~/websocket_context"
import { PersonIcon } from "@radix-ui/react-icons";

const User = ({ username } : { username: string }) => (
  <Flex justify={"center"} align={"center"} gap={"2"}>
    <PersonIcon /><Text>{ username }</Text>
  </Flex>
);

export default function UserList() {
  const [usernames, setUsernames] = useState<string[]>([]);

  useEffect(() => {
    const handler = (message: ServerMessage) => {
      if (message.type === "connected_users" && message.users) {
        // Replace the whole list of connected users.
        setUsernames(message.users);
      } else if (message.type === "user_status") {
        // If user_status indicates the user went offline, remove them.
        if (!message.status) {
          setUsernames((prev) => prev.filter((u) => u !== message.username));
        } else {
          // Optionally, add the user if they're online and not already in the list.
          setUsernames((prev) =>
            prev.includes(message.username) ? prev : [...prev, message.username]
          );
        }
      }
    };

    // Listen to both events.
    emitter.on("connected_users", handler);
    emitter.on("user_status", handler);

    return () => {
      emitter.off("connected_users", handler);
      emitter.off("user_status", handler);
    };
  }, []);

  return (
      <Flex direction={"column"} minWidth={"150px"} pl={"4"} pt={"4"} align={"center"} gap={"2"}>
        <Heading size={"2"}>Connected Users</Heading>
        {usernames.map((username) => (
          <User key={username} username={username} />
        ))}
      </Flex>
  )
}