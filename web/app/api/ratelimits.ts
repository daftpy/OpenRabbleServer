const hostname = import.meta.env.VITE_HOSTNAME;

export async function fetchRateLimits() {
  const response = await fetch(`https://chat.${hostname}/ratelimits`)

  if (!response.ok) {
    throw new Response("Failed to load rate limits", { status: response.status });
  }

  const data = await response.json();

  return data;
}

export async function updateRateLimits({
  id,
  messageLimit,
  windowSeconds,
}: {
  id: number;
  messageLimit: number;
  windowSeconds: number;
}) {
  const response = await fetch(`https://chat.${hostname}/ratelimits`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      id,
      message_limit: messageLimit,
      window_seconds: windowSeconds,
    }),
  });

  if (!response.ok) {
    throw new Response("Failed to update rate limits", { status: response.status });
  }

  const data = await response.json();
  return data;
}
