import { MessagesPerChannel } from "~/components/analysis/messages_per_channel";
import ChannelInput from "~/components/channel/channel_input";
import ChannelList from "~/components/channel/channel_list";
import "chart.js/auto"
import { Box, Flex, Heading, Text } from "@radix-ui/themes";
import type { Channel } from "~/types/components/channel";
import type { ChannelMessageCount } from "~/types/api/activity";

export function ChannelPage({ channels, channelActivity }: { channels: Channel[], channelActivity: ChannelMessageCount[] }) {
  return (
    <>
      <Box>
        <Heading color="indigo">Channel Management</Heading>
        <Text>Manage your channels here. You can add and remove channels, change their names, or set role based access.</Text>
      </Box>
      <Flex direction={"column"} gap={"6"} pt={"6"}>
        <Flex direction={"column"} gap={"2"}>
          <ChannelInput />
          <ChannelList channels={channels} />
        </Flex>
          <MessagesPerChannel channelData={channelActivity}  />
      </Flex>
  </>
  )
}
