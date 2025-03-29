import { Box, Heading } from "@radix-ui/themes";
import { useEffect, useState } from "react";
import { BarChart } from "../charts/bar_chart";
import type { ChannelMessageCount } from "~/types/api/activity";

export function MessagesPerChannel({ channelData } : {channelData : ChannelMessageCount[]}) {
  // const [channelData, setChannelData] = useState<ChannelMessageCount[]>([]);
  const [chartData, setChartData] = useState<any>(null);
  
  useEffect(() => {
    // Use memo to prevent rerenders when data does not change
    if (channelData) {
      const labels = channelData.map((channel : ChannelMessageCount) => channel.channel);
      const data = channelData.map((channel : ChannelMessageCount) => channel.message_count);

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