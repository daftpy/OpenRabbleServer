export type Channel = {
  id: number;
  name: string;
  description: string | null;
}

export enum ChannelListDialogs {
  EDIT_CHANNEL = "edit",
  REORDER_CHANNEL = "reorder",
  DELETE_CHANNEL = "delete"
}
