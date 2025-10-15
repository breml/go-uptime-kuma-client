# PEP 0005: Add Monitor Push Type

Extend the `monitor` package in @monitor to support the monitor type "push".

Follow the existing patterns for other monitor types (e.g. HTTP), in particular:

- Embed the `Base` struct.
- Add tests to the client it self (@monitor_test.go) as well as to the `monitor` package in @monitor/monitor_push_test.go.
- Make sure, the generated code does build without errors, does not have any linting issues and that all tests (unit and integration) pass.

Consider @.scratch/uptime-kuma/server/model/monitor.js and @.scratch/uptime-kuma/src/pages/EditMonitor.vue for reference of the implementation details, in particular the handling of different monitor types and their specific properties.
