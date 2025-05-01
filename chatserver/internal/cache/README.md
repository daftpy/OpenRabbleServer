# Chat Message Caching System

This project implements a hybrid caching mechanism for a chat application, combining in-memory caching with persistent storage to ensure efficient message handling and durability. It leverages Valkey (a Redis-compatible in-memory data store) for rapid access and PostgreSQL for long-term storage.

---

## Architecture

### Components

- **Valkey (In-Memory Cache):**
  - `recent_messages`: A circular buffer maintaining the most recent messages for quick retrieval.
  - `flush_messages`: A queue storing messages awaiting persistence to the database.
  - `cache_message_id`: A counter generating unique identifiers for each message.

- **PostgreSQL (Persistent Storage):**
  - `chat_messages` table: Stores all chat messages with associated metadata for long-term access.

---

## 🔁 Workflow

1. **Message Ingestion:**
   - Incoming messages are serialized to JSON and enriched with a unique `cache_id`.
   - A Lua script atomically:
     - Increments the `cache_message_id`.
     - Appends the message to both `recent_messages` and `flush_messages`.
     - Trims `recent_messages` to maintain a fixed size.

2. **Rate Limiting:**
   - Each user has an associated rate limit key (`ratelimit:<userID>`).
   - A Lua script increments the user's message count and sets an expiration based on the defined window.
   - If the user exceeds the allowed message count within the window, further messages are rejected.

3. **Flushing to Database:**
   - When `flush_messages` reaches a predefined size or at regular intervals, messages are batch-inserted into the `chat_messages` table.
   - After successful insertion, `flush_messages` is cleared to prevent duplicate entries.

---

## ⚙️ Configuration

- **`maxCacheSize`**: Maximum number of messages stored in `recent_messages` (default: 500).
- **`flushInterval`**: Interval for periodic flushing of messages to the database (default: 2 minutes).
- **`MessageLimit`**: Maximum number of messages a user can send within the rate limit window.
- **`WindowSeconds`**: Duration of the rate limit window in seconds.

## 📝 TODO

- [ ] **Scope cache keys by channel or scene**  
  Currently, all messages are stored under global keys (`recent_messages`, `flush_messages`, etc.). To support multi-channel or multi-room functionality, prepend keys with the channel ID or scene ID to isolate caches per chat context.

- [ ] **Improve error recovery on DB flush**  
  Right now, a single failed `INSERT` can abort the entire flush. Consider:
  - Logging and skipping bad messages (`continue` instead of `return`)
  - Adding a fallback queue for failed messages

- [ ] **Unit and integration test coverage**  
  Write tests for:
  - Lua script execution via Valkey client mocks
  - Rate limit enforcement
  - Cache serialization and flush logic

- [ ] **Graceful shutdown and flush**  
  Ensure that on server shutdown, any remaining `flush_messages` are flushed to the DB before exit.
