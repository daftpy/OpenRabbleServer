import { Box, Flex, ScrollArea } from "@radix-ui/themes";
import { useEffect, useRef, useState } from "react";
import type { ServerMessage } from "~/messages";
import { emitter } from "~/root";
import { MessageList } from "./message_list";

// TODO: refactor this out
export interface ChatMessageType {
  username: string;
  channel: string;
  message: string;
}

export function LiveView() {
  const [messages, setMessages] = useState<ChatMessageType[]>([]);
  // Track whether messages have been hydrated
  const [isHydrated, setIsHydrated] = useState(false);
  const [observerEnabled, setObserverEnabled] = useState(false);
  
  // Refs to the bottom and top of the list
  const bottomRef = useRef<HTMLDivElement>(null);
  const topRef = useRef<HTMLDivElement>(null);
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Once messages are loaded, set isHydrated to true
    if (messages.length > 0 && !isHydrated) {
      setIsHydrated(true);

      // Delay enabling the observer slightly to avoid smooth scroll interference
      setTimeout(() => setObserverEnabled(true), 500);
    }

    // Auto-scroll only when new messages arrive (but not on initial load)
    if (isHydrated && scrollRef.current) {
      scrollRef.current.scrollTo({ top: scrollRef.current.scrollHeight, behavior: "smooth" });
    }
  }, [messages, isHydrated]);

  useEffect(() => {
    if (!topRef.current || !observerEnabled) return; // Ensure observer runs only after delay

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          console.log("Top of the message list is now visible!");
          // Fetch older messages here if needed
        }
      },
      { root: null, threshold: 0.1 }
    );

    observer.observe(topRef.current);

    return () => observer.disconnect();
  }, [observerEnabled]); // Runs only after observer is enabled
  
  useEffect(() => {
    console.log("loaded");
    const handler = (message: ServerMessage) => {
      if (message.type === "chat_message") {
        console.log("channel", message.payload.channel);
        setMessages((prev) => [
          ...prev,
          { username: message.payload.username, channel: message.payload.channel, message: message.payload.message },
        ]);
      } else if (message.type === "bulk_chat_messages") {
        setMessages((prev) => [
          ...prev,
          ...message.payload.messages // Spread the array correctly here
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

  return (
    <Box className="rounded-sm p-2" style={{ border: "2px solid var(--indigo-3)" }}>
      <ScrollArea
        className="bg-blue-50 p-2"
        style={{
          minHeight: "50px",
          maxHeight: "150px",
          backgroundColor: "var(--indigo-2)",
        }}
        ref={scrollRef}
      >
        <Flex direction="column" gap="2">
          <div ref={topRef} />
            <MessageList messages={messages} />
          <div ref={bottomRef} />
        </Flex>
      </ScrollArea>
    </Box>
  )
}