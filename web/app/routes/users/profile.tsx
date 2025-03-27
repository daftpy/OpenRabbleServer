import RouteProtector from "~/components/route_protector";
import { UserPage } from "~/pages/users/profile";
import type { Route } from "./+types/profile";
import { banUser, fetchUser } from "~/api/users";

// Get the users information by username.
export async function loader({ params }: Route.LoaderArgs) {
  return await fetchUser(params.userId); // userId is actually username
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
  // Get the formdata
  const formData = await request.formData();

  // Check for banishedId
  const banishedId = formData.get("banishedId")?.toString().trim();
  if (!banishedId) {
    throw new Response("Missing required field: banishedId", { status: 400 });
  }

  // Get the reason or null
  const reason = formData.get("reason")?.toString().trim() || null;

  // Get the duration
  const duration = formData.get("duration")?.toString().trim() || null;
  const parsedDuration = duration ? parseInt(duration, 10) : null;

  return await banUser({ banished_id: banishedId, reason, duration: parsedDuration });
}

export default function UserRoute({loaderData,} : Route.ComponentProps) {
  const { username, id, is_banned, messages, has_more, session_activity } = loaderData;
  console.log("USERNAME", username);
  return (
    <RouteProtector>
      <UserPage username={username} id={id} messages={messages} session_activity={session_activity} isBanned={is_banned} hasMore={has_more} />
    </RouteProtector>
  )
}