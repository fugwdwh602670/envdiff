# Redaction

`envdiff` can automatically redact sensitive values from output so that secrets never appear in logs, CI output, or shared reports.

## Default Patterns

The following key substrings trigger redaction (case-insensitive):

- `secret`
- `password` / `passwd`
- `token`
- `apikey` / `api_key`
- `private`

Any key containing one of these substrings will have its values replaced with `***REDACTED***` in all output formats (text, JSON, CSV).

## CLI Flag

Pass `--redact` to enable redaction when running a diff:

```bash
envdiff --redact .env.staging .env.production
```

Redaction is **off by default** to preserve full diff visibility in local development.

## Config File

You can enable redaction permanently via the config file:

```bash
envdiff config set redact true
```

Or add custom patterns:

```bash
envdiff config set redact-patterns "secret,token,internal"
```

## Notes

- Redaction applies to output only; the original `.env` files are never modified.
- Both `ValueA` and `ValueB` are redacted when a key matches, including `missing_in_b` and `missing_in_a` statuses.
