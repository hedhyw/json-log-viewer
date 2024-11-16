# Features

It is an **interactive** tool for viewing and analyzing complex structured [json-log](../assets/example.log) files.

Main features:
1. It is interactive.
2. It displays a compact list of log entries.
3. It is possible to expand the log and see the full prettified JSON tree.
4. All non-json logs are captured.
5. Fields are [customizable](./customization.md).
6. Filtering is easy to use.
7. Log levels are colorized.
8. Transforming numeric timestamps.

It uses [antonmedv/fx](https://github.com/antonmedv/fx) for viewing JSON records and [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) for organizing the terminal UI. The tool is inspired by the project [json-log-viewer](https://github.com/gistia/json-log-viewer) which is unfortunately outdated and deserted.

The application is designed to help in visualization, navigation, and analyzing of JSON-formatted log data in a user-friendly and interactive manner. It provides a structured and organized view of the JSON logs, making it easier to comprehend the hierarchical nature of the data. It uses collapsible/expandable tree structures, indentation, and color-coded syntax to represent the JSON objects and arrays. It is possible to search for specific keywords, phrases, or patterns within the JSON logs. So it helps to significantly simplify the process of working with JSON logs, making it more intuitive and efficient. It is easy to troubleshoot issues, monitor system performance, or gain a deeper understanding of the application's behavior by analyzing its log data in post-mortem.