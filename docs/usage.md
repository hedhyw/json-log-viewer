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
| ↑↓     | Line Up / Down    |
| Home   | Navigate to Start |
| End    | Navigate to End   |
| ?      | Show/Hide help    |

> Attempting to navigate past the last line in the log will put you in follow mode.

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
