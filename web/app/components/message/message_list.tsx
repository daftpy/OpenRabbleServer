import { Flex, Text } from "@radix-ui/themes";
import { Message, type MessageType } from "./message";
// Remove the 'isLast' field from the MessageType
export type MessageListType = Omit<MessageType, "isLast">;

type Props = {
  messages: MessageListType[];
  hidePermaLink: boolean;
};

export function MessageList({ messages, hidePermaLink }: Props) {

  return (
    <>
      {messages ? messages.map((message, index) => (
        <Message
          key={index}
          isLast={index === messages.length - 1}
          meessage={message}
          hidePermaLink={hidePermaLink}
        />
      )) : (
        <Flex justify={"center"} style={{color: "var(--muted-text-color)"}}>
          <Text weight={"bold"}>No Messages</Text>
        </Flex>
      )}
    </>
  );
}
