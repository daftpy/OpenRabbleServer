import RouteProtector from "~/components/route_protector";
import { BansPage } from "~/pages/users/bans";
import type { Route } from "./+types/bans";
import { fetchBans } from "~/api/users";

export async function loader({ params }: Route.LoaderArgs) {
  return await fetchBans();
}

export default function BansRoute({loaderData,} : Route.ComponentProps) {
  const { records } = loaderData;
  return (
    <RouteProtector>
      <BansPage records={records} />
    </RouteProtector>
  )
}