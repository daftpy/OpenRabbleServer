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

export async function clientAction({ request }: Route.ActionArgs) {
  if (request.method === "DELETE") {
    const formData = await request.formData();
    const messageId = formData.get("id");

    if (!messageId) {
      return new Response("Missing message ID", { status: 400 });
    }

    const response = await fetch(`https://chat.localhost/messages?id=${messageId}`, {
      method: "DELETE",
    });

    if (!response.ok) {
      return new Response("Failed to delete message", { status: 500 });
    }

    return new Response(null, { status: 204 }); // No Content (Success)
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
