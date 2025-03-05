import { Cross1Icon } from "@radix-ui/react-icons";
import { Badge, Button, DropdownMenu, Flex, Heading, Text, TextField } from "@radix-ui/themes";
import { useState } from "react";
import { Link } from "react-router";
import { MessageList } from "~/components/message/message_list";

export function SearchMessagesPage({ messages, channels } : any) {
  const [filters, setFilters] = useState<string[]>([]);

  const addFilter = (filter : string) => {
    setFilters((prev) => [...prev, filter]);
  }

  const removeFilter = (targetFilter : string) => {
    setFilters((prev) => prev.filter((filter : string) => filter !== targetFilter));
  }

  return (
    <Flex direction={"column"} maxWidth={"990px"} m={"auto"} px={"4"} py={"6"}>
      <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{color: "var(--indigo-9)"}}>
        <Link to="/">Your OnRabble Server</Link>
      </Heading>
      <Text>Search messages stored in your cache or database. Filter by <strong style={{color: "var(--link-color)"}}>channel </strong>or <strong style={{color: "var(--link-color)"}}>keyword</strong>.</Text>
      <Flex py={"4"} gap={"4"}>
      <TextField.Root 
        placeholder="keyword"
        className="grow"
      />
      <DropdownMenu.Root>
        <DropdownMenu.Trigger>
          <Button>
            Channel
          </Button>
        </DropdownMenu.Trigger>
        <DropdownMenu.Content>
          {channels.map((channel : any) => (
            <DropdownMenu.Item onClick={() => addFilter(channel.name)}>{ channel.name }</DropdownMenu.Item>
          ))}
        </DropdownMenu.Content>
      </DropdownMenu.Root>
      <Button color="grass">Search</Button>
      </Flex>
      <Flex pb={"2"} gap={"2"} align={"center"}>
        <Text weight={"bold"} style={{color: "var(--subheading-color)"}}>Filters: </Text>
        {filters && filters.map((filter : any) => (
          <Badge size={"2"} color="ruby" onClick={() => removeFilter(filter)}>{ filter }<Cross1Icon /></Badge>
        ))}
      </Flex>
      <MessageList messages={messages} />
    </Flex>
  )
}