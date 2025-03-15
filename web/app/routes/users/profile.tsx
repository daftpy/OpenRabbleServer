import RouteProtector from "~/components/route_protector";
import { UserPage } from "~/pages/users/profile";
import type { Route } from "./+types/profile";
import { useLoaderData } from "react-router";

// Get the users information by username. Here we need the userID
export async function loader({ params }: Route.LoaderArgs) {
  const response = await fetch(`https://chat.localhost/users?username=${params.userId}`);
  if (!response.ok) {
    throw new Response("Failed to load users", { status: response.status });
  }

  const data = await response.json();
  if (data.payload.users.length == 0) {
    return [];
  }
  return data.payload.users[0];
}

// Use the userID to retrieve the users meessages
export async function clientLoader({
  serverLoader,
  params,
}: Route.ClientLoaderArgs) {
  
  const serverData = await serverLoader();
  const messagesRes = await fetch(`https://chat.localhost/messages?user_id=${serverData.id}&limit=10&offset=0`);
  const messagesData = await messagesRes.json();

  const activityRes = await fetch(`https://chat.localhost/activity/sessions?user_id=${serverData.id}`)
  const activityData = await activityRes.json();
  console.log("ACTIVITY: ", activityData);
  console.log("MESSAGES:", messagesData);

  return { ...serverData, ...messagesData.payload, ...activityData.payload };
}

// force the client loader to run during hydration
clientLoader.hydrate = true as const; // `as const` for type inference

export function HydrateFallback() {
  return <div>Loading...</div>;
}

export async function clientAction({ params, request }: Route.ActionArgs) {
  if (request.method !== "POST") {
    return new Response("Method Not Allowed", { status: 405 });
  }

  const formData = await request.formData();
  const reason = formData.get("reason");
  const duration = formData.get("duration");

  let data = {
    banished_id: formData.get("banishedId"),
    reason: reason,
    duration: 0
  }

  if (duration) {
    data.duration = parseInt(duration.toString(), 10);
  }
  
  const banResponse = await fetch("https://chat.localhost/users/ban", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

  if (!banResponse.ok) {
    return new Response("Failed to ban user", { status: 500 });
  }
  console.log("User banned");

  return new Response(
    JSON.stringify({ message: "User banned successfully" }),
    { status: 200 }
  );
}

export default function UserRoute({loaderData,} : Route.ComponentProps) {
  const { username, id, messages, has_more, session_activity } = loaderData;
  console.log("USERNAME", username);
  return (
    <RouteProtector>
      <UserPage username={username} id={id} messages={messages} session_activity={session_activity} hasMore={has_more} />
    </RouteProtector>
  )
}