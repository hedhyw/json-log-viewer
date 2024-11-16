# Customization

The application will look for the config `.jlv.jsonc` in the working directory or in the home directory:
- `$PWD/.jlv.jsonc`;
- `$HOME/.jlv.jsonc`.

It's also possible to define the path to the configuration using "-config" flag.

The Json path supports the described in [yalp/jsonpath](https://github.com/yalp/jsonpath#jsonpath-quick-intro) syntax.

Example configuration: [example.jlv.jsonc](../example.jlv.jsonc).

## Time Formats
JSON Log Viewer can handle a variety of datetime formats when parsing your logs.
The value is formatted by default in the "[RFC3339](https://www.rfc-editor.org/rfc/rfc3339)" format. The format is configurable, see the `time_format` field in the [config](../example.jlv.jsonc).

### `time`
This will return the exact value that was set in the JSON document.

### `numerictime`
This is a "smart" parser. It can accept an integer, a float, or a string. If it is numeric (`1234443`, `1234443.589`, `"1234443"`, `"1234443.589"`), based on the number of digits, it will parse as seconds, milliseconds, or microseconds. The output is a UTC-based RFC 3339 datetime.

If a string such as `"2023-05-01T12:00:34Z"` or `"---"` is used, the value will just be carried forward to your column.  

If you find that the smart parsing is giving unwanted results or you need greater control over how a datetime is parsed, considered using one of the other time formats instead.

### `secondtime`
This will attempt to parse the value as number of seconds and render as a UTC-based RFC 3339. Values accepted are integer, string, or float.

### `millitime`
Similar to `secondtime`, this will attempt to parse the value as number of milliseconds. Values accepted are integer, string, or float.

### `microtime`
Similar to `secondtime` and `millistime`, this will attempt to parse the value as number of microseconds. Values accepted are integer, string, or float.
