export async function fetchRateLimits() {
  const response = await fetch("https://chat.localhost/ratelimits")

  if (!response.ok) {
    throw new Response("Failed to load rate limits", { status: response.status });
  }

  const data = await response.json();

  return data;
}