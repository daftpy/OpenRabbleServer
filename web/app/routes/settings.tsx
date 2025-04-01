import { SettingsPage } from "~/pages/settings";
import type { Route } from "./+types/settings";
import { fetchRateLimits } from "~/api/ratelimits";

export async function loader({ params }: Route.LoaderArgs) {
  const data = await fetchRateLimits();

  return data;
}

// in routes/settings.tsx
import { updateRateLimits } from "~/api/ratelimits";

export async function clientAction({ request }: Route.ActionArgs) {
  const formData = await request.formData();

  const id = 1;
  const wordLimit = formData.get("wordLimit")?.toString().trim();
  const windowSeconds = formData.get("windowSeconds")?.toString().trim();

  if (!wordLimit || !windowSeconds) {
    throw new Error("Both word limit and window seconds must be provided");
  }

  const messageLimit = parseInt(wordLimit);
  const windowSecs = parseInt(windowSeconds);

  if (isNaN(messageLimit) || isNaN(windowSecs)) {
    throw new Error("Rate limit values must be numeric");
  }

  const result = await updateRateLimits({
    id,
    messageLimit,
    windowSeconds: windowSecs,
  });
  console.log("RESULT", result);
  return result;
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