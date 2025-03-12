import { Box, Button, Flex, Heading, Text } from "@radix-ui/themes";
import { useEffect, useState } from "react";
import { Link, useFetcher } from "react-router";
import { MessageList } from "~/components/message/message_list";
import { MessageSearchInput } from "~/components/message/search_input";

export function SearchMessagesPage({ messages, hasMore }: {messages: any, hasMore: boolean}) {
  const [page, setPage] = useState<number>(0);
  const messageFetcher = useFetcher();
  const [filteredMessages, setFilteredMessages] = useState<any>(null);
  const [moreMessages, setMoreMessages] = useState<boolean>(hasMore);
  const [keyword, setKeyword] = useState<string>("");
  const [filters, setFilters] = useState<any>([]);

  const handleMessageUpdate = (messages: any, hasMore: boolean, keyword: string, filters: string[]) => {
    setFilteredMessages(messages);
    setMoreMessages(hasMore);
    console.log("has more?", hasMore);
    setKeyword(keyword);
    setFilters(filters);
    setPage(0);
  }

  // Fetch messages from the API
  const fetchMessages = (newPage: number) => {
    const limit = 10; // Adjust
    const offset = newPage * limit;
    
    console.log("Fetching messages with offset:", offset);
  
    // Build a query string for the loader
    const params = new URLSearchParams();
    params.append("keyword", keyword);
    params.append("limit", limit.toString());
    params.append("offset", offset.toString());

    filters.forEach((filter: any) => {
      params.append("channel", filter.name);
    });

    // Trigger the loader by calling fetcher.load
    messageFetcher.load(`/messages?${params.toString()}`);
  
    setPage(newPage);
  };

  useEffect(() => {
    console.log("messageFetcher", messageFetcher.data);
    // Update messages when fetcher gets new data
    if (messageFetcher.data && messageFetcher.data.messages) {
      console.log("updating messages:", messageFetcher.data);
      console.log("has more?", messageFetcher.data.hasMore);
      setFilteredMessages(messageFetcher.data.messages);
      setMoreMessages(messageFetcher.data.hasMore);
    }
  }, [messageFetcher.data])

  return (
    <Flex direction={"column"} maxWidth={"990px"} m={"auto"} px={"4"} py={"6"}>
      <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{ color: "var(--slate-12)" }}>
        <Link to="/">Your OnRabble Server</Link>
      </Heading>
      <Box pt={"4"}>
        <Heading style={{color: "var(--indigo-10)"}}>Search</Heading>
        <Text>Search messages stored in your cache or database. Filter by <strong style={{ color: "var(--link-color)" }}>channel </strong>or <strong style={{ color: "var(--link-color)" }}>keyword</strong>.</Text>
      </Box>
      {/* Search Input & Filters */}
      <MessageSearchInput onMessagesUpdate={handleMessageUpdate} />

      {/* Show messages */}
      <MessageList messages={filteredMessages ? filteredMessages : messages} hidePermaLink={true} />
      <Flex justify={"between"} pt={"4"}>
      <Button 
          disabled={page === 0} 
          onClick={() => fetchMessages(page - 1)}
        >
          Previous Page
        </Button>
        <Button 
          disabled={!moreMessages}
          onClick={() => fetchMessages(page + 1)}
        >
          Next Page
        </Button>
      </Flex>
    </Flex>
  );
}
