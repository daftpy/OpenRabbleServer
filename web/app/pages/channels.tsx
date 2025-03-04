import { useState } from "react";
import { MessagesPerChannel } from "~/components/analysis/messages_per_channel";
import ChannelInput from "~/components/channel/channel_input";
import type { Channel } from "~/components/channel/channel_list";
import ChannelList from "~/components/channel/channel_list";
import "chart.js/auto"
import { Box, Container, Flex } from "@radix-ui/themes";

export function ChannelPage({ channels }: { channels: Channel[] }) {
  const [channelList, setChannelList] = useState<Channel[]>(channels);
  return (
    <Container className="min-h-full bg-red-100">
      <Flex direction={"column"} gap={"6"}>

          <Box>
            <ChannelInput channelList={channelList} setChannelList={setChannelList} />
            <ChannelList channels={channelList} />
          </Box>
          <MessagesPerChannel />

      </Flex>
  </Container>
  )
}