import type { Route } from "./+types/index";
import { useEffect } from "react";
import { fetchChannels } from "~/api/fetchChannels";
import RouteProtector from "~/components/route_protector";
import { HomePage } from "~/pages/home";

export type SessionActivity = {
  session_date: string; // e.g., "2025-02-23"
  session_count: number; // Number of sessions for that day
  total_duration: string; // e.g., "15 hours 30 minutes"
};

type SessionActivityResult = {
  type: string;
  sender: string;
  payload: {
    session_activity: SessionActivity[];
  };
};


export async function loader({ params }: Route.LoaderArgs) {
  return fetchChannels();
}

export async function clientLoader({
  serverLoader,
  params,
}: Route.ClientLoaderArgs) {
  
  const serverData = await serverLoader();

  const activityRes = await fetch(`https://chat.localhost/activity/sessions`)
  const activityData : SessionActivityResult = await activityRes.json();
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
