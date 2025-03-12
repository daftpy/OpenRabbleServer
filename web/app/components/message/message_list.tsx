import { Flex, Text } from "@radix-ui/themes";
import { Message, type MessageType } from "./message";
import { useEffect } from "react";
// Remove the 'isLast' field from the MessageType
export type MessageListType = Omit<MessageType, "isLast">;

type Props = {
  messages: MessageListType[];
  hidePermaLink: boolean;
};

import { memo } from "react";

export const MessageList = memo(({ messages, hidePermaLink }: Props) => {
  console.log("MessageList re-rendered", messages);
  return (
    <>
      {messages.map((message, index) => (
        <Message
          key={index}
          isLast={index === messages.length - 1}
          meessage={message}
          hidePermaLink={hidePermaLink}
        />
      ))}
    </>
  );
});
