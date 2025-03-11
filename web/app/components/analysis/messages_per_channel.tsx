import { Box, Heading } from "@radix-ui/themes";
import { useEffect, useState, useMemo } from "react";
import type { ChannelMessageCount, ServerMessage } from "~/messages";
import { emitter } from "~/root";
import { BarChart } from "../charts/bar_chart";

export function MessagesPerChannel({ channelData } : {channelData : ChannelMessageCount[]}) {
  // const [channelData, setChannelData] = useState<ChannelMessageCount[]>([]);
  const [chartData, setChartData] = useState<any>(null);
  
  useEffect(() => {
    // Use memo to prevent rerenders when data does not change
    if (channelData) {
      console.log("CHANNEL DATA", channelData);
      const labels = channelData.map((channel : any) => channel.channel);
      const data = channelData.map((channel : any) => channel.message_count);

      // Data for the bar chart, use memo here too
      setChartData({
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
      });
      }
  }, [channelData])

  // default aspect ratio, this could change
  const aspectRatio = 2;

  return (
    <Box>
      <Box pb={"2"}>
        <Heading style={{color: "var(--subheading-color)"}}>Messages By Channel</Heading>
      </Box>
      <Box style={{ border: "2px solid var(--indigo-3)", borderRadius: 4 }} p={"2"}>
        <Box px={"6"} py={"2"} style={{aspectRatio: aspectRatio ? aspectRatio : "auto", backgroundColor: "var(--indigo-2)"}} className="rounded">
          {chartData && (<BarChart data={chartData} options={{ maintainAspectRatio: true, responsive: true }} />)}
        </Box>
      </Box>
    </Box>
  )
}