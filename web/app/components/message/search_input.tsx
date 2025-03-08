import { Cross1Icon, MagnifyingGlassIcon } from "@radix-ui/react-icons";
import { Badge, Button, DropdownMenu, Flex, Text, TextField } from "@radix-ui/themes";
import { useEffect, useReducer, useState } from "react";
import { useFetcher } from "react-router";

function reducer(state : any, action: any) {
  switch (action.type) {
    case "set_keyword":
      return {...state, keyword: action.keyword}
    case "add_filter":
      console.log("Adding filter:", action.filter);
      return {...state, activeFilters: [...state.activeFilters, action.filter]}
    case "remove_filter":
      return {...state, activeFilters: [...state.activeFilters.filter((filter : any ) => filter != action.filter)]}
    default:
      return state;
  }
}

export function MessageSearchInput({
  userId,
  onMessagesUpdate 
} : {
  userId?: string,
  onMessagesUpdate: (messages: any) => void;
}) {
  const [state, dispatch] = useReducer(reducer, {keyword: "", activeFilters: []});
  const [availableFilters, setAvailableFilters] = useState<string[]>([]);
  const messageFetcher = useFetcher();
  const channelFetcher = useFetcher();

  // Search messages using keywords and filters if available
  const handleSearch = () => {
    let formData = new FormData(); // Create the new form

    // Add the conditional parameters for the search
    if (state.keyword) {
      formData.append("keyword", state.keyword);
    }
    if (userId) {
      formData.append("user_id", userId);
    }
    state.activeFilters.forEach((filter: any) => {
      formData.append("channel", filter.name);
    })

    messageFetcher.submit(formData, {method: "post", action: "/messages"});
  }

  // Fetch the server channels to use as filters
  useEffect(() => {
    channelFetcher.submit({}, {method: "post", action: "/channels"});
  }, [])

  // Retrieve the server channels and set them as available fiilters
  useEffect(() => {
    if (channelFetcher.data?.channels) {
      setAvailableFilters(channelFetcher.data.channels);
    }
  }, [channelFetcher.data])

  // Update the fetchedMessages when the handleSearch is over
  useEffect(() => {
    console.log("NEW searched meessages:", messageFetcher.data);
    if (messageFetcher.data?.messages) {
      onMessagesUpdate(messageFetcher.data.messages);
    }
  }, [messageFetcher.data]);

  return (
    <Flex direction={"column"} gap={"1"}>
      <Flex gap={"4"} py={"2"}>
        <TextField.Root
            placeholder="keyword"
            className="grow"
            value={state.keyword}
            onChange={(e) => dispatch({type: "set_keyword", keyword: e.target.value})}
          />
          <DropdownMenu.Root>
            <DropdownMenu.Trigger>
              <Button color="amber">Filter</Button>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content>
              {availableFilters.length > 0 && availableFilters.map((filter: any) => (
                <DropdownMenu.Item key={filter.name} onClick={(e) => dispatch({ type: "add_filter", filter: filter })}>
                  { filter.name }
                </DropdownMenu.Item>
              ))}
            </DropdownMenu.Content>
          </DropdownMenu.Root>
          <Button color="blue" onClick={() => handleSearch()}><MagnifyingGlassIcon />Search</Button>
      </Flex>
      <Flex gap={"2"} align={"baseline"}>
        <Text weight={"bold"} style={{ color: "var(--indigo-12)" }}>Filters: </Text>
        {state.activeFilters && state.activeFilters.map((filter: any) => (
          <Badge key={filter.name} size={"2"} color="tomato" onClick={(e) => dispatch({type: "remove_filter", filter: filter})}>
            {filter.name} <Cross1Icon />
          </Badge>
        ))}
      </Flex>
    </Flex>
  )
}