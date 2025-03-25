import { Box, Button, Flex, Text } from "@radix-ui/themes";
import React from "react";
import type { Channel } from "./channel_list";
import { CaretSortIcon, Cross2Icon, GearIcon } from "@radix-ui/react-icons";


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

export default ChannelRow;
