import { fetchSessionsActivity } from "~/api/activity";
import type { Route } from "./+types/index";
import { useEffect } from "react";
import { fetchChannels } from "~/api/channels";
import RouteProtector from "~/components/route_protector";
import { HomePage } from "~/pages/home";

export type SessionActivity = {
  session_date: string; // e.g., "2025-02-23"
  session_count: number; // Number of sessions for that day
  total_duration: string; // e.g., "15 hours 30 minutes"
};

export async function loader({ params }: Route.LoaderArgs) {

}

export async function clientLoader({
  serverLoader,
  params,
}: Route.ClientLoaderArgs) {
  
  const serverData = await fetchChannels();

  const activityData = await fetchSessionsActivity()

  return { channels: serverData, session_activity: activityData.session_activity };
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
  const {channels, session_activity} = loaderData;
  useEffect(() => {
    console.log("Test home", channels);
  }, []);

  return (
    <RouteProtector>
      <HomePage channels={channels} session_activity={session_activity} />
    </RouteProtector>
  );
}
