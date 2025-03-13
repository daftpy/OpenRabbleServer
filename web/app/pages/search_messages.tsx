import { CheckIcon, Cross2Icon, MagicWandIcon } from "@radix-ui/react-icons";
import { Box, Button, Container, DropdownMenu, Flex, Heading, Text } from "@radix-ui/themes";
import { Link } from "react-router";
import { MessageList } from "~/components/message/message_list";
import { MessageSearchInput } from "~/components/message/search_input";
import { useMessageSearch } from "~/hooks/useMessageSearch";

export function SearchMessagesPage({ messages, hasMore }: {messages: any, hasMore: boolean}) {
  const { state, messageFetcher, dispatch, nextPage, prevPage } = useMessageSearch({ messages, hasMore });
  console.log("SearchMessagesPage says hello");
  return (
    <Container px={"4"} py={"6"}>
      <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{ color: "var(--slate-12)" }}>
        <Link to="/">Your OnRabble Server</Link>
      </Heading>
      <Flex direction={"column"} gap={"4"}>
        <Box pt={"4"}>
          <Heading style={{color: "var(--indigo-10)"}}>Search</Heading>
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
            disabled={!state.hasMore}
            onClick={() => nextPage()}
          >
            Next Page
          </Button>
        </Flex>
      </Flex>
    </Container>
  );
}
