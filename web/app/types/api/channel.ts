import type { Channel } from "../components/channel";

// The required fields to add a channel
export type AddChannelPayload = {
    name: string;
    description: string | null;
}

// The required fields to edit a channel
export type EditChannelPayload = {
  id: number;
  name: string | null;
  description: string | null;
}

// The required fields to reorder a channel
export type ReorderChannelPayload = {
  id: number,
  before_id: number
}

// Response after fetching channels
export type FetchChannelsResponse = {
    channels: Channel[];
};

// Response after editing a channel
export type EditChannelResponse = {
    message: string;
}

// Repsonse after reordering a channel
export type ReorderChannelResponse = {
    message: string;
}

// Response after deleting a channel
export type DeleteChannelResponse = {
    message: string;
}

// Response after adding a channel
export type AddChannelResponse = {
    message: string;
    name: string;
}
