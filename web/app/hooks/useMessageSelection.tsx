import { useReducer } from "react";
import { useFetcher } from "react-router";
import type { Message } from "~/types/components/message";
import { MessageSelectActions, type MessageSelectAction, type MessageSelectState } from "~/types/reducers/messageSelectReducer";

function messageSelectReducer(
    state: MessageSelectState,
    action: MessageSelectAction
): MessageSelectState {
  switch (action.type) {
    // If the message is already selected, de-select, otherwise, select.
    case MessageSelectActions.SELECT_MESSAGE:
      return state.selected.includes(action.id)
        ? { ...state, selected: state.selected.filter((x) => x !== action.id) }
        : { ...state, selected: [...state.selected, action.id] };
    // Select all
    case MessageSelectActions.SELECT_MESSAGES:
      return { ...state, selected: [...state.selected, ...action.ids] };
    default:
      return state;
  }
}

export function useMessageSelection(messages: Message[]) {
  // Start with an empty selection
  const [state, dispatch] = useReducer(messageSelectReducer, { selected: [] });
  const messageFetcher = useFetcher({ key: "my-key" });

  // Select a singular message
  const selectMessage = (id: number) => {
    dispatch({ type: MessageSelectActions.SELECT_MESSAGE, id });
  }

  // Select all messages
  const selectAllMessages = () => {
    const ids = messages.map((msg: Message) => msg.id);
    dispatch({ type: MessageSelectActions.SELECT_MESSAGES, ids });
  }

  // Delete messages in the selection
  const deleteMessages = () => {
    if (state.selected.length === 0) return;

    messageFetcher.submit(
      { ids: JSON.stringify(state.selected) },
      {
        method: "DELETE",
        action: "/messages"
      }
    );
  };

  return {
    selected: state.selected,
    selectMessage,
    selectAllMessages,
    deleteMessages
  };
}