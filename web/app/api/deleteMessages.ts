// Deletes messages from the API in bulk, given an array of IDs.
export async function deleteMessagesFromAPI(ids: number[]): Promise<Response> {
  if (!Array.isArray(ids) || ids.length === 0) {
    // Return a custom error response or throw — whichever pattern you prefer
    throw new Error("No valid IDs provided");
  }

  const body = JSON.stringify({ ids });
  const response = await fetch("https://chat.localhost/messages", {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
    body,
  });

  // Return the raw response so the caller can check response.ok and status code.
  return response;
}