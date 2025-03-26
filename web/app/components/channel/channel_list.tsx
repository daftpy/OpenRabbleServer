import { Box, Heading, Grid } from "@radix-ui/themes";
import { useEffect, useReducer } from "react";
import { EditChannelDialog } from "./dialog/edit_channel_dialog";
import { ReorderChannelDialog } from "./dialog/reorder_channel_dialog";
import { DeleteChannelDialog } from "./dialog/delete_channel_dialog";
import { ChannelRow } from "./channel_list_row";

export type Channel = {
  id?: number;
  name: string;
  description: string | null;
}

// Shape of the reducer state
export type ChannelReducerState = {
  id: number | null;
  name: string | null;
  description: string | null;
  dialog: ChannelListDialogs | null;
}

// Initial state
const defaultState: ChannelReducerState = {
  id: null,
  name: null,
  description: null,
  dialog: null
}

// Supported reducer actions
export enum ChannelListActions {
  SELECT_CHANNEL = "select_channel",
  CLEAR_SELECTION = "clear_selection",
  SET_NAME = "set_name",
  SET_DESCRIPTION = "set_description"
}

export enum ChannelListDialogs {
  EDIT_CHANNEL = "edit",
  REORDER_CHANNEL = "reorder",
  DELETE_CHANNEL = "delete"
}

export type ChannelAction =
  | { type: ChannelListActions.SELECT_CHANNEL; id: number; dialog: ChannelListDialogs}
  | { type: ChannelListActions.CLEAR_SELECTION; }
  | { type: ChannelListActions.SET_NAME; name: string }
  | { type: ChannelListActions.SET_DESCRIPTION; description: string | null };


function reducer(state: ChannelReducerState, action: ChannelAction) {
  switch (action.type) {
    case ChannelListActions.SELECT_CHANNEL: {
      return { ...state, id: action.id, dialog: action.dialog }
    }
    case ChannelListActions.CLEAR_SELECTION: {
      return { id: null, name: null, description: null, dialog: null }
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

  // Reset selection if channels change
  useEffect(() => {
    dispatch({ type: ChannelListActions.CLEAR_SELECTION });
  }, [channels]);

  // Populate the inputs when a channel is selected
  useEffect(() => {
    if (selectedChannel) {
      dispatch({type: ChannelListActions.SET_NAME, name: selectedChannel.name });
      dispatch({type: ChannelListActions.SET_DESCRIPTION, description: selectedChannel.description });
    }
  }, [selectedChannel]);
  
  return (
    <Box flexGrow={"1"} width={"100%"}>
      <EditChannelDialog state={state} dispatch={dispatch} />
      <ReorderChannelDialog channels={channels} state={state} dispatch={dispatch} />
      <DeleteChannelDialog state={state} dispatch={dispatch} />
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
