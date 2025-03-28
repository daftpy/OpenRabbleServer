// Represents a user on the chat server
export type User = {
  id: string;
  username: string;
  is_banned: boolean;
}

/*
  Represents a ban record on the site. owner_id is
  the user that initiated the ban. banished_id is
  the user id of the user that is banned. Basic
  details like start, end, duration, and pardon 
  status are available.
*/
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
