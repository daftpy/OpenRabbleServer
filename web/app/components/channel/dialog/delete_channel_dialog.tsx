import { Button, Dialog, Flex, Switch, Text } from "@radix-ui/themes";
import { useFetcher } from "react-router";
import { useState } from "react";
import { ChannelListActions, type ChannelListAction, type ChannelReducerState } from "~/types/reducers/channelReducer";
import { ChannelListDialogs } from "~/types/components/channel";

type props = {
  state: ChannelReducerState;
  dispatch: React.Dispatch<ChannelListAction>;
}

export function DeleteChannelDialog({ state, dispatch } : props) {
  const [purge, setPurge] = useState<boolean>(false); // Tracks if we're deleting channels messages
  const channelFetcher = useFetcher();

  const deleteChannel = () => {
    channelFetcher.submit(
      {
        id: state.id,
        purge: purge ? "1" : "0",
        intent: "delete"
      },
      {
        method: "POST",
        action: "/channels",
        encType: "application/x-www-form-urlencoded"
      }
    )
  }

  return (
    <Dialog.Root open={state.dialog == ChannelListDialogs.DELETE_CHANNEL && state.id !== null}
      onOpenChange={() => dispatch({type: ChannelListActions.CLEAR_SELECTION})}>
      <Dialog.Content align="center">
        <Dialog.Title align={"center"}>Delete Channel {state.name}</Dialog.Title>
        <Dialog.Description align={"center"}>Remove the selected channel from the databasee.</Dialog.Description>
        <Flex direction={"column"} gap={"3"} align={"center"} pt={"4"}>
          <Flex gap={"3"} align={"center"}>
            <Text>Delete all associated channels.</Text>
            <Switch color="red" onCheckedChange={setPurge} />
          </Flex>
          <Button color="red" onClick={() => deleteChannel()}>Delete</Button>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  )
}
