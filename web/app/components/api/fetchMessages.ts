export async function fetchMessagesFromAPI(params: { keyword?: string; channels?: string[]; user_id?: string; limit?: string; offset?: string }) {
    const { keyword = "", channels = [], user_id, limit = "2", offset = "0" } = params;
  
    const queryParams = new URLSearchParams();
    if (keyword) queryParams.append("keyword", keyword);
    if (user_id) queryParams.append("user_id", user_id);
    channels.forEach(channel => queryParams.append("channel", channel));
    queryParams.append("limit", limit);
    queryParams.append("offset", offset);
  
    console.log("Fetching messages with params:", queryParams.toString());
  
    try {
      const response = await fetch(`https://chat.localhost/messages?${queryParams.toString()}`);
      if (!response.ok) throw new Response("Failed to fetch messages", { status: response.status });
  
      const messageData = await response.json();
      if (!messageData.payload || !Array.isArray(messageData.payload.messages)) {
        console.error("Unexpected response format:", messageData);
        throw new Response("Invalid response format", { status: 500 });
      }
      return { messages: messageData.payload.messages ?? [], hasMore: messageData.payload.has_more };
    } catch (error) {
      throw new Response("Error fetching messages", { status: 500 });
    }
  }
  