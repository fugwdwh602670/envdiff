# Watch Mode

The `watch` command polls two `.env` files at a configurable interval and prints a diff whenever their contents change.

## Usage

```bash
envdiff watch <fileA> <fileB> [--interval <ms>]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval`, `-i` | `2000` | Poll interval in milliseconds |

## Example

```bash
envdiff watch .env.staging .env.production --interval 1000
```

Output when a change is detected:

```
Watching .env.staging and .env.production (interval 1000ms)...
--- diff update ---
[MISMATCH]       DB_HOST          staging-db       prod-db
[MISSING_IN_B]   NEW_FLAG         true             -
```

## Behavior

- On startup, an initial diff is emitted immediately.
- Subsequent diffs are only printed when the file contents actually change.
- Press `Ctrl+C` to stop watching.
- Files are compared using MD5 hashing to detect changes efficiently before re-parsing.
