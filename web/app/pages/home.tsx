/*
    The home page of the admin dashboard. Landing here will display the current channels
    of the chatserver as well as connected users.
*/
import { useEffect, useState } from "react"
import { useNavigate } from "react-router";
import { Button, Container, Flex, Heading, Text } from "@radix-ui/themes";

import type { Channel } from "~/components/channel_list";
import ChannelInput from "~/components/channel_input";
import ChannelList from "~/components/channel_list";

/*
  TODO: Added a test navigate button here to move between pages. It works properly
  and does not accidentally trigger a refresh of the auth provider. Perfect!
*/

export function HomePage({ channels }: { channels: Channel[] }) {
    const [channelList, setChannelList] = useState<Channel[]>(channels); // Local state for channels
    useEffect(() => {
        // Use effect
    },[])
    const navigate = useNavigate();
    return (
      <main className="p-4">
        <Container>
          <Flex direction="column" gap="3">
            <div>
              <Button onClick={() => navigate("/about")}>About</Button>
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
