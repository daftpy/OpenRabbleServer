import type { Channel } from "./components/channel_list";
import type { ChatMessageType } from "./components/chat_message_list";

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

export interface BulkChatMessages extends Message {
  type: "bulk_chat_messages";
  // Use the meessage type for react state
  messages: ChatMessageType[];
}
  
  // Create a union type for all server messages
export type ServerMessage =
  | ConnectedUsersMessage
  | UserStatusMessage
  | ChatMessage
  | BulkChatMessages
  | ActiveChannelsMessage;

  export type EmitterEvents = {
    connected_users: ServerMessage;
    user_status: ServerMessage;
    active_channels: ServerMessage;
    chat_message: ServerMessage;
    bulk_chat_messages: ServerMessage;
  };