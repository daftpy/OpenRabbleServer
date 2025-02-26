import { PersonIcon } from "@radix-ui/react-icons";
import { Box, Flex, Heading, ScrollArea, Text } from "@radix-ui/themes";
import { useEffect, useRef, useState } from "react";
import type { ServerMessage } from "~/messages";
import { emitter } from "~/root";

interface ChatMessageType {
  username: string;
  channel: string;
  content: string;
}

const Message = ({ username, channel, content, isLast }: { username: string, channel: string, content: string, isLast: boolean }) => {
  let borderStyle;
  if (isLast) {
    borderStyle = "none";
  } else {
    borderStyle = "2px solid var(--indigo-3)";
  }
  console.log("message channel", channel);
  return (
    <Flex direction="column" gap="1" style={{borderBottom: borderStyle}} p={"1"} pb={"2"}>
      
      <Flex gap={"1"} align={"center"}>
        <Flex gap={"1"} align={"center"}><PersonIcon /><Heading size="1">{username}</Heading></Flex>
        <Text size="1" weight={"bold"} color="indigo">{ channel } #</Text>
      </Flex>
      <Text size="1" wrap={"wrap"}>{content}</Text>
    </Flex>
  );
};

export default function ChatMessageList() {
  const [messages, setMessages] = useState<ChatMessageType[]>([]);
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handler = (message: ServerMessage) => {
      if (message.type === "chat_message") {
        console.log("channel", message.channel);
        setMessages((prev) => [
          ...prev,
          { username: message.username, channel: message.channel, content: message.message },
        ]);
      }
    };

    emitter.on("chat_message", handler);
    return () => {
      emitter.off("chat_message", handler);
    };
  }, []);

  // Whenever messages change, scroll to the bottom.
  useEffect(() => {
    // If you want a smooth scroll, you can do:
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  return (
    <Box className="rounded-sm border-solid border-2 border-gray-100 p-2">
      <ScrollArea
        className="bg-blue-50 p-2"
        style={{
          minHeight: "50px",
          maxHeight: "150px",
          backgroundColor: "var(--indigo-2)",
        }}
      >
        <Flex direction="column" gap="2">
          {messages.map((message, index) => (
            <Message
              key={index}
              username={message.username}
              channel={message.channel}
              content={message.content}
              isLast={index == messages.length - 1 ? true : false}
            />
          ))}
          <div ref={bottomRef} />
        </Flex>
      </ScrollArea>
    </Box>
  );
}
