import type { SessionActivity } from "./routes/index";
import type { ChatMessageType } from "./components/message/live_view";
import type { Channel } from "./types/components/channel";
import type { ChannelMessageCount } from "./types/api/activity";

// Base message interface
export interface Message {
    type: string;
  }

export interface ConnectedUsersMessage extends Message {
  type: "connected_users";
  users: string[];
}

export interface UserStatusMessage extends Message {
  type: "user_status";
  payload: {
    username: string;
    status: boolean;
  }
}

export interface ChatMessage extends Message {
  type: "chat_message";
  payload: {
    message: string;
    username: string;
    channel: string;
    authored_at: string;
    id: number;
  }
}

export interface ActiveChannelsMessage extends Message {
  type: "active_channels";
  channels: Channel[]; 
}

export interface BulkChatMessages extends Message {
  type: "bulk_chat_messages";
  // Use the meessage type for react state
  payload: {
    messages: ChatMessageType[];
  }
}

export interface MessageCountByChannelMessage extends Message {
  type: "message_count_by_channel";
  payload: {
    channels: ChannelMessageCount[]
  }
}

export interface RecentActivityMessage extends Message {
  type: "session_activity";
  payload: {
    session_activity: SessionActivity[];
  }
}
  
  // Create a union type for all server messages
export type ServerMessage =
  | ConnectedUsersMessage
  | UserStatusMessage
  | ChatMessage
  | BulkChatMessages
  | MessageCountByChannelMessage
  | ActiveChannelsMessage
  | RecentActivityMessage;

  export type EmitterEvents = {
    connected_users: ServerMessage;
    user_status: ServerMessage;
    active_channels: ServerMessage;
    chat_message: ServerMessage;
    bulk_chat_messages: ServerMessage;
    message_count_by_channel: ServerMessage;
    session_activity: ServerMessage;
  };