import { useEffect } from "react";
import type { Route } from "./+types/search_messages";
import { SearchMessagesPage } from "~/pages/search_messages"
import { useLoaderData, type ShouldRevalidateFunctionArgs } from "react-router";
import { fetchMessagesFromAPI } from "~/api/fetchMessages";
import { deleteMessagesFromAPI } from "~/api/deleteMessages";

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

export async function clientAction({ request }: Route.ActionArgs) {
  if (request.method === "DELETE") {
    const formData = await request.formData();
    const rawIds = formData.get("ids");

    if (!rawIds) {
      return new Response("Missing message IDs", { status: 400 });
    }

    let messageIds: number[];
    try {
      messageIds = JSON.parse(rawIds.toString()); 
    } catch {
      return new Response("Invalid JSON for 'ids'", { status: 400 });
    }

    const response = await deleteMessagesFromAPI(messageIds);

    // Now response is a real Response object from fetch().
    if (!response.ok) {
      return new Response("Failed to delete messages", { status: 500 });
    }

    // Otherwise, success
    return new Response(null, { status: 204 });
  }

  return new Response("Method Not Allowed", { status: 405 });
}


export function meta({}: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Search Messages" },
  ];
}

export default function About({loaderData,}: Route.ComponentProps) {
  const { messages, hasMore } = useLoaderData() as { messages: any[], hasMore: boolean };
  return (
    <SearchMessagesPage messages={messages} hasMore={hasMore} />
  )
}
