import RouteProtector from "~/components/route_protector";
import type { Route } from "./+types/index";
import { UsersPage } from "~/pages/users";

export type User = {
  id: string;
  username: string;
}

type UsersSearchResult = {
  type: string;
  sender: string;
  payload: {
    users: User[];
  };
};

export async function loader({ params }: Route.LoaderArgs) {
  const response = await fetch("https://chat.localhost/users");
  if (!response.ok) {
    throw new Response("Failed to load users", { status: response.status });
  }

  const data : UsersSearchResult = await response.json();
  console.log(data);
  if (data.payload == null) {
    return [];
  }
  return data;
}

export default function UsersRoute({loaderData,} : Route.ComponentProps) {
const { payload: { users } } = loaderData as UsersSearchResult;

  return (
    <RouteProtector>
      <UsersPage users={users} />
    </RouteProtector>
  )
}