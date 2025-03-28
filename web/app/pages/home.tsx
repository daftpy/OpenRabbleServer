/*
    The home page of the admin dashboard. Landing here will display the current channels
    of the chatserver as well as connected users.
*/
import { Link, NavLink } from "react-router";
import ChannelInput from "~/components/channel/channel_input";
import ChannelList from "~/components/channel/channel_list";

import { Box, Flex, Heading, Text } from "@radix-ui/themes";
import "chart.js/auto"
import { RecentActivity } from "~/components/analysis/recent_activity";
import { LiveView } from "~/components/message/live_view";
import type { SessionActivity } from "../routes/index";
import UserList from "~/components/user/user_list";
import type { Channel } from "~/types/components/channel";

export function HomePage({ channels, session_activity }: { channels: Channel[], session_activity: SessionActivity[] }) {
  return (
    <main style={{color: "var(--primary-text-color)"}}>
      <Flex direction="column" gap={"6"} height={"100%"} m={"auto"} flexGrow={"1"}>
        <Flex direction={"column"}>
          <Text>Welcome to your dashboard.</Text>
          <NavLink to="#" style={{color: "var(--link-color)"}}>
            <Text size={"2"}>
              Hide Live Chat
            </Text>
          </NavLink>
          
          <Box pt={"3"}>
            <LiveView />
          </Box>
        </Flex>

        <Flex direction={"column"} gap={"2"}>
          <div>
            <Heading size={"7"} weight={"bold"} color="indigo">
              <Link to="/channels">Channels</Link>
            </Heading>
            <Text>You can add a new channel or manage your channels below.</Text>
          </div>
          <ChannelInput />
          <Flex gap={"6"} direction={{initial: "column", sm: "row"}} align={{initial: "center", sm: "start"}}>
            <ChannelList channels={channels} />
            <UserList />
          </Flex>
        </Flex>
        
        <Box>
          <Heading size={"7"} color="indigo">Analytics</Heading>
          <Text>Essential server analytics are available. Track basic metrics like how many messages you serve, user statistics, and other activity.</Text>
        </Box>

        <RecentActivity session_activity={session_activity} />
      </Flex>
    </main>   
  )
}
