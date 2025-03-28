import type { SessionActivity } from "~/routes/index";

export type SessionActivityResult = {
  type: string;
  sender: string;
  payload: {
    session_activity: SessionActivity[];
  };
};

export type ChannelMessageCount = {
  channel: string;
  message_count: number;
}

export type ChannelsActivityResult = {
    type: string;
    sender: string;
    payload: {
        channels: ChannelMessageCount[];
    }
}
