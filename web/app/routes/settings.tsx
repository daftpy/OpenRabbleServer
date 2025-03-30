import { SettingsPage } from "~/pages/settings";
import type { Route } from "./+types/settings";
import { fetchRateLimits } from "~/api/ratelimits";

export async function loader({ params }: Route.LoaderArgs) {
  const data = await fetchRateLimits();

  return data;
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Welcome to React Router!" },
  ];
}

export default function SettingsRoute({loaderData,} : Route.ComponentProps ) {
  const { rate_limiter } = loaderData;
  return (
    <SettingsPage rate_limiter={rate_limiter}  />
  )
}