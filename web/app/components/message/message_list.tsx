import { Box, Button, DropdownMenu, Flex, Text } from "@radix-ui/themes";
import { memo } from "react";
import { CheckIcon, Cross2Icon } from "@radix-ui/react-icons";
import type { Message } from "~/types/components/message";
import { MessageRow, type UnifiedMessage } from "./message_row";
import { useMessageSelection } from "~/hooks/useMessageSelection";

type Props = {
  // messages: Message[] | ChatMessageType[];
  messages: UnifiedMessage[];
  hidePermaLink: boolean;
};

export const MessageList = memo(({ messages, hidePermaLink }: Props) => {
  // 1) A type guard to determine if Message or ChatMessageType (live messages)
  function isMessage(m: UnifiedMessage): m is Message {
    return m.id !== undefined;
  }

  // 2) Filter the messages with the id (Message type)
  const messagesWithId = messages.filter(isMessage); // now typed as Message[]

  const { selected, selectMessage, selectAllMessages, deleteMessages } = 
    useMessageSelection(messagesWithId);

  return (
    <Box>
      {hidePermaLink && (
        <Flex justify="end" pb="3">
          <DropdownMenu.Root>
            <DropdownMenu.Trigger>
              <Button>Action</Button>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content>
              <DropdownMenu.Item onClick={() => selectAllMessages()}>
                <CheckIcon /> Select All
              </DropdownMenu.Item>
              <DropdownMenu.Item color="tomato" onClick={() => deleteMessages()}>
                <Cross2Icon /> Delete Selected
              </DropdownMenu.Item>
            </DropdownMenu.Content>
          </DropdownMenu.Root>
        </Flex>
      )}

      {messages.length > 0 && messages ? (messages.map((message, index) => (
        <MessageRow
          key={index}
          isLast={index === messages.length - 1}
          meessage={message}
          hidePermaLink={hidePermaLink}
          // Pass a boolean to indicate whether this message is in the selected array
          isSelected={selected.includes(message.id?? 0)}
          onSelect={() => selectMessage(message.id?? 0)}
        />
      ))) : (
        <Box py={"4"} key="none">
          <Text align={"center"} as="div" size="5" style={{color: "var(--muted-text-color)"}}>No Messages</Text>
        </Box>
      )}
    </Box>
  );
});
