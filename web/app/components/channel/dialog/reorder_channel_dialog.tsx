import { Button, Dialog, Flex, Select, Text } from "@radix-ui/themes";
import { ChannelListActions, ChannelListDialogs, type Channel, type ChannelAction, type ChannelReducerState } from "../channel_list";
import type React from "react";
import { useEffect, useState } from "react";
import { useFetcher } from "react-router";

type props = {
  channels: Channel[];
  // selectedChannel: Channel | undefined;
  state: ChannelReducerState;
  dispatch: React.Dispatch<ChannelAction>;
}

export function ReorderChannelDialog({ channels, state, dispatch } : props) {
  const [selected, setSelected] = useState<string>("");
  const channelFetcher = useFetcher();
  useEffect(() => {
    setSelected(state.id?.toString()?? "");
  }, [state.id])

  
  const reorder = (beforeId: number) => {
    if (state.id == null) return;
    console.log("REORDER", beforeId);

    channelFetcher.submit(
      {
        id: state.id,
        beforeId: beforeId,
        intent: "reorder"
      },
      {
        method: "post",
        action: "/channels",
        encType: "application/x-www-form-urlencoded"
      }
    )
  }
  return (
    <Dialog.Root open={state.dialog == ChannelListDialogs.REORDER_CHANNEL && state.id !== null} onOpenChange={(open) => {
      if (!open) {
        dispatch({type: ChannelListActions.CLEAR_SELECTION });
      }
    }}>
      <Dialog.Content>
        <Dialog.Title align={"center"}>
          Reorder Channel
        </Dialog.Title>
        <Dialog.Description align={"center"}>
          Change the order your channel appears in clients.
        </Dialog.Description>
        <Flex gap={"4"} direction={"column"} pt={"3"} >
          <Flex gap={"2"} align={"baseline"} justify={"center"}>
            <Text>Move before channel</Text>
            <Select.Root defaultValue={state.id?.toString() || ""} onValueChange={(value) => {
              setSelected(value);
            }}>
              <Select.Trigger />
              <Select.Content>
                {channels && channels.map(channel => (
                  <Select.Item key={channel.id} value={channel.id ? channel.id.toString() : channel.name}>{channel.name}</Select.Item>
                ))}
              </Select.Content>
            </Select.Root>
          </Flex>
          <Button onClick={() => {
    if (selected != "") {
      console.log("SELECTED", parseInt(selected));
      reorder(parseInt(selected));
    }
  }}>Reorder</Button>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  )
}