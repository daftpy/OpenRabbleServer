import { CircleBackslashIcon, PersonIcon, TimerIcon } from "@radix-ui/react-icons";
import { Box, Button, Container, DropdownMenu, Flex, Heading, ScrollArea, Text } from "@radix-ui/themes";
import { Link } from "react-router";
import { RecentActivity } from "~/components/analysis/recent_activity";
import { MessageList } from "~/components/message/message_list";
import { MessageSearchInput } from "~/components/message/search_input";
import "chart.js/auto"
import type { SessionActivity } from "~/routes/index";
import { useMessageSearch } from "~/hooks/useMessageSearch";

export function UserPage({ username, id, messages, session_activity } : { username: string, id: string, messages: any, session_activity: SessionActivity[] }) {
  const { state, dispatch, nextPage, prevPage } = useMessageSearch({ messages, userId: id, hasMore: false });

  return (
    <Container p={"6"}>
      <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{ color: "var(--slate-12)" }}>
        <Link to="/">Your OnRabble Server</Link>
      </Heading>
      <Flex direction="column" pt={"4"} gap={"6"}>
        <Box>
          <Heading style={{color: "var(--indigo-10)"}}>User Information</Heading>
          <Text>Here you can set a users roles, inspect their activity, and ban them from the chatserver if needed.</Text>
        </Box>
        <Flex align={"center"} gap={"4"}>
          <Flex direction={"column"} gap={"2"} flexGrow={"1"}>
            <Flex align={"center"} gap={"2"}>
              <PersonIcon style={{width: "2em", height: "2em", color: "var(--indigo-12)", padding: "0.2em", border: "1px solid var(--indigo-12)", borderRadius: "15px"}} />
              <Heading style={{color: "var(--indigo-12)"}} size={"8"}>{username}</Heading>
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
            <MessageSearchInput state={state} dispatch={dispatch}  />
          </Box>
          {/* <ScrollArea style={{height: "600px", border: "1px solid var(--indigo-4)", padding: "1em", borderRadius: "4px"}}> */}
            <MessageList messages={state.messages} hidePermaLink={true} />
          {/* </ScrollArea> */}
          <Flex justify={"between"} pt={"4"}>
            <Button 
              disabled={state.page === 0} 
              onClick={() => prevPage()}
            >
              Previous Page
            </Button>
            <Button 
              disabled={!state.hasMore}
              onClick={() => nextPage()}
            >
              Next Page
            </Button>
          </Flex>
        </Box>
        <Box>
          <Heading style={{color: "var(--subheading-color)"}}>Session Activity</Heading>
          <Text>You can review a users session activity, showing you when sessions begin, end, and their duration.</Text>
        </Box>
        <RecentActivity session_activity={session_activity} />
      </Flex>
    </Container>
  )
}