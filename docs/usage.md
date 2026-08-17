# Advanced usage

## Hotkeys

| Key    | Action            |
|--------|-------------------|
| Enter  | Open log          |
| Esc    | Back              |
| F      | Filter            |
| R      | Reverse           |
| Ctrl+C | Exit              |
| F10    | Exit              |
| ↑↓ / jk| Line Up / Down    |
| PgUp   | Page Up           |
| PgDown | Page Down         |
| Home   | Navigate to Start |
| End / G| Navigate to End   |
| ?      | Show/Hide help    |

> Attempting to navigate past the last line in the log will put you in follow mode.

## Filtering

The filter matches a case-insensitive substring by default. Wrap the query in
slashes to match a regular expression instead:

```text
/timeout|deadline exceeded/
```

Both forms work for a full-text filter and for a filter by field, but they
match different text:

- A full-text filter matches the raw JSON line, so it sees field names, the
  quoting and the JSON escaping rather than what the table shows. `/^error/`
  never matches because every line starts with `{`, and a quote has to be
  written as `\"`. Fields are matched in the order they appear in the log, so
  a pattern spanning two of them, like `/error.*timeout/`, only matches the
  lines that happen to store them in that order.
- A filter by field matches the value as it is rendered in the column, after
  the level mapping and the time formatting have been applied. A pattern can
  only span one field, so anchors like `/^error$/` are useful here.

A regular expression is case-insensitive like the default filter. Prefix it
with `(?-i)` to make it case-sensitive:

```text
/(?-i)ERROR/
```

A term is only treated as a regular expression when it both starts and ends
with a slash, so `/` and `//` stay plain substring searches. To search for a
literal term that is wrapped in slashes, escape it with `\Q...\E`:

```text
/\Q/api/v1/\E/
```

## Configuration

```shell
jlv -config example.jlv.jsonc assets/example.log
jlv -config example.jlv.jsonc < assets/example.log
```

## Pull logs by URL

```shell
URL="https://raw.githubusercontent.com/hedhyw/json-log-viewer/main/assets/example.log"
curl "$URL" | jlv
```

## Preview logs from string

```shell
jlv << EOF
{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "day 1"}
{"time":"1970-01-02T00:00:00.00","level":"INFO","message": "day 2"}
EOF
```

## Show kubernetes logs

```shell
kubectl logs pod/POD_NAME -f | jlv
```

## View docker logs

```shell
docker logs -f 000000000000 2>&1 | jlv
```
