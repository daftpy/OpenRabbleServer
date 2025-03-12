import { useCallback, useEffect, useReducer, useState } from "react";
import { useFetcher } from "react-router";

// An enumeration of possible action types
export enum MessageSearchActionType {
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
  hasMore: boolean;
  page: number;
  searching: boolean;
}

// The actions that can be performed.
type SetKeywordAction = { type: MessageSearchActionType.SetKeyword; keyword: string };
type SetTemporaryKeywordAction = {type: MessageSearchActionType.SetTemporaryKeyword; keyword: string}
type AddFilterAction = { type: MessageSearchActionType.AddFilter; filter: string };
type RemoveFilterAction = { type: MessageSearchActionType.RemoveFilter; filter: string };
type SetMessagesAction = { type: MessageSearchActionType.SetMessages; messages: any, hasMore: boolean };
type NextPageAction = { type: MessageSearchActionType.NextPage; }
type PrevPageAction = { type: MessageSearchActionType.PrevPage; }
type ExecuteSearchAction = {type: MessageSearchActionType.ExecuteSearch };
type SearchFinished = {type: MessageSearchActionType.SearchFinished };

export type MessageSearchAction = SetKeywordAction | SetTemporaryKeywordAction | AddFilterAction | RemoveFilterAction | SetMessagesAction | NextPageAction | PrevPageAction | ExecuteSearchAction | SearchFinished;

function reducer(state: MessageSearchState, action: MessageSearchAction) {
  switch (action.type)  {
    case MessageSearchActionType.SetKeyword: 
      return { ...state, keyword: action.keyword, page: 0 };
    case MessageSearchActionType.SetTemporaryKeyword:
      return { ...state, temporaryKeyword: action.keyword }
    case MessageSearchActionType.AddFilter:
      return { ...state, activeFilters: [...state.activeFilters, action.filter]};
    case MessageSearchActionType.RemoveFilter:
      return { ...state, activeFilters: state.activeFilters.filter((filter: string) => filter !== action.filter) };
    case MessageSearchActionType.SetMessages:
      return { ...state, messages: action.messages, hasMore: action.hasMore }
    case MessageSearchActionType.NextPage:
      return { ...state, page: state.page + 1, searching: true };
    case MessageSearchActionType.PrevPage:
      return { ...state, page: state.page - 1, searching: true };
    case MessageSearchActionType.ExecuteSearch:
      return { ...state, keyword: state.temporaryKeyword, page: 0, searching: true };
    case MessageSearchActionType.SearchFinished: 
      return { ...state, searching: false};
    default:
      return state;
  }
}

export function useMessageSearch({userId, messages, hasMore } : {userId?: string, messages: any, hasMore: boolean}) {
  const [state, dispatch] = useReducer(reducer, { keyword: "", temporaryKeyword: "", activeFilters: [], availableFilters: [], messages: messages, hasMore: hasMore, page: 0, searching: false });
  // Used to fetch messages
  const messageFetcher = useFetcher();

  // Perform the search
  const searchMessages = useCallback((limit = 10, offset = state.page * limit) => {
    console.log(`Search messages with ${limit} limit and ${offset} offset`);
    const params = new URLSearchParams();
    
    if (userId) {
      params.append("user_id", userId);
    }
    
    if (state.keyword) {
      params.append("keyword", state.keyword);
    }
  
    state.activeFilters.forEach((filter) => {
      params.append("channel", filter);
    });
  
    params.append("limit", limit.toString());
    params.append("offset", offset.toString());
  
    messageFetcher.load(`/messages?${params.toString()}`);
  }, [userId, state.keyword, state.activeFilters, state.page, messageFetcher]);
  

  useEffect(() => {
    // If data was received, update the messages
    if (messageFetcher.data?.messages) {
      dispatch({ type: MessageSearchActionType.SetMessages, messages: messageFetcher.data.messages, hasMore: messageFetcher.data.hasMore });
    }
  }, [messageFetcher.data]);

  const nextPage = () => {
    dispatch({ type: MessageSearchActionType.NextPage });
  }

  const prevPage = () => {
    dispatch({ type: MessageSearchActionType.PrevPage });
  }

  // When the page or keyword changes, execute a search
  useEffect(() => {
    if (state.searching) {
      searchMessages(10, state.page * 10);
      dispatch({type: MessageSearchActionType.SearchFinished});
    }
  }, [state.searching]);

  return {
    state,
    dispatch,
    searchMessages,
    nextPage,
    prevPage
  }
}
