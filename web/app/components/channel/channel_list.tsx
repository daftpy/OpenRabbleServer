import { CaretSortIcon, Cross2Icon, GearIcon } from "@radix-ui/react-icons";
import { Flex, Button, Text, Box, Heading, Grid } from "@radix-ui/themes";
import React from "react";

export interface Channel {
  name: string;
  description: string | null;
}

const ChannelRow = ({ channel, isLast } : { channel : Channel, isLast: boolean }) => {
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
          <Button color="blue" size={"1"} radius="full" style={{ boxShadow: "var(--shadow-1)", height: "20px", width: "20px"}}><Box><GearIcon style={{width: "14px", height: "14px"}} /></Box></Button>
          <Button color="red" size={"1"} radius="full" style={{ boxShadow: "var(--shadow-1)", height: "20px", width: "20px"}}><Box><Cross2Icon style={{width: "14px", height: "14px"}} /></Box></Button>
        </Flex>
      </Flex>
    </React.Fragment>
  )
}

export default function ChannelList({ channels }  : { channels: Channel[] }) {
  console.log("LOOK HERE THRICE", channels);
  return (
    <Box flexGrow={"1"} width={"100%"}>
      <Grid columns="1fr 3fr" width={"100%"} gapY={"2"} pt={"2"}>
        <Heading size={"3"} style={{color: "var(--subheading-color)"}}>Channels</Heading>
        <Heading size={"3"} style={{color: "var(--subheading-color)"}}>Description</Heading>
        { channels && channels.map((channel, index) => (
          <ChannelRow key={index} channel={channel} isLast={index == channels.length -1} />
        ))}
      </Grid>
    </Box>
  )
}