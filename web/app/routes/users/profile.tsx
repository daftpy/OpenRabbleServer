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
  const messagesRes = await fetch(`https://chat.localhost/messages?user_id=${serverData.id}`);
  const messagesData = await messagesRes.json();

  const activityRes = await fetch(`https://chat.localhost/activity/sessions?user_id=${serverData.id}`)
  const activityData = await activityRes.json();
  console.log("ACTIVITY: ", activityData);

  return { ...serverData, ...messagesData.payload, ...activityData.payload };
}

// force the client loader to run during hydration
clientLoader.hydrate = true as const; // `as const` for type inference

export function HydrateFallback() {
  return <div>Loading...</div>;
}

export default function UserRoute({loaderData,} : Route.ComponentProps) {
  const { username, id, messages, session_activity } = loaderData;
  console.log("USERNAME", username);
  return (
    <RouteProtector>
      <UserPage username={username} id={id} messages={messages} session_activity={session_activity}/>
    </RouteProtector>
  )
}