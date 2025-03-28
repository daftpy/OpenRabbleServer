import type { MessageType } from "../components/message";

export type FetchMessagesResponse = {
    type: string;
    sender: string;
    payload: {
        messages: MessageType[];
        has_more: boolean;
    }
}
