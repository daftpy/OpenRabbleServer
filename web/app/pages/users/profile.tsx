import { CircleBackslashIcon, CrossCircledIcon, PersonIcon, TimerIcon } from "@radix-ui/react-icons";
import { Box, Button, Container, Dialog, DropdownMenu, Flex, Heading, ScrollArea, Text, TextField } from "@radix-ui/themes";
import { Link, useFetcher } from "react-router";
import { RecentActivity } from "~/components/analysis/recent_activity";
import { MessageList } from "~/components/message/message_list";
import { MessageSearchInput } from "~/components/message/search_input";
import "chart.js/auto"
import type { SessionActivity } from "~/routes/index";
import { useMessageSearch } from "~/hooks/useMessageSearch";
import { useState } from "react";

export function UserPage({ username, id, isBanned, messages, hasMore, session_activity } : { username: string, id: string, isBanned: boolean, messages: any, hasMore: boolean, session_activity: SessionActivity[] }) {
  console.log("HAS MORE?", hasMore);
  type banState = "temp" | "perm";
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [banState, setBanState] = useState<banState>("temp");
  const [banDuration, setBanDuration] = useState<string>("");
  const [banReason, setBanReason] = useState<string>("");
  const { state, messageFetcher, dispatch, nextPage, prevPage } = useMessageSearch({ messages, userId: id, hasMore: hasMore });
  const fetcher = useFetcher();
  

  const handleBanUser = (duration? : string) => {
    const formData = new FormData();
    formData.append("banishedId", id);
    if (banReason.trim().length != 0) {
      formData.append("reason", banReason);
    }

    if (duration !== undefined) {
      formData.append("duration", duration.toString());
    }

    fetcher.submit(formData, { method: "POST" });
  };

  const banButton = () => {
    if (isBanned) {
      return (
        <Button color="red">Unban</Button>
      )
    }
    return (
      <DropdownMenu.Root>
        <DropdownMenu.Trigger>
          <Button color="red">Ban</Button>
        </DropdownMenu.Trigger>
        <DropdownMenu.Content>
            <DropdownMenu.Item onClick={() => {
              setIsOpen(true);
              setBanState("perm");
            }}>
              <CircleBackslashIcon /> Permanent
            </DropdownMenu.Item>
            <DropdownMenu.Item onClick={() => {
              setIsOpen(true);
              setBanState("temp");
            }}>
              <TimerIcon /> Temporary
            </DropdownMenu.Item>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    )
  }

  const tempForm = () => {
    return (
      <Dialog.Content>
        <Flex direction={"column"} py={"6"} gap={"6"}>
          <Flex direction={"column"} gap={"2"}>
            <Text align={"center"} weight={"bold"}>Temporary Ban</Text>
            <Text size={"2"}>Reason</Text>
            <TextField.Root placeholder="Reason" />
          </Flex>
          <Flex direction={"column"} gap={"2"}>
            <Text size={"2"}>Time (in hours):</Text>
            <TextField.Root placeholder={"e.g. 3"} onChange={(e) => setBanDuration(e.target.value)} />
          </Flex>
          <Button color="red" onClick={() => {handleBanUser(banDuration); setIsOpen(false);}}>Ban</Button>
        </Flex>
      </Dialog.Content>
    )
  }

  const permForm = () => {
    return (
      <Dialog.Content>
        <Flex direction={"column"} py={"6"} gap={"6"}>
          <Flex direction={"column"} gap={"2"}>
            <Text align={"center"} weight={"bold"}>Permanent Ban</Text>
            <Text size={"2"}>Reason</Text>
            <TextField.Root placeholder="Reason" />
          </Flex>
          <Button color="red" onClick={() => {handleBanUser(); setIsOpen(false);}}>Ban</Button>
        </Flex>
      </Dialog.Content>
    )
  }

  return (
    <Container p={"6"}>
      <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{ color: "var(--slate-12)" }}>
        <Link to="/">OnRabble Server</Link>
      </Heading>
      <Flex direction="column" pt={"4"} gap={"6"}>
        <Box>
          <Heading style={{color: "var(--indigo-10)"}} size={"7"}>User Information</Heading>
          <Text>Here you can set a users roles, inspect their activity, and ban them from the chatserver if needed.</Text>
        </Box>
        <Dialog.Root open={isOpen} onOpenChange={setIsOpen}>
          {banState == "temp" ? tempForm() : permForm()}
        </Dialog.Root>
        <Flex align={"center"} gap={"4"}>
          <Flex direction={"column"} gap={"3"} flexGrow={"1"}>
            <Flex align={"center"} gap={"2"}>
              {isBanned ? 
                <CrossCircledIcon style={{width: "2em", height: "2em", color: "var(--red-9)", padding: "0.2em"}} /> 
              : <PersonIcon style={{width: "2em", height: "2em", color: "var(--indigo-12)", padding: "0.2em", border: "1px solid var(--indigo-12)", borderRadius: "15px"}} />
              }
              <Heading style={{color: isBanned ? "var(--red-9)" : "var(--indigo-12)"}} size={"8"}>{username}</Heading>
            </Flex>
            <Flex gap={"2"}>
              <Text weight={"bold"} size={"2"}>ID: </Text>
              <Text size={"2"} style={{color: "var(--muted-text-color)"}}>{ id }</Text>
            </Flex>
          </Flex>
          { banButton() }
        </Flex>
        <Box pt={"2"}>
          <Box pb={"4"}>
            <Heading style={{color: "var(--subheading-color)"}} size={"6"}>Message History</Heading>
            <Text>You can search through a users chat history and filter by channel or keyword.</Text>
            <MessageSearchInput state={state} dispatch={dispatch}  />
          </Box>
          {/* <ScrollArea style={{height: "600px", border: "1px solid var(--indigo-4)", padding: "1em", borderRadius: "4px"}}> */}
            <MessageList messages={messageFetcher.data?.messages ? messageFetcher.data.messages : messages} hidePermaLink={true} />
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