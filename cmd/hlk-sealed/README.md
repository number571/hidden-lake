# HLK-sealed

> Hidden Lake Kernel without secrets on disk

`hlk-sealed` runs an HLK node and a TCP connection to relays in one process for a desktop client. All secrets (private key, friends, tokens) are received from stdin and are kept in memory only. The only file is an anti-replay journal whose keys are `HMAC-SHA256(seen_key, hash)`.

## Contract (v3)

The client starts `hlk-sealed` without arguments, writes **one line** of JSON + `\n` to stdin and keeps stdin open while the process lives:

```json
{
  "priv_key": "PrivKey{...} or empty (generate a new key)",
  "network_key": "network from build (hybrid, not the default network) or empty",
  "friends": { "alias": "PubKey{...}" },
  "endpoints": ["relay-host:port"],
  "service_addr": "127.0.0.1:port (receives POST /push)",
  "api_token": "64 hex",
  "push_token": "64 hex",
  "seen_db_path": "path to the anti-replay journal file (directory must exist)",
  "seen_key": "64 hex"
}
```

stdout: `PRIVKEY <k>` (only if `priv_key` is empty), `PUBKEY <k>`, `READY <host:port>` (internal API, loopback). Nothing is written to stdout after `READY`; logs go to stderr.

Internal API (`Hlk-Auth: <api_token>` and `Host: <READY address>` are required, otherwise 401; body limit 64KiB):

- `GET /api/config/settings`
- `POST/DELETE /api/config/friends` — `{"alias_name","friend_key"}` / `{"alias_name"}`
- `GET/POST/DELETE /api/config/connects` — body `host:port`
- `GET /api/network/online`
- `POST /api/network/request?friend=<alias>` — HL request `POST hls-messenger /push`, sent through the queue without waiting for a response

Incoming requests of friends (`POST /push` of `hls-messenger` only) are forwarded to `http://service_addr/push` with headers `Hlk-Sender-Name: <alias>` and `Hlk-Auth: <push_token>` only; headers of the peer are never copied.

EOF on stdin stops the process (exit within 1 second).

The parent must read stderr continuously or redirect it to null: logs are written there during the whole life of the process, and a full pipe would block the node.

Exit codes (nothing is written to stdout before an error):

- `0` — stopped (EOF on stdin or SIGINT/SIGTERM);
- `1` — runtime error: the seen db is locked by another process, damaged or its directory does not exist; the node, the TCP adapter or the internal API failed; writing to stdout failed;
- `2` — invalid boot line (not a single line, invalid JSON, unknown field or invalid value).

## Build

```bash
$ go build -trimpath -ldflags="-s -w -buildid=" ./cmd/hlk-sealed
```

`-trimpath` removes absolute paths of the build machine (with the user name) from the binary, `-s -w` removes symbol and debug tables, `-buildid=` makes the build reproducible.
