import type { Message } from "../components/message";

export type FetchMessagesResponse = {
    type: string;
    sender: string;
    payload: {
        messages: Message[];
        has_more: boolean;
    }
}
