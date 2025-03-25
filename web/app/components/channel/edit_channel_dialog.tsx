import { Box, Button, Dialog, Flex, Text, TextField } from "@radix-ui/themes";
import type React from "react";
import { ChannelListActions, type ChannelAction, type ChannelReducerState } from "./channel_list";

type props = {
  state: ChannelReducerState
  dispatch: React.Dispatch<ChannelAction>
  update: () => void;
}

const EditChannelDialog = ({state, dispatch, update} : props) => {
  return (
    <Dialog.Root
        open={state.id !== null}
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

export default EditChannelDialog;