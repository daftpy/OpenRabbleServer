import type { Channel } from "./components/channel_list";

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
  username: string;
  status: boolean;
}

export interface ChatMessage extends Message {
  type: "chat_message";
  message: string;
  username: string;
  channel: string;
}

export interface ActiveChannelsMessage extends Message {
  type: "active_channels";
  channels: Channel[]; 
}
  
  // Create a union type for all server messages
export type ServerMessage =
  | ConnectedUsersMessage
  | UserStatusMessage
  | ChatMessage
  | ActiveChannelsMessage;