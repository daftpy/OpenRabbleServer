import type { Route } from "./+types/index";
import { useLoaderData } from "react-router";
import { useEffect } from "react";
import RouteProtector from "~/components/route_protector";
import { HomePage } from "~/pages/home";

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

export async function clientLoader({
  serverLoader,
  params,
}: Route.ClientLoaderArgs) {
  
  const serverData = await serverLoader();

  const activityRes = await fetch(`https://chat.localhost/activity/sessions`)
  const activityData = await activityRes.json();
  console.log("ACTIVITY: ", serverData);

  return { channels: serverData, ...activityData.payload };
}

// force the client loader to run during hydration
clientLoader.hydrate = true as const; // `as const` for type inference

export function HydrateFallback() {
  return <div>Loading...</div>;
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Welcome to React Router!" },
  ];
}

export default function Index({loaderData,}: Route.ComponentProps) {
  const {channels, session_activity} = useLoaderData();
  useEffect(() => {
    console.log("Test home", channels);
  }, []);

  return (
    <RouteProtector>
      <HomePage channels={channels} session_activity={session_activity} />
    </RouteProtector>
  );
}
