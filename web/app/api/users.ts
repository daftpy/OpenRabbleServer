import type { BanUserPayload, BanUserResponse, FetchBansResponse, FetchUsersReponse } from "~/types/api/users";
import type { BanRecord, User } from "~/types/components/users";

const hostname = import.meta.env.VITE_HOSTNAME;

export async function fetchUser(username: string) : Promise<User> {
  const response = await fetch(`https://chat.${hostname}/users?username=${username}`);
  if (!response.ok) {
    throw new Response("Failed to load users", { status: response.status });
  }

  const data : FetchUsersReponse = await response.json();
  const user = data.payload.users?.[0];

  if (!user) {
    throw new Response("User now found", { status: 404 });
  }

  return user;
}

export async function fetchUsers() : Promise<User[]> {
  const response = await fetch(`https://chat.${hostname}/users`);
  if (!response.ok) {
    throw new Response("Failed to load users", { status: response.status });
  }

  const data : FetchUsersReponse = await response.json();
  
  return data.payload.users ?? [];
}

export async function fetchBans() : Promise<FetchBansResponse["payload"]>  {
  const response = await fetch(`https://chat.${hostname}/users/bans`);
  if (!response.ok) {
    throw new Response("Failed to load bans", { status: response.status });
  }

  const data : FetchBansResponse = await response.json();

  return {
    records: data.payload?.records ?? [],
    has_more: data.payload?.has_more ?? false
  };
}

export async function banUser(payload: BanUserPayload) : Promise<BanUserResponse> {
  console.log("PAYLOAD", payload);
  const response = await fetch(`https://chat.${hostname}/users/ban`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });

  if (!response.ok) {
    throw new Response("Failed to ban user", { status: 500 });
  }

  return { message: "User banned succesfully" };
}

export async function pardonUser(banId: number): Promise<{ message: string }> {
  const response = await fetch(`https://chat.${hostname}/users/ban?ban_id=${banId}`, {
    method: "DELETE",
  });

  if (!response.ok) {
    throw new Response("Failed to pardon user", { status: response.status });
  }
  console.log(response);

  return { message: "User pardoned successfully" };
}
