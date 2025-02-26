/*
    The home page of the admin dashboard. Landing here will display the current channels
    of the chatserver as well as connected users.
*/
import { useState } from "react"
import { useNavigate } from "react-router";
import ChannelInput from "~/components/channel_input";
import ChannelList from "~/components/channel_list";

import { Button, Container, Flex, Heading, Link, Text } from "@radix-ui/themes";
import type { Channel } from "~/components/channel_list";
import UserList from "~/components/user_list";
import ChatMessageList from "~/components/chat_message_list";
import { GearIcon } from "@radix-ui/react-icons";


/*
  TODO: Added a test navigate button here to move between pages. It works properly
  and does not accidentally trigger a refresh of the auth provider. Perfect!
*/

export function HomePage({ channels }: { channels: Channel[] }) {
  const [channelList, setChannelList] = useState<Channel[]>(channels);
  const navigate = useNavigate();

    return (
      <main className="p-4">
        <Container>
          <Flex direction="column" gap="4">
            <Flex direction={"column"}>
              <Heading weight={"bold"} className="text-xl">Your OnRabble Server</Heading>
              <Text m="0">Welcome to your dashboard.</Text>
              <Link href="#">Hide Live Chat</Link>
            </Flex>
            <ChatMessageList />
            <Flex>
              <Button onClick={() => navigate("/about")}><GearIcon /> Settings</Button>
            </Flex>
            <div>
              <Heading weight={"bold"} style={{ color: "var(--indigo-9)" }}>Channels</Heading>
              <Text m="0">You can add a new channel or manage your channels below.</Text>
            </div>
            <ChannelInput channelList={channelList} setChannelList={setChannelList} />
            <Flex gap={"6"}>
              <ChannelList channels={channelList} />
              <UserList />
            </Flex>
          </Flex>
        </Container>
      </main>   
    )
}
