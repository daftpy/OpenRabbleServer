import { Cross1Icon, MagnifyingGlassIcon } from "@radix-ui/react-icons";
import { Badge, Box, Button, DropdownMenu, Flex, Heading, Text, TextField } from "@radix-ui/themes";
import { useState } from "react";
import { Link, useFetcher } from "react-router";
import { MessageList } from "~/components/message/message_list";

export function SearchMessagesPage({ messages, channels }: any) {
  const fetcher = useFetcher();
  const [filters, setFilters] = useState<string[]>([]);
  const [keyword, setKeyword] = useState<string>("");

  const addFilter = (filter: string) => {
    setFilters((prev) => [...prev, filter]);
  };

  const removeFilter = (targetFilter: string) => {
    setFilters((prev) => prev.filter((filter) => filter !== targetFilter));
  };

  const handleSearch = () => {
    let formData = new FormData();
    formData.append("keyword", keyword);
    filters.forEach((filter) => formData.append("channel", filter));
    fetcher.submit(formData, { method: "post" });
  };

  return (
    <Flex direction={"column"} maxWidth={"990px"} m={"auto"} px={"4"} py={"6"}>
      <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{ color: "var(--indigo-9)" }}>
        <Link to="/">Your OnRabble Server</Link>
      </Heading>
      <Box pt={"4"}>
        <Heading style={{color: "var(--subheading-color)"}}>Search</Heading>
        <Text>Search messages stored in your cache or database. Filter by <strong style={{ color: "var(--link-color)" }}>channel </strong>or <strong style={{ color: "var(--link-color)" }}>keyword</strong>.</Text>
      </Box>
      {/* Search Input & Filters */}
      <Flex gap={"4"} py={"2"}>
        <TextField.Root
          placeholder="keyword"
          className="grow"
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
        />
        <DropdownMenu.Root>
          <DropdownMenu.Trigger>
            <Button color="amber">Filter</Button>
          </DropdownMenu.Trigger>
          <DropdownMenu.Content>
            {channels.map((channel: any) => (
              <DropdownMenu.Item key={channel.name} onClick={() => addFilter(channel.name)}>
                {channel.name}
              </DropdownMenu.Item>
            ))}
          </DropdownMenu.Content>
        </DropdownMenu.Root>
        <Button color="blue" onClick={handleSearch}><MagnifyingGlassIcon />Search</Button>
      </Flex>

      {/* Active Filters */}
      <Flex gap={"2"} align={"center"} pb={"4"}>
        <Text weight={"bold"} style={{ color: "var(--indigo-12)" }}>Filters: </Text>
        {filters.map((filter) => (
          <Badge key={filter} size={"2"} color="tomato" onClick={() => removeFilter(filter)}>
            {filter} <Cross1Icon />
          </Badge>
        ))}
      </Flex>

      {/* Show messages - if fetcher is pending, show loading */}
      {fetcher.state === "submitting" || fetcher.state === "loading" ? (
        <Text>Loading messages...</Text>
      ) : (
        <MessageList messages={fetcher.data?.messages ?? messages} hidePermaLink={true} />
      )}
    </Flex>
  );
}
