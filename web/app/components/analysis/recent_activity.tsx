import { Box, Button, Dialog, Flex, Heading, useThemeContext } from "@radix-ui/themes"
import { LineChart } from "../charts/line_chart"
import { emitter } from "~/root";
import { useEffect, useState } from "react";
import type { RecentActivityMessage, ServerMessage } from "~/messages";
import { DayPicker } from "react-day-picker";
import "react-day-picker/style.css";
import { ChatBubbleIcon, LapTimerIcon } from "@radix-ui/react-icons";

export type SessionActivity = {
  session_date: string; // e.g., "2025-02-23"
  session_count: number; // Number of sessions for that day
  total_duration: string; // e.g., "15 hours 30 minutes"
};


export function RecentActivity() {
  const [selected, setSelected] = useState<Date>();
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
              backgroundColor: "rgb(98, 132, 244)",
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

  // Colors for buttons
  const {appearance} = useThemeContext();
  let buttonStyle = {backgroundColor: "var(--secondary-button-color)", border: ""};
  if (appearance === "dark") {
    buttonStyle.border = "1px solid var(--indigo-4)";
  }
  return (
    <Box>
      <Flex justify={"between"} pb={"2"}>
        <Heading style={{color: "var(--subheading-color)"}}>Recent Activity</Heading>
        <Flex  gap={"2"}>
          <Button size="2"><LapTimerIcon /> Sessions</Button>
          <Button size="2" color="gray" style={buttonStyle}><ChatBubbleIcon /> Messages</Button>
        </Flex>
      </Flex>
      <Box style={{ border: "2px solid var(--indigo-3)", borderRadius: 4 }} p={"2"}>
        <Flex justify={"between"} pb={"2"} px={"1"}>
          <Dialog.Root>
            <Dialog.Trigger>
              <Button size={"1"} color={"gray"} style={buttonStyle}>Select Range</Button>
            </Dialog.Trigger>

            <Dialog.Content maxWidth={"350px"}>
              <Box className="size-fit m-auto">
                <DayPicker
                  mode="single"
                  selected={selected}
                  onSelect={setSelected}
                  footer={
                    selected ? `Selected: ${selected.toLocaleDateString()}` : "Pick a day."
                  }
                />
                <Flex justify={"end"}>
                  <Button>Select</Button>
                </Flex>
              </Box>
            </Dialog.Content>
          </Dialog.Root>
          <Flex gap={"2"}>
            <Button size={"1"} color={"gray"} style={buttonStyle}>Day</Button>
            <Button size={"1"} color={"indigo"} >Week</Button>
            <Button size={"1"} color={"gray"} style={buttonStyle}>Month</Button>
            <Button size={"1"} color={"gray"} style={buttonStyle}>Year</Button>
          </Flex>
        </Flex>
        <Box px={"4"} py={"2"} style={{aspectRatio: aspectRatio ? aspectRatio : "auto", backgroundColor: "var(--indigo-2)"}} className="rounded">
          <LineChart data={lineData} options={{ maintainAspectRatio: true, responsive: true, scales: {y: {beginAtZero: true}} }} />
        </Box>
      </Box>
    </Box>
  )
}