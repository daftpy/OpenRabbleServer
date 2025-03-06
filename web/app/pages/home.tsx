/*
    The home page of the admin dashboard. Landing here will display the current channels
    of the chatserver as well as connected users.
*/
import { useState } from "react"
import { Link, NavLink, useNavigate } from "react-router";
import ChannelInput from "~/components/channel/channel_input";
import ChannelList from "~/components/channel/channel_list";

import { Box, Button, Flex, Heading, Text } from "@radix-ui/themes";
import type { Channel } from "~/components/channel/channel_list";
import UserList from "~/components/user_list";
import { GearIcon, LockClosedIcon, MagnifyingGlassIcon, PersonIcon } from "@radix-ui/react-icons";
import "chart.js/auto"
import { MessagesPerChannel } from "~/components/analysis/messages_per_channel";
import { RecentActivity } from "~/components/analysis/recent_activity";
import { LiveView } from "~/components/message/live_view";

/*
  TODO: Added a test navigate button here to move between pages. It works properly
  and does not accidentally trigger a refresh of the auth provider. Perfect!
*/
export function HomePage({ channels }: { channels: Channel[] }) {
  const [channelList, setChannelList] = useState<Channel[]>(channels);
  const navigate = useNavigate();
  
  return (
    <main style={{color: "var(--primary-text-color)"}}>
      <Flex direction="column" gap={"6"} height={"100%"} maxWidth={"900px"} m={"auto"} flexGrow={"1"} px={"4"} py={"6"}>

        <Flex direction={"column"}>
          <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{color: "var(--indigo-9)"}}>Your OnRabble Server</Heading>
          <Text>Welcome to your dashboard.</Text>
          <NavLink to="#" style={{color: "var(--link-color)"}}>Hide Live Chat</NavLink>
        </Flex>

        <LiveView />

        <Flex gap={"3"} direction={{initial: "column", sm: "row"}}>
          <Button onClick={() => navigate("/users")}><PersonIcon /> User Management</Button>
          <Button onClick={() => navigate("/messages")}><MagnifyingGlassIcon /> Messages</Button>
          <Button onClick={() => navigate("/about")}><GearIcon /> Settings</Button>
          <Button onClick={() => navigate("/about")}><LockClosedIcon /> Keycloak</Button>
        </Flex>

        <Flex direction={"column"} gap={"2"}>
          <div>
            <Heading size={"7"} weight={"bold"} style={{ color: "var(--indigo-10)" }}>
              <Link to="/channels">Channels</Link>
            </Heading>
            <Text>You can add a new channel or manage your channels below.</Text>
          </div>
          <ChannelInput channelList={channelList} setChannelList={setChannelList} />
          <Flex gap={"6"} direction={{initial: "column", sm: "row"}} align={{initial: "center", sm: "start"}}>
            <ChannelList channels={channelList} />
            <UserList />
          </Flex>
        </Flex>
        
        <Box>
          <Heading size={"7"} style={{ color: "var(--indigo-10)" }}>Analytics</Heading>
          <Text>Essential server analytics are available. Track basic metrics like how many messages you serve, user statistics, and other activity.</Text>
        </Box>

        <RecentActivity />

      </Flex>
    </main>   
  )
}
