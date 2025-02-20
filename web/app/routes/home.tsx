import type { Route } from "./+types/home";
import { Dash } from "../dash/dash"
import { useLoaderData } from "react-router";

export async function loader({ params }: Route.LoaderArgs) {
  const response = await fetch("https://chat.localhost/channels"); // ✅ Update to use the correct service name
  
  if (!response.ok) {
    throw new Response("Failed to load channels", { status: response.status });
  }

  const data = await response.json();
  console.log(data);
  return data.channels;
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Welcome to React Router!" },
  ];
}

export default function Home({loaderData,}: Route.ComponentProps) {
  const channels = useLoaderData();

  return <Dash channels={channels} />;
}
