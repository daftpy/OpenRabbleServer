import { Box, Button, DropdownMenu, Flex, Text } from "@radix-ui/themes";
import { Message, type MessageType } from "./message";
import { useReducer, memo } from "react";
import { CheckIcon, Cross2Icon } from "@radix-ui/react-icons";
import { useFetcher } from "react-router";
import { MessageSelectActions, type MessageSelectAction, type MessageSelectState } from "~/types/reducers/messageSelectReducer";

export type MessageListType = Omit<MessageType, "isLast">;

type Props = {
  messages: MessageListType[];
  hidePermaLink: boolean;
};

function reducer(state: MessageSelectState, action: MessageSelectAction) {
  switch (action.type) {
    case MessageSelectActions.SELECT_MESSAGE:
      // Example: toggling selection if already selected, otherwise adding it
      return state.selected.includes(action.id)
        ? { ...state, selected: state.selected.filter(x => x !== action.id) }
        : { ...state, selected: [...state.selected, action.id] };
    case MessageSelectActions.SELECT_MESSAGES:
      return {...state, selected: [...state.selected, ...action.ids]}

    default:
      return state;
  }
}

export const MessageList = memo(({ messages, hidePermaLink }: Props) => {
  const [state, dispatch] = useReducer(reducer, { selected: [] });
  const messageFetcher = useFetcher({key: "my-key"});

  const selectMessage = (id: number) => {
    dispatch({ type: MessageSelectActions.SELECT_MESSAGE, id });
  };

  const selectAllMessages = () => {
    const ids = messages.map((msg) => msg.id);
    
    dispatch({ type: MessageSelectActions.SELECT_MESSAGES, ids });
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
