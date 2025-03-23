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