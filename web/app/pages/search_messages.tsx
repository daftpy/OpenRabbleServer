import { Box, Button, Container, Flex, Heading, Text } from "@radix-ui/themes";
import { Link } from "react-router";
import { MessageList } from "~/components/message/message_list";
import { MessageSearchInput } from "~/components/message/search_input";
import { useMessageSearch } from "~/hooks/useMessageSearch";

export function SearchMessagesPage({ messages, hasMore }: {messages: any, hasMore: boolean}) {
  const { state, messageFetcher, dispatch, nextPage, prevPage } = useMessageSearch({ messages, hasMore });
  console.log("SearchMessagesPage says hello", hasMore);
  return (
    <>
      <Flex direction={"column"} gap={"4"}>
        <Box>
          <Heading color="indigo">Search</Heading>
          <Text>Search messages stored in your cache or database. Filter by <strong style={{ color: "var(--link-color)" }}>channel </strong>or <strong style={{ color: "var(--link-color)" }}>keyword</strong>.</Text>
          {/* Search Input & Filters */}
          <MessageSearchInput state={state} dispatch={dispatch} />
        </Box>

        {/* Show messages */}
        <Box>
          <MessageList messages={messageFetcher.data?.messages ? messageFetcher.data.messages : messages} hidePermaLink={true} />
        </Box>
        <Flex justify={"between"} pt={"4"}>
        <Button 
            disabled={state.page === 0} 
            onClick={() => prevPage()}
          >
            Previous Page
          </Button>
          <Button 
            disabled={!state.has_more}
            onClick={() => nextPage()}
          >
            Next Page
          </Button>
        </Flex>
      </Flex>
    </>
  );
}
