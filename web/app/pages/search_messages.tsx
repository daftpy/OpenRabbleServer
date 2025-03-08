import { Cross1Icon, MagnifyingGlassIcon } from "@radix-ui/react-icons";
import { Badge, Box, Button, DropdownMenu, Flex, Heading, Text, TextField } from "@radix-ui/themes";
import { useState } from "react";
import { Link, useFetcher } from "react-router";
import { MessageList } from "~/components/message/message_list";
import { MessageSearchInput } from "~/components/message/search_input";

export function SearchMessagesPage({ messages }: any) {
  const [filteredMessages, setFilteredMessages] = useState<any>(null);

  const handleMessageUpdate = (messages: any) => {
    setFilteredMessages(messages);
  }

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

    </Flex>
  );
}
