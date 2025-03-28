import { CircleBackslashIcon, CrossCircledIcon, PersonIcon, TimerIcon } from "@radix-ui/react-icons";
import { Box, Button, Container, DropdownMenu, Flex, Heading, Text } from "@radix-ui/themes";
import { Link } from "react-router";
import { RecentActivity } from "~/components/analysis/recent_activity";
import { MessageList } from "~/components/message/message_list";
import { MessageSearchInput } from "~/components/message/search_input";
import "chart.js/auto"
import type { SessionActivity } from "~/routes/index";
import { useMessageSearch } from "~/hooks/useMessageSearch";
import { useEffect, useMemo, useState } from "react";
import { TempBanDialog } from "~/components/user/temp_ban_dialog";
import { PermBanDialog } from "~/components/user/perm_ban_dialog";

// Ban dialog types that can appear
export enum BanDialog {
  PermaBan,
  TempBan
}

export function UserPage({ username, id, isBanned, messages, hasMore, session_activity } : { username: string, id: string, isBanned: boolean, messages: any, hasMore: boolean, session_activity: SessionActivity[] }) {
  const { state, messageFetcher, dispatch, nextPage, prevPage } = useMessageSearch({ messages, userId: id, hasMore: hasMore });
  const [dialog, setDialog] = useState<BanDialog | null>(null); // Whether temp or perm ban forms appear

  // Reset the dialog when ban status changes
  useEffect(() => {
    setDialog(null);
  }, [isBanned])

  const banButton = useMemo(() => {
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
            setDialog(BanDialog.PermaBan);
          }}>
            <CircleBackslashIcon /> Permanent
          </DropdownMenu.Item>
          <DropdownMenu.Item onClick={() => {
            setDialog(BanDialog.TempBan);
          }}>
            <TimerIcon /> Temporary
          </DropdownMenu.Item>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    )
  }, [isBanned])

  return (
    <>
      <Flex direction="column" gap={"6"}>
        <Box>
          <Heading color="indigo" size={"7"}>User Information</Heading>
          <Text>Here you can set a users roles, inspect their activity, and ban them from the chatserver if needed.</Text>
        </Box>
        <TempBanDialog id={id} username={username} dialog={dialog?? undefined} setDialog={setDialog} />
        <PermBanDialog id={id} username={username} dialog={dialog?? undefined} setDialog={setDialog}/>
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
          { banButton }
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
              disabled={!state.has_more}
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
    </>
  )
}
