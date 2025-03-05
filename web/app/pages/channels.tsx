import { useState } from "react";
import { MessagesPerChannel } from "~/components/analysis/messages_per_channel";
import ChannelInput from "~/components/channel/channel_input";
import type { Channel } from "~/components/channel/channel_list";
import ChannelList from "~/components/channel/channel_list";
import "chart.js/auto"
import { Box, Container, Flex, Heading } from "@radix-ui/themes";

export function ChannelPage({ channels }: { channels: Channel[] }) {
  const [channelList, setChannelList] = useState<Channel[]>(channels);

  return (
    <Container className="min-h-full" p={"6"}>
      <Flex direction={"column"} gap={"6"}>
        <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{color: "var(--indigo-9)"}}>Your OnRabble Server</Heading>
        <Flex direction={"column"} gap={"2"}>
          <ChannelInput channelList={channelList} setChannelList={setChannelList} />
          <ChannelList channels={channelList} />
        </Flex>
        <MessagesPerChannel />
      </Flex>
  </Container>
  )
}