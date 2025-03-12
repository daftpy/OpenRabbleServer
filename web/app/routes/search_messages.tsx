import { useEffect } from "react";
import type { Route } from "./+types/search_messages";
import { SearchMessagesPage } from "~/pages/search_messages"
import { useLoaderData, type ShouldRevalidateFunctionArgs } from "react-router";
import { fetchMessagesFromAPI } from "~/components/api/fetchMessages";

// TODO the default limit is hardcoded twice here and again in the SearchMessagesPage
// It is currently fragile and should possibly be refactored
export async function loader({ request }: Route.LoaderArgs) {
  const url = new URL(request.url);

  return fetchMessagesFromAPI({
    keyword: url.searchParams.get("keyword") ?? "",
    channels: url.searchParams.getAll("channel"),
    user_id: url.searchParams.get("user_id") ?? undefined,
    limit: url.searchParams.get("limit") ?? "10",
    offset: url.searchParams.get("offset") ?? "0"
  });
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Search Messages" },
  ];
}

export default function About({loaderData,}: Route.ComponentProps) {
  const { messages, channels, hasMore } = useLoaderData() as { messages: any[], channels: any[], hasMore: boolean };
  useEffect(() => {
    console.log("Test search page");
  }, []);
  return (
    <SearchMessagesPage messages={messages} hasMore={hasMore} />
  )
}
