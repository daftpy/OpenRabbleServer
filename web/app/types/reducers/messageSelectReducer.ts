export type MessageSelectState = {
  selected: number[];
};

export enum MessageSelectActions {
  SELECT_MESSAGE = "select_message",
  SELECT_MESSAGES = "select_messages"
}

export type MessageSelectAction =
  | { type: MessageSelectActions.SELECT_MESSAGE; id: number }
  | { type: MessageSelectActions.SELECT_MESSAGES; ids: number[] };

