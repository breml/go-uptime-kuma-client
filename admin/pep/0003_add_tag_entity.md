# PEP 0003: Add Tag Entity

Extend the client with CRUD operations for the tag entity.

Add the file `tag.go` to contain `GetTags`, `GetTag`, `CreateTag`, `UpdateTag`, and `DeleteTag` methods.
Extend @client.go where necessary.
Add tests to `tag_test.go`.
Add the `tag` package in @tag containing the type definitions and implementations. Since tags are simple and don't have subtypes, only `tag.go` is needed, no `tag_base.go`.

Consider @.scratch/uptime-kuma/server/server.js and @.scratch/uptime-kuma/src/components/TagEditDialog.vue for reference of the socket.io events and payloads.
