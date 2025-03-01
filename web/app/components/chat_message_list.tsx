import { PersonIcon } from "@radix-ui/react-icons";
import { Box, Flex, Heading, ScrollArea, Text } from "@radix-ui/themes";
import { useEffect, useRef, useState } from "react";
import type { ServerMessage } from "~/messages";
import { emitter } from "~/root";

export interface ChatMessageType {
  username: string;
  channel: string;
  message: string;
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
        <Text size="1" weight={"bold"} color="indigo">{ channel } #</Text>
        <Flex gap={"1"} align={"center"}><PersonIcon /><Heading size="1">{username}</Heading></Flex>
      </Flex>
      <Text size="1" wrap={"wrap"}>{content}</Text>
    </Flex>
  );
};

export default function ChatMessageList() {
  const [messages, setMessages] = useState<ChatMessageType[]>([]);
  const bottomRef = useRef<HTMLDivElement>(null);
  const topRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    console.log("loaded");
    const handler = (message: ServerMessage) => {
      if (message.type === "chat_message") {
        console.log("channel", message.channel);
        setMessages((prev) => [
          ...prev,
          { username: message.username, channel: message.channel, message: message.message },
        ]);
      } else if (message.type === "bulk_chat_messages") {
        setMessages((prev) => [
          ...prev,
          ...message.messages // Spread the array correctly here
        ]);
        console.log("messages added");
      }
    };

    emitter.on("chat_message", handler);
    emitter.on("bulk_chat_messages", handler);
    return () => {
      emitter.off("chat_message", handler);
      emitter.off("bulk_chat_messages", handler);
    };
  }, []);

  // Whenever messages change, scroll to the bottom.
  useEffect(() => {
    // If you want a smooth scroll, you can do:
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // Observe topRef visibility
  useEffect(() => {
    if (!topRef.current) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          console.log("Top of the message list is now visible!");
        }
      },
      { root: null, threshold: 0.1 } // Adjust threshold as needed
    );

    observer.observe(topRef.current);

    return () => observer.disconnect();
  }, []);

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
          <div ref={topRef} />
          {messages.map((message, index) => (
            <Message
              key={index}
              username={message.username}
              channel={message.channel}
              content={message.message}
              isLast={index == messages.length - 1 ? true : false}
            />
          ))}
          <div ref={bottomRef} />
        </Flex>
      </ScrollArea>
    </Box>
  );
}
