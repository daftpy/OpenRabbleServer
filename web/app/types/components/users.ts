export type User = {
  id: string;
  username: string;
  is_banned: boolean;
}

export type BanRecord = {
  id: string;
  owner_id: string;
  banished_id: string;
  banished_username: string;
  start: string;    // ISO timestamp
  end: string;      // ISO timestamp
  duration: string; // "HH:MM:SS.ssssss"
  pardoned: boolean;
};
