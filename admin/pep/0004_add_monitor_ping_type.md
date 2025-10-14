# PEP 0004: Add Monitor Ping Type

Extend the `monitor` package in @monitor to support the monitor type "ping".

Follow the existing patterns for other monitor types (e.g. HTTP), in particular:

- Embed the `Base` struct.
- Add tests to the client it self (@monitor_test.go) as well as to the `monitor` package in @monitor/monitor_ping_test.go.

Consider @.scratch/uptime-kuma/server/model/monitor.js and @.scratch/uptime-kuma/src/pages/EditMonitor.vue for reference of the implementation details, in particular the handling of different monitor types and their specific properties.
