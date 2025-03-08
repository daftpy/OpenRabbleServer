import { CircleBackslashIcon, PersonIcon, TimerIcon } from "@radix-ui/react-icons";
import { Box, Button, Container, DropdownMenu, Flex, Heading, ScrollArea, Text } from "@radix-ui/themes";
import { useState } from "react";
import { Link } from "react-router";
import { MessageList } from "~/components/message/message_list";
import { MessageSearchInput } from "~/components/message/search_input";

export function UserPage({ username, id, messages } : { username: string, id: string, messages: any }) {
  const [filteredMessages, setFilteredMessages] = useState<any>(null);

  const handleMessageUpdate = (messages: any) => {
    setFilteredMessages(messages);
  }

  return (
    <Container p={"6"}>
      <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{ color: "var(--slate-12)" }}>
        <Link to="/">Your OnRabble Server</Link>
      </Heading>
      <Flex direction="column" pt={"4"} gap={"4"}>
        <Box>
        <Heading style={{color: "var(--indigo-10)"}}>User Information</Heading>
        <Text>Here you can set a users roles, inspect their activity, and ban them from the chatserver if needed.</Text>
        </Box>
        <Flex align={"center"} gap={"4"}>
          <Flex direction={"column"} gap={"2"} flexGrow={"1"}>
            <Flex align={"center"} gap={"2"}>
              <PersonIcon style={{width: "1.5em", height: "1.5em", color: "var(--slate-11)", padding: "0.1em", border: "1px solid var(--slate-11)", borderRadius: "15px"}} />
              <Heading style={{color: "var(--link-color)"}}>{username}</Heading>
            </Flex>
            <Flex gap={"2"}>
              <Text weight={"bold"} size={"1"}>ID: </Text>
              <Text size={"1"} style={{color: "var(--muted-text-color)"}}>{ id }</Text>
            </Flex>
          </Flex>
          <DropdownMenu.Root>
            <DropdownMenu.Trigger>
              <Button color="red">Ban</Button>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content>
                <DropdownMenu.Item>
                  <CircleBackslashIcon /> Permanent
                </DropdownMenu.Item>
                <DropdownMenu.Item>
                  <TimerIcon /> Temporary
                </DropdownMenu.Item>
            </DropdownMenu.Content>
          </DropdownMenu.Root>
        </Flex>
        <Box pt={"2"}>
          <Box pb={"4"}>
            <Heading style={{color: "var(--subheading-color)"}}>Message History</Heading>
            <Text>You can search through a users chat history and filter by channel or keyword.</Text>
            <MessageSearchInput userId={id} onMessagesUpdate={handleMessageUpdate} />
          </Box>
          <ScrollArea style={{maxHeight: "300px", border: "1px solid var(--indigo-4)", padding: "1em", borderRadius: "4px"}}>
            <MessageList messages={filteredMessages ? filteredMessages : messages} hidePermaLink={true} />
          </ScrollArea>
        </Box>
        <Box>
          <Heading style={{color: "var(--subheading-color)"}}>Session Activity</Heading>
          <Text>You can review a users session activity, showing you when sessions begin, end, and their duration.</Text>
        </Box>
      </Flex>
    </Container>
  )
}