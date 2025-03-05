import { useEffect } from "react";
import type { Route } from "./+types/search_messages";
import { SearchMessagesPage } from "~/pages/search_messages"
import { useLoaderData } from "react-router";

export async function loader({ params }: Route.LoaderArgs) {
  const response = await fetch("https://chat.localhost/messages?channel=General&limit=10&offset=0");
  
  if (!response.ok) {
    throw new Response("Failed to load channels", { status: response.status });
  }

  const data = await response.json();
  console.log(data);
  if (data.messages == null) {
    return [];
  }
  return data.messages;
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Search Messages" },
  ];
}

export default function About({loaderData,}: Route.ComponentProps) {
  const messages = useLoaderData();
  useEffect(() => {
    console.log("Test search page");
  }, []);
  return (
    <SearchMessagesPage messages={messages} />
  )
}
