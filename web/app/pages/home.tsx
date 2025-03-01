/*
    The home page of the admin dashboard. Landing here will display the current channels
    of the chatserver as well as connected users.
*/
import { useEffect, useState } from "react"
import { useNavigate } from "react-router";
import ChannelInput from "~/components/channel_input";
import ChannelList from "~/components/channel_list";
import { Bar } from "react-chartjs-2";

import { Box, Button, Flex, Heading, Link, Text } from "@radix-ui/themes";
import type { Channel } from "~/components/channel_list";
import UserList from "~/components/user_list";
import ChatMessageList from "~/components/chat_message_list";
import { GearIcon, LockClosedIcon, MagnifyingGlassIcon, PersonIcon } from "@radix-ui/react-icons";
import { emitter } from "~/root";
import type { ChannelMessageCount, ServerMessage } from "~/messages";
import "chart.js/auto"

/*
  TODO: Added a test navigate button here to move between pages. It works properly
  and does not accidentally trigger a refresh of the auth provider. Perfect!
*/
export function HomePage({ channels }: { channels: Channel[] }) {
  const [channelList, setChannelList] = useState<Channel[]>(channels);
  const navigate = useNavigate();

  // State for the bar chart
  const [barData, setBarData] = useState({
    labels: [] as string[],
    datasets: [
      {
        label: "Messages per Channel",
        data: [] as number[],
        backgroundColor: "rgba(54, 162, 235, 0.6)",
        borderColor: "rgba(54, 162, 235, 1)",
        borderWidth: 1,
        borderRadius: 3
      },
    ],
  });
  

  useEffect(() => {
    console.log("loaded");
    const handler = (message: ServerMessage) => {
      if (message.type === "message_count_by_channel") {
        const channels: ChannelMessageCount[] = message.channels;
        console.log("Updating analytics");
        setBarData({
          labels: channels.map((c) => c.channel),
          datasets: [
            {
              label: "Messages per Channel",
              data: channels.map((c) => c.message_count),
              backgroundColor: "rgba(62, 99, 221, 1)",
              borderColor: "rgb(50, 54, 176)",
              borderWidth: 1,
              borderRadius: 3
            },
          ],
        });
      }
    };
  
    emitter.on("message_count_by_channel", handler);
    return () => {
      emitter.off("message_count_by_channel", handler);
    };
  }, []);
  return (
    <main>
      <Flex direction="column" gap={"6"} height={"100%"} maxWidth={"900px"} m={"auto"} flexGrow={"1"} px={"4"} py={"6"}>
        <Flex direction={"column"}>
          <Heading weight={"bold"} className="text-xl pb-1">Your OnRabble Server</Heading>
          <Text>Welcome to your dashboard.</Text>
          <Link href="#">Hide Live Chat</Link>
        </Flex>
        <ChatMessageList />
        <Flex gap={"3"} direction={{initial: "column", sm: "row"}}>
          <Button onClick={() => navigate("/about")}><PersonIcon /> User Management</Button>
          <Button onClick={() => navigate("/about")}><MagnifyingGlassIcon /> Messages</Button>
          <Button onClick={() => navigate("/about")}><GearIcon /> Settings</Button>
          <Button onClick={() => navigate("/about")}><LockClosedIcon /> Keycloak</Button>
        </Flex>
        <Flex direction={"column"} gap={"2"}>
          <div>
            <Heading weight={"bold"} style={{ color: "var(--indigo-9)" }}>Channels</Heading>
            <Text>You can add a new channel or manage your channels below.</Text>
          </div>
          <ChannelInput channelList={channelList} setChannelList={setChannelList} />
        </Flex>
        <Flex gap={"6"} direction={{initial: "column", sm: "row"}} align={{initial: "center", sm: "start"}}>
          <ChannelList channels={channelList} />
          <UserList />
        </Flex>
        <Box>
          <Heading color="indigo">Activity</Heading>
          <Text>Essential server analytics are available. Track basic metrics like how many messages you serve, user statistics, and other activity.</Text>
        </Box>
        <Box style={{ border: "2px solid var(--indigo-3)", borderRadius: 4 }} p={"2"}>
          <Box px={"6"} py={"2"} style={{aspectRatio: "4 / 2", backgroundColor: "var(--indigo-2)"}} className="rounded">
            <Bar data={barData} options={{ maintainAspectRatio: true, responsive: true }} />
          </Box>
        </Box>
      </Flex>
    </main>   
  )
}
