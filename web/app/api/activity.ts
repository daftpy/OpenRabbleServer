import type { ChannelsActivityResult, SessionActivityResult } from "~/types/api/activity";

export async function fetchSessionsActivity() {
  const response = await fetch("https://chat.localhost/activity/sessions")

  if (!response.ok) {
    throw new Response("Failed to load session activity", { status: response.status });
  }

  const data : SessionActivityResult = await response.json();

  return data.payload;
}

export async function fetchChannelsActivity() {
  const response = await fetch("https://chat.localhost/activity/channels");

  if (!response.ok) {
    throw new Response("Failed to load channel activity", { status: response.status });
  }

  const data : ChannelsActivityResult = await response.json();

  return data.payload.channels;
}
