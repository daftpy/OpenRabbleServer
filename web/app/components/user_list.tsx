import { Box, Flex, Heading, Text } from "@radix-ui/themes"
import { useEffect, useState } from "react"
import type { ServerMessage } from "~/messages";
import { emitter } from "~/root";
import { PersonIcon } from "@radix-ui/react-icons";

const User = ({ username } : { username: string }) => (
  <Flex justify={"center"} align={"center"} gap={"2"}  maxWidth={"225px"} minWidth={"175px"} overflow={"hidden"}>
    <Box><PersonIcon/></Box><Text truncate>{ username }</Text>
  </Flex>
);

export default function UserList() {
  const [usernames, setUsernames] = useState<string[]>([]);

  useEffect(() => {
    const handler = (message: ServerMessage ) => {
      if (message.type === "connected_users" && message.users) {
        // Replace the whole list of connected users.
        setUsernames(message.users);
      } else if (message.type === "user_status") {
        // If user_status indicates the user went offline, remove them.
        if (!message.payload.status) {
          setUsernames((prev) => prev.filter((u) => u !== message.payload.username));
        } else {
          // Optionally, add the user if they're online and not already in the list.
          setUsernames((prev) =>
            prev.includes(message.payload.username) ? prev : [...prev, message.payload.username]
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
      <Flex direction={"column"} minHeight={"fit-content"}  px={"4"} pt={"2"} align={"center"} gap={"2"}>
        <Heading size={"3"} wrap={"nowrap"} style={{padding: "0 10px", color: "var(--subheading-color)"}}>Connected Users</Heading>
        {usernames.length != 0 ? (
          usernames.map((username) => (
            <User key={username} username={username} />
          ))
        ) : (
        <Text size={"2"}>
          No users connected
        </Text>
      )}
      </Flex>
  )
}