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

export default function UserRoute() {
  const { username, id } = useLoaderData();
  console.log("USERNAME", username);
  return (
    <RouteProtector>
      <UserPage username={username} id={id} />
    </RouteProtector>
  )
}