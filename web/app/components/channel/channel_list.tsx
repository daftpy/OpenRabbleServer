import { Flex, Button, Text, Box, Heading, Grid, Dialog, TextField } from "@radix-ui/themes";
import { useEffect, useReducer } from "react";
import { useFetcher } from "react-router";
import ChannelRow from "./channel_list_row";

export type Channel = {
  id?: number;
  name: string;
  description: string | null;
}

// Shape of the reducer state
type ReducerState = {
  id: number | null;
  name: string | null;
  description: string | null;
}

// Initial state
const defaultState: ReducerState = {
  id: null,
  name: null,
  description: null
}

// Supported reducer actions
export enum ChannelListActions {
  SELECT_CHANNEL = "select_channel",
  CLEAR_SELECTION = "clear_selection",
  SET_NAME = "set_name",
  SET_DESCRIPTION = "set_description"
}

export type ChannelAction =
  | { type: ChannelListActions.SELECT_CHANNEL; id: number}
  | { type: ChannelListActions.CLEAR_SELECTION; }
  | { type: ChannelListActions.SET_NAME; name: string }
  | { type: ChannelListActions.SET_DESCRIPTION; description: string | null };


function reducer(state: ReducerState, action: ChannelAction) {
  switch (action.type) {
    case ChannelListActions.SELECT_CHANNEL: {
      return { ...state, id: action.id }
    }
    case ChannelListActions.CLEAR_SELECTION: {
      return { id: null, name: null, description: null }
    }
    case ChannelListActions.SET_NAME: {
      return { ...state, name: action.name }
    }
    case ChannelListActions.SET_DESCRIPTION: {
      return { ...state, description: action.description }
    }
    default: {
      throw Error("Unknown action");
    }
  }
}

export default function ChannelList({ channels }  : { channels: Channel[] }) {
  const [state, dispatch] = useReducer(reducer, defaultState);
  const selectedChannel = channels.find(c => c.id === state.id);
  const channelFetcher = useFetcher();

  // Reset selection if channels change
  useEffect(() => {
    dispatch({ type: ChannelListActions.CLEAR_SELECTION });
  }, [channels]);

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

  // Populate the inputs when a channel is selected
  useEffect(() => {
    if (selectedChannel) {
      dispatch({type: ChannelListActions.SET_NAME, name: selectedChannel.name });
      dispatch({type: ChannelListActions.SET_DESCRIPTION, description: selectedChannel.description });
    }
  }, [selectedChannel]);
  
  return (
    <Box flexGrow={"1"} width={"100%"}>
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
      <Grid columns="1fr 3fr" width={"100%"} gapY={"2"} pt={"2"}>
        <Heading size={"3"} style={{color: "var(--subheading-color)"}}>Channels</Heading>
        <Heading size={"3"} style={{color: "var(--subheading-color)"}}>Description</Heading>
        { channels && channels.map((channel, index) => (
          <ChannelRow key={index} channel={channel} isLast={index == channels.length -1} dispatch={dispatch} />
        ))}
      </Grid>
    </Box>
  )
}
