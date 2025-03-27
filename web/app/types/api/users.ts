import type { BanRecord, User } from "../components/users";

export type BanUserPayload = {
  banished_id: string;
  reason: string | null;
  duration: number | null;
}

export type BanUserResponse = {
  message: string;
}

export type FetchUsersReponse = {
  type: string;
  sender: string;
  payload: {
    users: User[];
  };
};

export type FetchBansResponse = {
  type: "ban_records_result";
  sender: string;
  payload: {
    records: BanRecord[];
    has_more: boolean;
  };
};
