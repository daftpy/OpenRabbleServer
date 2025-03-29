import { useCallback, useEffect, useReducer } from "react";
import { useFetcher } from "react-router";
import { MessageSearchActions, type MessageSearchAction, type MessageSearchState } from "~/types/reducers/messageSearchReducer";

function reducer(state: MessageSearchState, action: MessageSearchAction) {
  switch (action.type)  {
    case MessageSearchActions.SetKeyword: 
      return { ...state, keyword: action.keyword, page: 0 };
    case MessageSearchActions.SetTemporaryKeyword:
      return { ...state, temporaryKeyword: action.keyword }
    case MessageSearchActions.AddFilter:
      if (!state.activeFilters.includes(action.filter)) {
        return { ...state, activeFilters: [...state.activeFilters, action.filter]};
      }
      return state;
    case MessageSearchActions.RemoveFilter:
      return { ...state, activeFilters: state.activeFilters.filter((filter: string) => filter !== action.filter) };
    case MessageSearchActions.SetMessages:
      return { ...state, messages: action.messages, has_more: action.hasMore }
    case MessageSearchActions.NextPage:
      return { ...state, page: state.page + 1, searching: true };
    case MessageSearchActions.PrevPage:
      if (state.page <= 0) return state; // Return if already at page 0
      return { ...state, page: state.page - 1, searching: true };
    case MessageSearchActions.ExecuteSearch:
      return { ...state, keyword: state.temporaryKeyword, page: 0, searching: true };
    case MessageSearchActions.SearchFinished: 
      return { ...state, searching: false};
    default:
      return state;
  }
}

export function useMessageSearch({userId, messages, hasMore } : {userId?: string, messages: any, hasMore: boolean}) {
  const [state, dispatch] = useReducer(reducer, { keyword: "", temporaryKeyword: "", activeFilters: [], availableFilters: [], messages: messages, has_more: hasMore, page: 0, searching: false });
  // Used to fetch messages
  const messageFetcher = useFetcher({key: "my"});

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
      console.log("Messages fetched?");
      const newMessages = messageFetcher.data.messages;
      const hasMore = messageFetcher.data.has_more;

      // If no messages and we're not on page 0, go back a page
      if (newMessages.length === 0 && state.page > 0) {
        dispatch({ type: MessageSearchActions.PrevPage });
        return;
      }

      dispatch({ type: MessageSearchActions.SetMessages, messages: newMessages, hasMore: hasMore });
    }
  }, [messageFetcher.data]);

  const nextPage = useCallback(() => {
    dispatch({ type: MessageSearchActions.NextPage });
  }, []);

  const prevPage = useCallback(() => {
    dispatch({ type: MessageSearchActions.PrevPage });
  }, []);

  // When the page or keyword changes, execute a search
  useEffect(() => {
    if (state.searching) {
      searchMessages(10, state.page * 10);
      dispatch({type: MessageSearchActions.SearchFinished});
    }
  }, [state.searching]);

  return {
    state,
    messageFetcher,
    dispatch,
    searchMessages,
    nextPage,
    prevPage
  }
}
