import { Cross1Icon, MagnifyingGlassIcon } from "@radix-ui/react-icons";
import { Badge, Button, DropdownMenu, Flex, Text, TextField } from "@radix-ui/themes";
import { useEffect, useMemo } from "react";
import { useFetcher } from "react-router";
import { MessageSearchActionType, type MessageSearchAction, type MessageSearchState } from "~/hooks/useMessageSearch";
import type { Channel } from "~/types/components/channel";

export function MessageSearchInput({
  state,
  dispatch,
} : {
  state: MessageSearchState;
  dispatch: React.Dispatch<MessageSearchAction>;
}) {
  const channelFetcher = useFetcher();

  const filters = useMemo(() => {
    if (channelFetcher.data?.channels) {
      return channelFetcher.data.channels.map((channel: Channel) => channel.name);
    }
    return [];
  }, [channelFetcher.data]);

  useEffect(() => {
    channelFetcher.load("/channels");
  }, []);

  useEffect(() => {
    console.log("Wow wow wow, state changed", state.messages);
  }, [state.messages]);
  return (
    <Flex direction={"column"} gap={"1"}>
      <Flex gap={"4"} py={"2"}>
        <TextField.Root
            placeholder="keyword"
            className="grow"
            value={state.temporaryKeyword}
            onChange={(e) => dispatch({ type: MessageSearchActionType.SetTemporaryKeyword, keyword: e.target.value })}
          />
          <DropdownMenu.Root>
            <DropdownMenu.Trigger>
              <Button color="amber">Filter</Button>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content>
              {filters.length > 0 && filters.map((filter: string) => (
                <DropdownMenu.Item key={filter} onClick={() => dispatch({ type: MessageSearchActionType.AddFilter, filter })}>
                  { filter }
                </DropdownMenu.Item>
              ))}
            </DropdownMenu.Content>
          </DropdownMenu.Root>
          <Button color="blue" onClick={() => dispatch({type: MessageSearchActionType.ExecuteSearch})}><MagnifyingGlassIcon />Search</Button>
      </Flex>
      <Flex gap={"2"} align={"baseline"}>
        <Text weight={"bold"} style={{ color: "var(--indigo-12)" }}>Filters: </Text>
        {state.activeFilters && state.activeFilters.map((filter: string) => (
          <Badge key={filter} size={"2"} color="tomato" onClick={(e) => dispatch({type: MessageSearchActionType.RemoveFilter, filter: filter})}>
            {filter} <Cross1Icon />
          </Badge>
        ))}
      </Flex>
    </Flex>
  )
}
