import { Box, Button, DropdownMenu, Flex, Text } from "@radix-ui/themes";
import { Message, type MessageType } from "./message";
import { useReducer, memo } from "react";
import { CheckIcon, Cross2Icon } from "@radix-ui/react-icons";
import { useFetcher } from "react-router";

export type MessageListType = Omit<MessageType, "isLast">;

type Props = {
  messages: MessageListType[];
  hidePermaLink: boolean;
};

type MessageSelectState = {
  selected: number[];
};

type Action =
  | { type: "select_message"; id: number }
  | { type: "select_messages"; ids: number[] };

function reducer(state: MessageSelectState, action: Action) {
  switch (action.type) {
    case "select_message":
      // Example: toggling selection if already selected, otherwise adding it
      return state.selected.includes(action.id)
        ? { ...state, selected: state.selected.filter(x => x !== action.id) }
        : { ...state, selected: [...state.selected, action.id] };
    case "select_messages":
      return {...state, selected: [...state.selected, ...action.ids]}

    default:
      return state;
  }
}

export const MessageList = memo(({ messages, hidePermaLink }: Props) => {
  const [state, dispatch] = useReducer(reducer, { selected: [] });
  const messageFetcher = useFetcher({key: "my-key"});

  const selectMessage = (id: number) => {
    dispatch({ type: "select_message", id });
  };

  const selectAllMessages = () => {
    const ids = messages.map((msg) => msg.id);
    
    dispatch({ type: "select_messages", ids });
  };

  const deleteMessages = () => {
    if (state.selected.length === 0) {
      return; // nothing to delete
    }
  
    // Submit the selected IDs as a JSON string in a form field
    messageFetcher.submit(
      { ids: JSON.stringify(state.selected) },
      {
        method: "delete",
        action: "/messages",
      }
    );
  };

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

      {messages.length > 0 ? (messages.map((message, index) => (
        <Message
          key={index}
          isLast={index === messages.length - 1}
          meessage={message}
          hidePermaLink={hidePermaLink}
          // Pass a boolean to indicate whether this message is in the selected array
          isSelected={state.selected.includes(message.id)}
          onSelect={() => selectMessage(message.id)}
        />
      ))) : (
        <Box py={"4"} key="none">
          <Text align={"center"} as="div" size="5" style={{color: "var(--muted-text-color)"}}>No Messages</Text>
        </Box>
      )}
    </Box>
  );
});
