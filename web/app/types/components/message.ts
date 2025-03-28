/*
  Represents a message on the chat server. Currently, 
  this represents messages that were pushed to the 
  database. Messages that exist in the cache do not 
  posess an id.
*/
export type Message = {
  id: number;
  owner_id: string;
  username: string;
  message: string;
  channel: string;
  authored_at: string;
}
