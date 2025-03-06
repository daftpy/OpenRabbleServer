import RouteProtector from "~/components/route_protector";
import type { Route } from "./+types/index";
import { UsersPage } from "~/pages/users";
import { useLoaderData } from "react-router";

export async function loader({ params }: Route.LoaderArgs) {
  const response = await fetch("https://chat.localhost/users");
  if (!response.ok) {
    throw new Response("Failed to load users", { status: response.status });
  }

  const data = await response.json();
  console.log(data);
  if (data.payload == null) {
    return [];
  }
  return data.payload;
}

export default function UsersRoute({loaderData,} : Route.ComponentProps) {
  const { users } = useLoaderData() as { users: any[] };
  return (
    <RouteProtector>
      <UsersPage users={users} />
    </RouteProtector>
  )
}