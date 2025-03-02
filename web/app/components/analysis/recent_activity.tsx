import { Box, Button, Flex } from "@radix-ui/themes"
import { LineChart } from "../charts/line_chart"
import { emitter } from "~/root";
import { useEffect, useState } from "react";
import type { RecentActivityMessage, ServerMessage } from "~/messages";

export type SessionActivity = {
  session_date: string; // e.g., "2025-02-23"
  session_count: number; // Number of sessions for that day
  total_duration: string; // e.g., "15 hours 30 minutes"
};


export function RecentActivity() {
  const [lineData, setLineData] = useState<any>({
    labels: ["Day 1", "Day 2", "Day 3", "Day 4", "Day 5", "Day 6", "Day 7"],
    datasets: [
      {
        label: "Chat Sessions",
        data: [1, 2, 3, 4, 5, 6, 7]
      }
    ]
  });

  useEffect(() => {
    const handler = (message: ServerMessage) => {
      if (message.type === "session_activity") {
        const recentActivity = message as RecentActivityMessage;
        console.log("Recent activity data:", recentActivity.session_activity);

        // Transform session activity into chart data format
        const labels = recentActivity.session_activity.map((entry) => entry.session_date);
        const data = recentActivity.session_activity.map((entry) => entry.session_count);

        // Update the chart data
        setLineData({
          labels,
          datasets: [
            {
              label: "Chat Sessions",
              data,
              borderColor: "rgb(62, 99, 221)",
              backgroundColor: "rgb(41, 78, 200, 0.4)",
              fill: true,
            }
          ]
        });
      }
    };

    emitter.on("session_activity", handler);
    return () => {
      emitter.off("session_activity", handler);
    };
  }, []);
  const data = {
    labels: ["Day 1", "Day 2", "Day 3", "Day 4", "Day 5", "Day 6", "Day 7"],
    datasets: [
      {
        label: "Chat Sessions",
        data: [1, 2, 3, 4, 5, 6, 7]
      }
    ]
  }
  const aspectRatio = 2;
  return (
    <Box style={{ border: "2px solid var(--indigo-3)", borderRadius: 4 }} p={"2"}>
      <Flex justify={"between"} pb={"2"} px={"1"}>
        <Box>
          Date Picker
        </Box>
        <Flex gap={"2"}>
          <Button size={"1"}>Day</Button>
          <Button size={"1"}>Week</Button>
          <Button size={"1"}>Month</Button>
          <Button size={"1"}>Year</Button>
        </Flex>
      </Flex>
      <Box px={"6"} py={"2"} style={{aspectRatio: aspectRatio ? aspectRatio : "auto", backgroundColor: "var(--indigo-2)"}} className="rounded">
        <LineChart data={lineData} options={{ maintainAspectRatio: true, responsive: true, scales: {y: {beginAtZero: true}} }} />
      </Box>
    </Box>
  )
}