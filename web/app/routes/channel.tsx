import { fetchChannelsActivity } from "~/api/activity";
import type { Route } from "./+types/channel";
import { useEffect } from "react";
import { addChannel, deleteChannel, editChannel, fetchChannels, redorderChannel } from "~/api/channels";
import RouteProtector from "~/components/route_protector";
import { ChannelPage } from "~/pages/channels";
import type { AddChannelPayload, EditChannelPayload, ReorderChannelPayload } from "~/types/api/channel";

export async function loader({ params }: Route.LoaderArgs) {
  return fetchChannels();
}

export async function clientLoader({
  serverLoader,
  params,
}: Route.ClientLoaderArgs) {
  
  const serverData = await serverLoader();

  const activityData = await fetchChannelsActivity();

  return { channels: serverData, channel_activity: activityData };
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

  const id = parseInt(formData.get("id") as string);
  const name = formData.get("name")?.toString().trim();
  const description = formData.get("description")?.toString().trim();
  const beforeId = parseInt(formData.get("beforeId") as string);
  const purge = parseInt(formData.get("purge") as string);

  switch (intent) {
    case "edit": {
      const payload: EditChannelPayload = {
        id,
        name: name || null,
        description: description || null,
      };
      return await editChannel(payload);
    }

    case "reorder": {
      const payload: ReorderChannelPayload = {
        id,
        before_id: beforeId,
      }
      return await redorderChannel(payload);
    }

    case "delete": {
      if (id < 0 || ![0, 1].includes(purge)) {
        throw new Response("Invalid request for deletion", { status: 400 });
      }
      return await deleteChannel(id, purge);
    }

    case "add": {
      if (!name) {
        throw new Response("Channel name is required", { status: 400 });
      }

      const payload: AddChannelPayload = {
        name,
        description: description || null,
      }
      return await addChannel(payload);
    }

    default:
      console.warn("Unrecognized action:", intent);
      return new Response("Invalid intent", { status: 400 });
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
