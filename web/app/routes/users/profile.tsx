import RouteProtector from "~/components/route_protector";
import { UserPage } from "~/pages/users/profile";
import type { Route } from "./+types/profile";
import { useLoaderData } from "react-router";

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

export async function clientLoader({
  serverLoader,
  params,
}: Route.ClientLoaderArgs) {
  
  const serverData = await serverLoader();
  const res = await fetch(`https://chat.localhost/messages?user_id=${serverData.id}`);
  const data = await res.json();
  console.log("CLIENT ACTION RES", data.payload);
  return { ...serverData, ...data.payload };
}

// force the client loader to run during hydration
clientLoader.hydrate = true as const; // `as const` for type inference

export function HydrateFallback() {
  return <div>Loading...</div>;
}

export default function UserRoute() {
  const { username, id, messages } = useLoaderData();
  console.log("USERNAME", username);
  return (
    <RouteProtector>
      <UserPage username={username} id={id} messages={messages} />
    </RouteProtector>
  )
}