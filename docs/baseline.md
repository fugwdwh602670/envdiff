# Baseline Feature

The `baseline` command allows you to snapshot the current diff between two `.env` files and later compare against it to detect regressions or track progress.

## Commands

### `envdiff baseline save <fileA> <fileB>`

Runs a diff between `fileA` and `fileB` and saves the results as a baseline JSON file.

```sh
envdiff baseline save .env .env.prod --baseline .envdiff-baseline.json
```

### `envdiff baseline check <fileA> <fileB>`

Runs the current diff and compares it against the saved baseline, reporting:
- **New issues**: keys that are now problematic but weren't in the baseline
- **Resolved**: keys that were problematic in the baseline but are now clean

Exits with code `1` if there are any new issues.

```sh
envdiff baseline check .env .env.prod --baseline .envdiff-baseline.json
```

## Flags

| Flag         | Default                      | Description              |
|--------------|------------------------------|--------------------------|
| `--baseline` | `.envdiff-baseline.json`     | Path to baseline file    |

## Baseline File Format

The baseline is stored as JSON:

```json
{
  "created_at": "2024-01-01T00:00:00Z",
  "file_a": ".env",
  "file_b": ".env.prod",
  "results": [
    {"key": "DB_HOST", "status": "mismatch", "value_a": "localhost", "value_b": "prod-db"}
  ]
}
```
