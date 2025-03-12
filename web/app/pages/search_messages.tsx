import { Box, Button, Flex, Heading, Text } from "@radix-ui/themes";
import { Link } from "react-router";
import { MessageList } from "~/components/message/message_list";
import { MessageSearchInput } from "~/components/message/search_input";
import { useMessageSearch } from "~/hooks/useMessageSearch";

export function SearchMessagesPage({ messages, hasMore }: {messages: any, hasMore: boolean}) {
  const { state, dispatch, nextPage, prevPage } = useMessageSearch({ messages, hasMore });
  console.log("SearchMessagesPage says hello");
  return (
    <Flex direction={"column"} maxWidth={"990px"} m={"auto"} px={"4"} py={"6"}>
      <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{ color: "var(--slate-12)" }}>
        <Link to="/">Your OnRabble Server</Link>
      </Heading>
      <Box pt={"4"}>
        <Heading style={{color: "var(--indigo-10)"}}>Search</Heading>
        <Text>Search messages stored in your cache or database. Filter by <strong style={{ color: "var(--link-color)" }}>channel </strong>or <strong style={{ color: "var(--link-color)" }}>keyword</strong>.</Text>
      </Box>
      {/* Search Input & Filters */}
      <MessageSearchInput state={state} dispatch={dispatch} />

      {/* Show messages */}
      <MessageList messages={state.messages} hidePermaLink={true} />
      <Flex justify={"between"} pt={"4"}>
      <Button 
          disabled={state.page === 0} 
          onClick={() => prevPage()}
        >
          Previous Page
        </Button>
        <Button 
          disabled={!state.hasMore}
          onClick={() => nextPage()}
        >
          Next Page
        </Button>
      </Flex>
    </Flex>
  );
}
