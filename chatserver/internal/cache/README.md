# Chat Message Caching System

This project implements a hybrid caching mechanism for a chat application, combining in-memory caching with persistent storage to ensure efficient message handling and durability. It leverages **Valkey** (a Redis-compatible in-memory data store) for rapid access and **PostgreSQL** for long-term storage.


## Architecture

### Valkey (In-Memory Cache)

- `recent_messages`: Circular buffer of recent public messages.
- `flush_messages`: Queue of public messages waiting to be persisted.
- `recent_private_messages:<userID>`: Per-user circular cache of private messages (both sent and received).
- `flush_private_messages`: Global queue of private messages pending database flush.
- `cache_message_id`: Auto-increment counter for public messages.
- `cache_private_message_id`: Auto-increment counter for private messages.
- `ratelimit:<userID>`: Tracks per-user message counts for rate limiting.


### PostgreSQL (Persistent Storage)

- `chat_messages`: Stores all flushed public chat messages.
- `private_messages`: Stores all flushed private messages.


## 🔁 Workflow

### 1. Message Ingestion

- Public and private messages are serialized as JSON and enriched with a unique `cache_id`.
- A Lua script:
  - Atomically increments the appropriate counter (`cache_message_id` or `cache_private_message_id`).
  - Stores messages in both:
    - Circular cache for recent access.
    - Flush queue for eventual DB persistence.
  - If the sender is also the recipient (a DM to self), the message is only stored once.


### 2. 🚦 Rate Limiting

- Each user is tracked via a key: `ratelimit:<userID>`.
- A Lua script:
  - Uses `INCR` to count messages.
  - Applies `EXPIRE` to set the time window.
  - Rejects messages once the limit is exceeded.
- Limits are customizable:
  - `MessageLimit`: max messages per window.
  - `WindowSeconds`: length of the rate window in seconds.


### 3. Flushing to Database

- Flush is triggered when:
  - The flush queue (`flush_messages` or `flush_private_messages`) reaches the threshold (`maxCacheSize`).
  - Or periodically using a timer (`flushInterval`).
- Messages are written to the database in a single transaction.
- Upon success, the flush queue is cleared to prevent duplicates.


## Configuration

| Setting           | Purpose                                         | Default        |
|------------------|--------------------------------------------------|----------------|
| `maxCacheSize`   | Max number of messages to keep before flush     | `500`          |
| `flushInterval`  | Interval to flush messages automatically        | `2 minutes`    |
| `MessageLimit`   | Max number of messages allowed per user         | configurable   |
| `WindowSeconds`  | Duration of message window for rate limiting    | configurable   |

You can update rate limits dynamically using:

```go
cache.UpdateRateLimitSettings(limit, window)
```


## Features

- Atomic caching and trimming via Lua.
- Full support for both public and private messages.
- Automatic and manual database flush control.
- Per-user rate limiting with Lua-based enforcement.
- Self-DMs are deduplicated to avoid storing duplicates.


## 📝 TODO

- [ ] **Improve error recovery on DB flush**  
  A single bad message aborts the entire transaction. Consider:
  - Logging and skipping failed inserts with `continue`
  - Fallback queue for unflushable messages

- [ ] **Add testing coverage**  
  Write unit and integration tests for:
  - Lua execution
  - Flush logic and failure paths
  - Rate limiting behavior
  - Serialization integrity

- [ ] **Graceful shutdown**  
  Ensure that any pending flush queues are committed to the DB on application shutdown.


## 🧪 Dev Notes

You can inspect and delete private message keys in Valkey using:

```bash
valkey-cli
127.0.0.1:6379> keys recent_private_messages:*
127.0.0.1:6379> del recent_private_messages:<userID>
```

To clear **all private message caches**:

```bash
valkey-cli --scan --pattern 'recent_private_messages:*' | xargs valkey-cli del
```

