import type { Route } from "./+types/index";
import { useLoaderData } from "react-router";
import { useEffect } from "react";
import RouteProtector from "~/components/route_protector";
import { ChannelPage } from "~/pages/channels";

export async function loader({ params }: Route.LoaderArgs) {
  const response = await fetch("https://chat.localhost/channels");
  
  if (!response.ok) {
    throw new Response("Failed to load channels", { status: response.status });
  }

  const data = await response.json();
  console.log(data);
  if (data.channels == null) {
    return [];
  }
  return data.channels;
}

export async function clientAction({ request }: Route.ActionArgs) {
  const response = await fetch("https://chat.localhost/channels");
  
  if (!response.ok) {
    throw new Response("Failed to fetch channels", { status: response.status });
  }
  
  const channelsData = await response.json();
  return { channels: channelsData.channels ? channelsData.channels : [] };
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Welcome to React Router!" },
  ];
}

export default function ChannelRoute({loaderData,}: Route.ComponentProps) {
  const channels = useLoaderData();
  useEffect(() => {
    console.log("Test home");
  }, []);

  return (
    <RouteProtector>
      <ChannelPage channels={channels} />
    </RouteProtector>
  );
}
