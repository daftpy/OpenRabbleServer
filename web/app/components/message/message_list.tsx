import { Message, type MessageType } from "./message";
// Remove the 'isLast' field from the MessageType
export type MessageListType = Omit<MessageType, "isLast">;

type Props = {
  messages: MessageListType[];
};

export function MessageList({ messages }: Props) {

  return (
    <>
      {messages.map((message, index) => (
        <Message
          key={index}
          username={message.username}
          channel={message.channel}
          message={message.message}
          authored_at={message.authored_at}
          isLast={index === messages.length - 1}
        />
      ))}
    </>
  );
}
