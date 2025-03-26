// const hostname = import.meta.env.VITE_HOSTNAME;

export async function fetchChannels() {
  const response = await fetch("https://chat.localhost/channels");
  if (!response.ok) {
    throw new Response("Failed too load channels", { status: response.status });
  }

  const data = await response.json();
  if (data.channels == null) {
    return [];
  }

  return data.channels;
}

export type EditChannelPayload = {
  id: number;
  name: string | null;
  description: string | null;
}

export type ReorderChannelPayload = {
  id: number,
  before_id: number
}

export async function editChannel(payload : EditChannelPayload) {
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

  return await response.json();
}

export const redorderChannel = async (payload: ReorderChannelPayload) => {
  const response = await fetch("https://chat.localhost/channels", {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload)
  });
  console.log("CHANNEL REORDER PAYLOAD", payload)

  if (!response.ok) {
    console.log(response);
    throw new Error("Failed to reorder the channel");
  }

  return await response.json();
}
