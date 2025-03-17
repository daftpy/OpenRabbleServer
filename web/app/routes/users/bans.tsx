import RouteProtector from "~/components/route_protector";
import { BansPage } from "~/pages/users/bans";
import type { Route } from "./+types/bans";

export async function loader({ params }: Route.LoaderArgs) {
  const response = await fetch("https://chat.localhost/users/bans");
  if (!response.ok) {
    throw new Response("Failed to load users", { status: response.status });
  }

  const data : any = await response.json();
  console.log("BANS DATA", data);
  if (data.payload == null) {
    return [];
  }
  return data.payload;
}

export default function BansRoute({loaderData,} : Route.ComponentProps) {
  const {records} = loaderData;
  console.log("RECORDS:", records);
  return (
    <RouteProtector>
      <BansPage records={records} />
    </RouteProtector>
  )
}