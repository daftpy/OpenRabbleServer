/*
  Represents a channel on the chat server.
  Has a simple name, description, and integer based id.
*/
export type Channel = {
  id: number;
  name: string;
  description: string | null;
}

// Dialogs available for interacting with channels
export enum ChannelListDialogs {
  EDIT_CHANNEL = "edit",
  REORDER_CHANNEL = "reorder",
  DELETE_CHANNEL = "delete"
}
