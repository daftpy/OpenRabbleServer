import type { BanRecord, User } from "../components/users";

// The required fields to ban a user
export type BanUserPayload = {
  banished_id: string;
  reason: string | null;
  duration: number | null;
}

// Response after banning a user
export type BanUserResponse = {
  message: string;
}

// Response after fetching a user
export type FetchUsersReponse = {
  type: string;
  sender: string;
  payload: {
    users: User[];
  };
};

// Response after banning a user
export type FetchBansResponse = {
  type: "ban_records_result";
  sender: string;
  payload: {
    records: BanRecord[];
    has_more: boolean;
  };
};
