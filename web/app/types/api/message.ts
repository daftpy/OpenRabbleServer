import type { Message } from "../components/message";

// Response after fetching messages
export type FetchMessagesResponse = {
    type: string;
    sender: string;
    payload: {
        messages: Message[];
        has_more: boolean;
    }
}
