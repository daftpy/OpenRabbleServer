/*
    The home page of the admin dashboard. Landing here will display the current channels
    of the chatserver as well as connected users.
*/
import ChannelInput from "~/components/channel_input";
import ChannelList from "~/components/channel_list";

import { Container, Flex, Heading, Text } from "@radix-ui/themes";
import { useEffect, useState } from "react"
import type { Channel } from "~/components/channel_list";

export function Home({ channels }: { channels: Channel[] }) {
    const [channelList, setChannelList] = useState<Channel[]>(channels); // Local state for channels
    useEffect(() => {
        // Use effect
    },[])
    return (
      <main className="p-4">
        <Container>
          <Flex direction="column" gap="3">
            <div>
              <Heading weight={"bold"} className="text-xl">Your OnRabble Server</Heading>
            </div>
            <Text>Welcome to your dashboard.</Text>
            <div>
              <Heading weight={"bold"} style={{ color: "var(--indigo-9)" }}>Channels</Heading>
              <Text m="0">You can add a new channel or manage your channels below.</Text>
            </div>
            <ChannelInput channelList={channelList} setChannelList={setChannelList} />
            <ChannelList channels={channelList} />
          </Flex>
        </Container>
      </main>   
    )
}
