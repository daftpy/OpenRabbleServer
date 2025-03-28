import { Box, Flex, ScrollArea } from "@radix-ui/themes";
import { useEffect, useRef, useState } from "react";
import { MessageList } from "./message_list";
import { useWebSocket } from "~/contexts/websocket_context";

// TODO: refactor this out
export interface ChatMessageType {
  username: string;
  channel: string;
  message: string;
  authored_at: string;
  id: number; // These do not actually exist on the ChatMessageType (live messages) needs a serious refactor
  owner_id: string; // These do not actually exist on the ChatMessageType (live messages) needs a serious refactor
}

export function LiveView() {
  // const [messages, setMessages] = useState<ChatMessageType[]>([]);
  const { messages } = useWebSocket();
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

  return (
    <Box className="rounded-sm p-2" style={{ border: "2px solid var(--indigo-3)" }}>
      <ScrollArea
        className="bg-blue-50 p-2"
        style={{
          minHeight: "50px",
          maxHeight: "190px",
          backgroundColor: "var(--indigo-2)",
        }}
        ref={scrollRef}
      >
        <Flex direction="column" gap="2" px={"2"}>
          <div ref={topRef} />
            <MessageList messages={messages} hidePermaLink={false} />
          <div ref={bottomRef} />
        </Flex>
      </ScrollArea>
    </Box>
  )
}