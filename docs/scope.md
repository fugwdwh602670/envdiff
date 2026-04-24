# `envdiff scope` — Filter Keys by Prefix Scope

The `scope` command lets you narrow down an `.env` file to only the keys that
belong to a specific subsystem or concern, identified by a common prefix.

## Usage

```
envdiff scope <file> [prefix...] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `<file>` | Path to the `.env` file to inspect |
| `[prefix...]` | One or more key prefixes to match (e.g. `DB_`, `REDIS_`) |

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--invert` | `false` | Show keys that do **not** match any of the given prefixes |

## Examples

### Show only database-related keys

```bash
envdiff scope .env DB_
```

### Show keys for multiple subsystems

```bash
envdiff scope .env DB_ REDIS_ CACHE_
```

### Show everything except secrets

```bash
envdiff scope .env SECRET_ --invert
```

## Output

```
Scope (include prefixes: DB_)
  Total: 12 | Included: 3 | Excluded: 9

  DB_HOST=localhost
  DB_NAME=mydb
  DB_PORT=5432
```

When no prefixes are supplied all keys are shown with a summary line indicating
that no filtering was applied.

## Notes

- Prefix matching is case-sensitive; use uppercase prefixes to match standard
  env key conventions.
- Combine `scope` with shell redirection to produce a filtered env file:
  ```bash
  envdiff scope .env.production DB_ > .env.db
  ```
