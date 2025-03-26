import { Box, Button, Dialog, Flex, Text, TextField } from "@radix-ui/themes";
import type React from "react";
import { useFetcher } from "react-router";
import { ChannelListActions, type ChannelListAction, type ChannelReducerState } from "~/types/reducers/channelReducer";
import { ChannelListDialogs } from "~/types/components/channel";

type props = {
  state: ChannelReducerState;
  dispatch: React.Dispatch<ChannelListAction>;
}

export function EditChannelDialog({state, dispatch} : props) {
  const channelFetcher = useFetcher();

  // Submit updated channel to the server
  const update = () => {
    if (state.id == null) return;
  
    channelFetcher.submit(
      {
        id: String(state.id),
        name: state.name ?? "",
        description: state.description ?? "",
        intent: "edit"
      },
      {
        method: "post",
        action: "/channels",
        encType: "application/x-www-form-urlencoded",
      }
    );
  };
  
  return (
    <Dialog.Root
        open={state.dialog == ChannelListDialogs.EDIT_CHANNEL && state.id !== null}
        onOpenChange={(open) => {
          if (!open) {
            dispatch({ type: ChannelListActions.CLEAR_SELECTION });
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
                value={state.name ?? ""}
                onChange={(e) => dispatch({type: ChannelListActions.SET_NAME, name: e.target.value })}
              />
            </Box>
            <Box>
              <Box pb="1"><Text>Description</Text></Box>
              <TextField.Root
                value={state.description ?? ""}
                onChange={(e) => dispatch({type: ChannelListActions.SET_DESCRIPTION, description: e.target.value })}
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
  )
}
