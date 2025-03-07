import { useState } from "react";
import { MessagesPerChannel } from "~/components/analysis/messages_per_channel";
import ChannelInput from "~/components/channel/channel_input";
import type { Channel } from "~/components/channel/channel_list";
import ChannelList from "~/components/channel/channel_list";
import "chart.js/auto"
import { Box, Container, Flex, Heading, Text } from "@radix-ui/themes";
import { Link } from "react-router";

export function ChannelPage({ channels }: { channels: Channel[] }) {
  const [channelList, setChannelList] = useState<Channel[]>(channels);

  return (
    <Container className="min-h-full" p={"6"}>
      <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{ color: "var(--slate-12)" }}>
        <Link to="/">Your OnRabble Server</Link>
      </Heading>
      <Heading style={{color: "var(--indigo-10)"}}>Channel Management</Heading>
      <Text>Manage your channels here. You can add and remove channels, change their names, or set role based access.</Text>
      <Flex direction={"column"} gap={"6"} pt={"6"}>
        <Flex direction={"column"} gap={"2"}>
          <ChannelInput channelList={channelList} setChannelList={setChannelList} />
          <ChannelList channels={channelList} />
        </Flex>
          <MessagesPerChannel />
      </Flex>
  </Container>
  )
}