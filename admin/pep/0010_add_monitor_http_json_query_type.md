# PEP 0010: Add Monitor HTTP JSON Query Type

Extend the `monitor` package in @monitor to support the monitor type "http_json_query".

Follow the existing patterns for other monitor types (e.g. HTTP), in particular:

- Embed the `Base` struct.
- Add tests to the client it self (@monitor_test.go) as well as to the `monitor` package in @monitor/monitor_http_json_query_test.go.
- Make sure, the generated code does build without errors, does not have any linting issues and that all tests (unit and integration) pass.

Consider @.scratch/uptime-kuma/server/model/monitor.js and @.scratch/uptime-kuma/src/pages/EditMonitor.vue for reference of the implementation details, in particular the handling of different monitor types and their specific properties.
