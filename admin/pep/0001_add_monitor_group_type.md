# PEP 0001: Add Monitor Group Type

Extend the `monitor` package in @monitor to support the "special" monitor type "monitor group".

Follow the existing patterns for other monitor types (e.g. HTTP), in particular:

- Embed the `Base` struct.
- Add tests to the client it self (@monitor_test.go) as well as to the `monitor` package in @monitor/monitor_group_test.go.

Consider @.scratch/uptime-kuma/server/model/monitor.js and @.scratch/uptime-kuma/src/pages/EditMonitor.vue for reference of the implementation details, in particular the handling of different monitor types and their specific properties.
