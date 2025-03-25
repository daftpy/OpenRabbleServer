import { CaretSortIcon, Cross2Icon, GearIcon } from "@radix-ui/react-icons";
import { Flex, Button, Text, Box, Heading, Grid, Dialog, TextField } from "@radix-ui/themes";
import React, { useEffect, useState } from "react";
import { useFetcher } from "react-router";

export interface Channel {
  id?: number;
  name: string;
  description: string | null;
}

const ChannelRow = ({ channel, isLast, setId } : { channel : Channel, isLast: boolean, setId: (id: number) => void; }) => {
  const noLine = "none";
  const line = "1px solid var(--indigo-4)";
  
  return (
    <React.Fragment>
      <Box style={{borderBottom: isLast ? noLine : line}} pr={"2"} py={"2"}>
        <Text style={{color: "var(--indigo-12)"}} weight={"bold"} size={"2"}>{ channel.name }</Text>
      </Box>
      <Flex flexGrow={"1"} gap={"4"} overflow={"hidden"} style={{borderBottom: isLast ? noLine : line}}  py={"2"}> 
        <Text truncate size={"2"}>{ channel.description }</Text>
        <Flex gap={"2"} flexGrow={"1"} justify={"end"}>
          <Button color="iris" size={"1"} radius="full" style={{ boxShadow: "var(--shadow-1)", height: "20px", width: "20px"}}><Box><CaretSortIcon style={{width: "14px", height: "14px"}} /></Box></Button>
          <Button color="blue" size={"1"} radius="full" style={{ boxShadow: "var(--shadow-1)", height: "20px", width: "20px"}} onClick={() => channel.id !== undefined && setId(channel.id)}><Box><GearIcon style={{width: "14px", height: "14px"}} /></Box></Button>
          <Button color="red" size={"1"} radius="full" style={{ boxShadow: "var(--shadow-1)", height: "20px", width: "20px"}}><Box><Cross2Icon style={{width: "14px", height: "14px"}} /></Box></Button>
        </Flex>
      </Flex>
    </React.Fragment>
  )
}

export default function ChannelList({ channels }  : { channels: Channel[] }) {
  const [channelId, setChannelId] = useState<number | null>(null);
  const [channelName, setChannelName] = useState<string | null>(null);
  const [channelDescription, setChannelDescription] = useState<string | null>(null);
  const selectedChannel = channels.find(c => c.id === channelId);
  const channelFetcher = useFetcher();

  useEffect(() => {
    // This closes the dialog box because of revalidation after fetcher update()
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
        method: "post", // Use POST here even though it's a PATCH to Go, because React Router forms only support GET/POST
        action: "/channels", // this matches your clientAction route
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