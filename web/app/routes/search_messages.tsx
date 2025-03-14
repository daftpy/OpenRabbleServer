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
    const rawIds = formData.get("ids");

    console.log("Deleting message IDs", rawIds);

    if (!rawIds) {
      console.log("Missing message ids");
      return new Response("Missing message IDs", { status: 400 });
    }

    // Parse the JSON string into an array
    let messageIds: number[];
    try {
      messageIds = JSON.parse(rawIds.toString()); 
      // At this point, messageIds could be something like [17,16].
    } catch (err) {
      console.log("Failed to parse 'ids' as JSON");
      return new Response("Invalid JSON for 'ids'", { status: 400 });
    }

    // Validate it's a non-empty array
    if (!Array.isArray(messageIds) || messageIds.length === 0) {
      console.log("No valid IDs in array");
      return new Response("No valid IDs", { status: 400 });
    }

    // Build the JSON body for bulk deletion
    const body = JSON.stringify({ ids: messageIds });

    // Send the DELETE request
    const response = await fetch("https://chat.localhost/messages", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body,
    });

    // 5) Check response
    if (!response.ok) {
      console.log("Failed to delete messages");
      return new Response("Failed to delete messages", { status: 500 });
    }

    console.log("Messages should be deleted");
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
