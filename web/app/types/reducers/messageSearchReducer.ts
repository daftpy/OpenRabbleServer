// An enumeration of possible action types
export enum MessageSearchActions {
  SetKeyword,
  SetTemporaryKeyword,
  AddFilter,
  RemoveFilter,
  SetMessages,
  NextPage,
  PrevPage,
  ExecuteSearch,
  SearchFinished
}

// The actions that can be performed.
type SetKeywordAction = { type: MessageSearchActions.SetKeyword; keyword: string };
type SetTemporaryKeywordAction = {type: MessageSearchActions.SetTemporaryKeyword; keyword: string}
type AddFilterAction = { type: MessageSearchActions.AddFilter; filter: string };
type RemoveFilterAction = { type: MessageSearchActions.RemoveFilter; filter: string };
type SetMessagesAction = { type: MessageSearchActions.SetMessages; messages: any, hasMore: boolean };
type NextPageAction = { type: MessageSearchActions.NextPage; }
type PrevPageAction = { type: MessageSearchActions.PrevPage; }
type ExecuteSearchAction = {type: MessageSearchActions.ExecuteSearch };
type SearchFinished = {type: MessageSearchActions.SearchFinished };

export type MessageSearchAction = SetKeywordAction | SetTemporaryKeywordAction | AddFilterAction | RemoveFilterAction | SetMessagesAction | NextPageAction | PrevPageAction | ExecuteSearchAction | SearchFinished;

/*
  The search state is simple. It contains a keyword and activeFilters collection.
  If the keyword is present (not, ""), it is used as a search term. If activeFilters
  is not empty, the filters are used to further filter messages by channel.
*/
export type MessageSearchState = {
  activeFilters: string[]
  keyword: string;
  temporaryKeyword: string;
  availableFilters: string[];
  messages: any;
  has_more: boolean;
  page: number;
  searching: boolean;
}
