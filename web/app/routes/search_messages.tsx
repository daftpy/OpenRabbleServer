import { useEffect } from "react";
import type { Route } from "./+types/search_messages";
import { SearchMessagesPage } from "~/pages/search_messages"
import { useLoaderData } from "react-router";

export async function loader({ params }: Route.LoaderArgs) {
  const messageResponse = await fetch("https://chat.localhost/messages?channel=General&limit=20&offset=0");
  const channelResponse  = await fetch("https://chat.localhost/channels");

  
  if (!messageResponse.ok) {
    throw new Response("Failed to load messages", { status: messageResponse.status });
  }
  if (!channelResponse.ok) {
    throw new Response("Failed to load channels", { status: channelResponse.status });
  }

  const messageData = await messageResponse.json();
  const channelData = await channelResponse.json();
  console.log(messageData);
  console.log(channelData);
  
  return {
    messages: messageData.messages ?? [],
    channels: channelData.channels ?? [],
  };
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Search Messages" },
  ];
}

export default function About({loaderData,}: Route.ComponentProps) {
  const { messages, channels } = useLoaderData() as { messages: any[], channels: any[] };
  useEffect(() => {
    console.log("Test search page");
  }, []);
  return (
    <SearchMessagesPage messages={messages} channels={channels} />
  )
}
