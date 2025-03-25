import type { Route } from "./+types/channel";
import { useEffect } from "react";
import { editChannel, fetchChannels, type EditChannelPayload } from "~/api/channels";
import RouteProtector from "~/components/route_protector";
import { ChannelPage } from "~/pages/channels";

export async function loader({ params }: Route.LoaderArgs) {
  return fetchChannels();
}

export async function clientLoader({
  serverLoader,
  params,
}: Route.ClientLoaderArgs) {
  
  const serverData = await serverLoader();

  const activityRes = await fetch(`https://chat.localhost/activity/channels`)
  const activityData = await activityRes.json();
  console.log("ACTIVITY: ", serverData);

  return { channels: serverData, channel_activity: activityData.payload.channels };
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Welcome to React Router!" },
  ];
}

// force the client loader to run during hydration
clientLoader.hydrate = true as const; // `as const` for type inference

export function HydrateFallback() {
  return <div>Loading...</div>;
}

// src/routes/channels/update.ts
export async function clientAction({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const intent = formData.get("intent");

  const payload: EditChannelPayload = {
    id: parseInt(formData.get("id") as string),
    name: formData.get("name")?.toString() ?? null,
    description: formData.get("description")?.toString() ?? null,
  };

  return await editChannel(payload);
}

export default function ChannelRoute({loaderData,}: Route.ComponentProps) {
  const {channels, channel_activity } = loaderData;
  useEffect(() => {
    console.log("Test home");
  }, [channels]);

  return (
    <RouteProtector>
      <ChannelPage channels={channels} channelActivity={channel_activity} />
    </RouteProtector>
  );
}
