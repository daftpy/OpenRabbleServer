import type { ChannelListDialogs } from "../components/channel";

// Shape of the reducer state
export type ChannelReducerState = {
  id: number | null;
  name: string | null;
  description: string | null;
  dialog: ChannelListDialogs | null;
}

// Supported reducer action types
export enum ChannelListActions {
  SELECT_CHANNEL = "select_channel",
  CLEAR_SELECTION = "clear_selection",
  SET_NAME = "set_name",
  SET_DESCRIPTION = "set_description"
}

export type ChannelListAction =
  | { type: ChannelListActions.SELECT_CHANNEL; id: number; dialog: ChannelListDialogs}
  | { type: ChannelListActions.CLEAR_SELECTION; }
  | { type: ChannelListActions.SET_NAME; name: string }
  | { type: ChannelListActions.SET_DESCRIPTION; description: string | null };