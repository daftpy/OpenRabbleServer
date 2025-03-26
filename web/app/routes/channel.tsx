import type { Route } from "./+types/channel";
import { useEffect } from "react";
import { editChannel, fetchChannels, redorderChannel, type EditChannelPayload, type ReorderChannelPayload } from "~/api/channels";
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

export async function clientAction({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const intent = formData.get("intent");

  switch (intent) {
    case "edit": {
      const payload: EditChannelPayload = {
        id: parseInt(formData.get("id") as string),
        name: formData.get("name")?.toString() ?? null,
        description: formData.get("description")?.toString() ?? null,
      };

      return await editChannel(payload);
    }
    case "reorder": {
      const payload: ReorderChannelPayload = {
        id: parseInt(formData.get("id") as string),
        before_id: parseInt(formData.get("beforeId") as string)
      }
      console.log("SENDING REORDER PAYLOAD", payload);
      return await redorderChannel(payload);
    }

    case "delete": {
      const id = parseInt(formData.get("id") as string)
      const purge = parseInt(formData.get("purge") as string)

      if (id < 0) {
        throw new Response("Channel ID cannot be negative.");
      }

      return await fetch(`https://chat.localhost/channels?id=${id}&purge=${purge}`, {method: "DELETE"});
    }

    default: {
      console.log("Action not recognized");
      return;
    }
  }
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
