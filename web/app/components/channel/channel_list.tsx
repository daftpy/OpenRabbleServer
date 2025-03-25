import { CaretSortIcon, Cross2Icon, GearIcon } from "@radix-ui/react-icons";
import { Flex, Button, Text, Box, Heading, Grid, Dialog, TextField } from "@radix-ui/themes";
import React, { useEffect, useState } from "react";
import { useFetcher } from "react-router";
import ChannelRow from "./channel_list_row";

export interface Channel {
  id?: number;
  name: string;
  description: string | null;
}

export default function ChannelList({ channels }  : { channels: Channel[] }) {
  const [channelId, setChannelId] = useState<number | null>(null);
  const [channelName, setChannelName] = useState<string | null>(null);
  const [channelDescription, setChannelDescription] = useState<string | null>(null);
  const selectedChannel = channels.find(c => c.id === channelId);
  const channelFetcher = useFetcher();

  useEffect(() => {
    // Reset selection if channels change
    setChannelId(null);
    setChannelName(null);
    setChannelDescription(null);
  }, [channels]);

  const update = () => {
    if (channelId == null) return;
  
    channelFetcher.submit(
      {
        id: String(channelId),
        name: channelName ?? "",
        description: channelDescription ?? "",
        intent: "edit"
      },
      {
        method: "post",
        action: "/channels",
        encType: "application/x-www-form-urlencoded",
      }
    );
  };

  useEffect(() => {
    if (selectedChannel) {
      setChannelName(selectedChannel.name);
      setChannelDescription(selectedChannel.description ?? "");
    }
  }, [selectedChannel]);
  
  return (
    <Box flexGrow={"1"} width={"100%"}>
      <Dialog.Root
        open={channelId !== null}
        onOpenChange={(open) => {
          if (!open) {
            setChannelId(null);
            setChannelName(null);
            setChannelDescription(null);
          }
        }}
      >
        <Dialog.Content>
          <Dialog.Title align="center">Update Channel</Dialog.Title>
          <Dialog.Description align={"center"}>
            Update the channel name and description below.
          </Dialog.Description>
          <Flex direction={"column"} gap={"6"} pt={"4"}>
            <Box>
              <Box pb="1"><Text>Name</Text></Box>
              <TextField.Root
                value={channelName ?? ""}
                onChange={(e) => setChannelName(e.target.value)}
              />
            </Box>
            <Box>
              <Box pb="1"><Text>Description</Text></Box>
              <TextField.Root
                value={channelDescription ?? ""}
                onChange={(e) => setChannelDescription(e.target.value)}
              />
            </Box>
            <Flex>
              <Button style={{ flexGrow: "1" }} onClick={() => update()}>
                Update
              </Button>
            </Flex>
          </Flex>
        </Dialog.Content>
      </Dialog.Root>
      <Grid columns="1fr 3fr" width={"100%"} gapY={"2"} pt={"2"}>
        <Heading size={"3"} style={{color: "var(--subheading-color)"}}>Channels</Heading>
        <Heading size={"3"} style={{color: "var(--subheading-color)"}}>Description</Heading>
        { channels && channels.map((channel, index) => (
          <ChannelRow key={index} channel={channel} isLast={index == channels.length -1} setId={setChannelId} />
        ))}
      </Grid>
    </Box>
  )
}
