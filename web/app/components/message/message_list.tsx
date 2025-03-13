import { Button, DropdownMenu, Flex, Text } from "@radix-ui/themes";
import { Message, type MessageType } from "./message";
import { useEffect } from "react";
// Remove the 'isLast' field from the MessageType
export type MessageListType = Omit<MessageType, "isLast">;

type Props = {
  messages: MessageListType[];
  hidePermaLink: boolean;
};

import { memo } from "react";
import { CheckIcon, Cross2Icon } from "@radix-ui/react-icons";

export const MessageList = memo(({ messages, hidePermaLink }: Props) => {
  console.log("MessageList re-rendered", messages);
  return (
    <>
      <Flex justify={"end"} pb={"3"}>
        <DropdownMenu.Root>
          <DropdownMenu.Trigger>
            <Button color="tomato">Action</Button>
          </DropdownMenu.Trigger>
          <DropdownMenu.Content>
            <DropdownMenu.Item>
              <CheckIcon/> Select All
            </DropdownMenu.Item>
            <DropdownMenu.Item color="tomato">
              <Cross2Icon /> Delete Selected
            </DropdownMenu.Item>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </Flex>
      {messages && messages.map((message, index) => (
        <Message
          key={message.id}
          isLast={index === messages.length - 1}
          meessage={message}
          hidePermaLink={hidePermaLink}
        />
      ))}
    </>
  );
});
