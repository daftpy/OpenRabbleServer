import RouteProtector from "~/components/route_protector";
import type { Route } from "./+types/index";
import { UsersPage } from "~/pages/users";
import { fetchUsers } from "~/api/users";

export type User = {
  id: string;
  username: string;
  is_banned: boolean;
}

export async function loader({ params }: Route.LoaderArgs) {
  return await fetchUsers();
}

export default function UsersRoute({loaderData,} : Route.ComponentProps) {
const users = loaderData;

  return (
    <RouteProtector>
      <UsersPage users={users} />
    </RouteProtector>
  )
}