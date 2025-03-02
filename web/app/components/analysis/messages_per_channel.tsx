import { Box, Heading } from "@radix-ui/themes";
import { useEffect, useState, useMemo } from "react";
import type { ChannelMessageCount, ServerMessage } from "~/messages";
import { emitter } from "~/root";
import { BarChart } from "../charts/bar_chart";

export function MessagesPerChannel() {
  const [channelData, setChannelData] = useState<ChannelMessageCount[]>([]);
  
  // Set up listeners on the emitter
  useEffect(() => {
    const handler = (message: ServerMessage) => {
      if (message.type == "message_count_by_channel") {
        console.log("Got analysis data", message.channels);
        setChannelData(message.channels);
      }
    }

    // Listen for messagers per channel message
    emitter.on("message_count_by_channel", handler);

    // Destroy the listeners
    return () => {
      emitter.off("message_count_by_channel", handler);
    }
    
  }, []);

  // Use memo to prevent rerenders when data does not change
  const labels = useMemo(() => channelData.map((channel) => channel.channel), [channelData]);
  const data = useMemo(() => channelData.map((channel) => channel.message_count), [channelData]);

  // Data for the bar chart, use memo here too
  const barData = useMemo(
    () => ({
      labels,
      datasets: [
        {
          label: "Messages Per Channel",
          data,
          backgroundColor: "#3e63dd",
          borderColor: "rgb(50, 54, 176)",
          borderWidth: 1,
          borderRadius: 3,
        },
      ],
    }),
    [labels, data] // Only recompute when labels or data change
  );

  // default aspect ratio, this could change
  const aspectRatio = 2;

  return (
    <Box>
      <Box pb={"2"}>
        <Heading style={{color: "#415187"}}>Messages By Channel</Heading>
      </Box>
      <Box style={{ border: "2px solid var(--indigo-3)", borderRadius: 4 }} p={"2"}>
        <Box px={"6"} py={"2"} style={{aspectRatio: aspectRatio ? aspectRatio : "auto", backgroundColor: "var(--indigo-2)"}} className="rounded">
          <BarChart data={barData} options={{ maintainAspectRatio: true, responsive: true }} />
        </Box>
      </Box>
    </Box>
  )
}