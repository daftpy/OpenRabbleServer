// const hostname = import.meta.env.VITE_HOSTNAME;

import type { AddChannelPayload, AddChannelResponse, ChannelResponse, EditChannelResponse, FetchChannelsResponse, ReorderChannelResponse } from "~/types/api/channel";
import type { EditChannelPayload, ReorderChannelPayload } from "~/types/api/channel";

export async function fetchChannels(): Promise<ChannelResponse[]> {
  const response = await fetch("https://chat.localhost/channels");
  if (!response.ok) {
    throw new Response("Failed too load channels", { status: response.status });
  }

  const data: FetchChannelsResponse = await response.json();
  return data.channels ?? [];
}

export async function editChannel(payload : EditChannelPayload): Promise<string> {
  const response = await fetch("https://chat.localhost/channels", {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload)
  });

  if (!response.ok) {
    throw new Error("Failed to update channel");
  }

  const data : EditChannelResponse = await response.json()
  return data.message;
}

export const redorderChannel = async (payload: ReorderChannelPayload): Promise<string> => {
  const response = await fetch("https://chat.localhost/channels", {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload)
  });

  if (!response.ok) {
    console.log(response);
    throw new Error("Failed to reorder the channel");
  }
  const data: ReorderChannelResponse = await response.json();
  return data.message;
}

export const deleteChannel = async (id: number, purge: number) => {
  if (![0, 1].includes(purge)) {
    throw new Error("Invalid purge value - must be 0 or 1");
  }
  return await fetch(`https://chat.localhost/channels?id=${id}&purge=${purge}`, {method: "DELETE"});
}

export const addChannel = async (payload: AddChannelPayload): Promise<AddChannelResponse> => {
  const response = await fetch(`https://chat.localhost/channels`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    throw new Response("Failed to create channel", { status: response.status });
  }

  const data : AddChannelResponse = await response.json();
  return { message: data.message, name: data.name };
}
