# go-uptime-kuma-client

[![Test Status](https://github.com/breml/go-uptime-kuma-client/actions/workflows/ci.yml/badge.svg)](https://github.com/breml/go-uptime-kuma-client/actions)
[![Go Report Card](https://goreportcard.com/report/github.com/breml/go-uptime-kuma-client)](https://goreportcard.com/report/github.com/breml/go-uptime-kuma-client)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Go client library for the Uptime Kuma API using Socket.IO for real-time communication.

## Installation

```bash
go get github.com/breml/go-uptime-kuma-client
```

## Requirements

Uptime Kuma 2.x and 3.x, tested against 2.5.6 and the 3.0 nightly image. The
two lines authenticate differently, see [Authentication](#authentication); the
client detects which one it talks to.

Uptime Kuma 1.x is not supported and connecting to it is not expected to work.
It speaks a different login protocol: it never announces `loginRequired`, so the
client cannot tell a server that wants a login from one that does not, and it
reports rejected logins with untranslated messages that the client does not map
to its sentinel errors.

## Features

- **Monitor Management**: HTTP, TCP, Ping, DNS, Redis, PostgreSQL, gRPC, Real Browser,
  and more
- **Notification Providers**: Ntfy, Slack, Teams, Generic, and other notification
  types
- **Tag Management**: Organize monitors with tags
- **Proxy Configuration**: Route monitor requests through HTTP/HTTPS/SOCKS proxies
- **Maintenance Windows**: Schedule maintenance periods
- **Status Pages**: Create and manage public status pages
- **Real-time Updates**: Socket.IO-based event system for state synchronization
- **Authentication**: Username and password, two-factor authentication without a
  human at the keyboard, session token reuse, and servers with authentication
  disabled

## Usage

Here's a complete example showing how to connect to Uptime Kuma, create a notification
provider, and set up an HTTP monitor with notifications:

```go
package main

import (
 "context"
 "log"

 kuma "github.com/breml/go-uptime-kuma-client"
 "github.com/breml/go-uptime-kuma-client/monitor"
 "github.com/breml/go-uptime-kuma-client/notification"
)

func main() {
 ctx := context.Background()

 // Create client and connect to Uptime Kuma
 client, err := kuma.New(ctx, "http://localhost:3001", "admin", "password")
 if err != nil {
  log.Fatal(err)
 }
 defer client.Disconnect()

 // Create a notification provider
 ntfyNotification := notification.Ntfy{
  Base: notification.Base{
   Name:     "My Ntfy Alert",
   IsActive: true,
  },
  NtfyDetails: notification.NtfyDetails{
   ServerURL:            "https://ntfy.sh",
   Topic:                "uptime-alerts",
   Priority:             5,
   AuthenticationMethod: "none",
  },
 }

 notificationID, err := client.CreateNotification(ctx, ntfyNotification)
 if err != nil {
  log.Fatal(err)
 }
 log.Printf("Created notification with ID: %d", notificationID)

 // Create an HTTP monitor that uses the notification
 httpMonitor := &monitor.HTTP{
  Base: monitor.Base{
   Name:            "Example Website",
   Interval:        60,
   NotificationIDs: []int64{notificationID},
  },
  HTTPDetails: monitor.HTTPDetails{
   URL:                 "https://example.com",
   Method:              "GET",
   AcceptedStatusCodes: []string{"200-299"},
  },
 }

 monitorID, err := client.CreateMonitor(ctx, httpMonitor)
 if err != nil {
  log.Fatal(err)
 }
 log.Printf("Created monitor with ID: %d", monitorID)
}
```

## Self-signed certificates

To trust a private CA, pass it as `RootCAs`. To accept a self-signed
certificate, pin its SHA-256 fingerprint (hex, `:` separators optional) with
`WithCertificateFingerprint`; the chain and host name are then not verified,
but the connection is refused unless the server presents exactly that
certificate:

```sh
openssl s_client -connect kuma.example.com:443 </dev/null 2>/dev/null |
  openssl x509 -noout -fingerprint -sha256
```

```go
client, err := kuma.New(ctx, url, "username", "password",
 kuma.WithTLSConfig(&tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12}))

client, err = kuma.New(ctx, url, "username", "password",
 kuma.WithCertificateFingerprint("AB:CD:..."))
```

Both apply to the HTTP requests and the WebSocket connection.

## Authentication

A username and password is the usual way in. An account with two-factor
authentication enabled additionally needs the shared secret, which the client
turns into the one-time code the server asks for:

```go
client, err := kuma.New(ctx, url, "username", "password",
 kuma.WithTOTPSecret(secret))
```

`WithTOTPCode` takes a callback instead, for a secret the client cannot read
itself.

A successful login answers with a session token. Persisting it lets a later
client connect with neither the password nor a one-time code (keep the username
and password configured as well, see below):

```go
token := client.SessionToken()

// ... later, in another process
client, err := kuma.New(ctx, url, "", "", kuma.WithSessionToken(token))
```

This is also how several clients start at once against an account with
two-factor authentication: the server refuses a one-time code it has already
accepted, but takes the same token any number of times.

The token is a bearer credential, and for an account with two-factor
authentication it is a stronger one than the password. Store it the way the
password would be stored; `WithSessionToken` documents what invalidates it. On
a 2.x server it never expires. On a 3.x server it is the session cookie, which
expires after seven days without use and on a password change, so a client
given the username and password alongside falls back to those and reports the
dead token through `SessionTokenRejected` or `WithSessionTokenRejectedCallback`.

A server with authentication disabled logs the client in by itself, so it takes
no credentials at all:

```go
client, err := kuma.New(ctx, url, "", "")
```

On Uptime Kuma 3.x an API key is a further way in, once "API and MCP"
(`apiEnabled`, alpha upstream) is turned on in the server's settings:

```go
client, err := kuma.New(ctx, url, "", "", kuma.WithAPIKey(key))
```

The key goes out as `Authorization: Bearer <key>` (pass it without the
`Bearer` prefix) and the server authenticates the Socket.IO connection by it,
with no login and no session: `SessionToken` returns the empty string and
`Resync` connects again with the key. It is the only credential such a client
has; combining it with a username and password, `WithSessionToken`, a TOTP
option or `WithAutosetup` fails `New`. `WithAPIKey("")` is a no-op.

`New` takes a session the server answers the handshake with as the key being
accepted. The server answers an unknown key and a disabled API alike, reported
as `ErrAPIKeyRejected`. A 3.x server with authentication disabled ignores the
key and opens a session for anyone, so there even an invalid key "succeeds".

A 2.x server accepts API keys for HTTP basic auth on `/metrics` only, never for
the Socket.IO API. `New` fails with `ErrAPIKeyUnsupported` when the scheme is
detected as, or pinned to, `AuthSchemeSocket`, and when a pinned
`AuthSchemeHTTP` meets a 2.x server with authentication disabled. Pinning
`AuthSchemeHTTP` against a 2.x server that wants a login surfaces as
`ErrAPIKeyRejected` instead. A 3.x server that is not set up yet fails with
`ErrAPIKeySetupRequired`. A base URL with userinfo (`user:password@`) cannot be
combined with a key. There is no OIDC, SSO, LDAP or client-certificate login.

### HTTP headers

`WithHeaders` adds headers to every request the client makes, including the
Socket.IO long polling and the WebSocket handshake, for an access proxy in
front of the server:

```go
headers := http.Header{}
headers.Set("CF-Access-Client-Id", id)
headers.Set("CF-Access-Client-Secret", secret)

client, err := kuma.New(ctx, url, "username", "password",
 kuma.WithHeaders(headers))
```

`New` rejects invalid header names and values, and the headers the client or
its transport set themselves: `Host`, `Origin`, `Connection`, `Upgrade`,
`Accept`, `Accept-Encoding`, `Content-Type`, `Content-Length`,
`Transfer-Encoding`, `TE`, `Trailer` and `Sec-WebSocket-*`. A configured
`Cookie` is appended to the client's own cookies. `Authorization` is accepted
for a scheme other than `Bearer`, for a reverse proxy that wants HTTP basic
auth; a Bearer token belongs in `WithAPIKey`.

The headers, the API key included, go to the server's own origin only; other
origins get the User-Agent alone.
`WithUserAgent` replaces the default User-Agent,
`go-uptime-kuma-client/<version>`.

### Redirects

The client follows no redirect that leaves the base URL's origin (scheme, host
and port), whether or not headers or a key are configured. That includes a
server or proxy that upgrades `http` to `https`: a redirect turns the POSTs of
the login and the Socket.IO long polling into bodyless GETs and the WebSocket
handshake follows none, so it never worked reliably. `New` fails with an error
naming both URLs and, where the redirect tells it, the base URL to use instead
(the `https` one, for an upgrade). URLs in errors never show a password from
the base URL.

### Uptime Kuma 2.x and 3.x

Uptime Kuma 3.0 moved its login from the Socket.IO connection to HTTP
([better-auth](https://www.better-auth.com)): the client signs in with a
request, receives a session cookie, and the server authenticates the Socket.IO
connection from that cookie on its handshake. `New` tells the two apart with a
single request before it connects and `Client.AuthScheme` reports the result;
`WithAuthScheme` pins it for a proxy that answers that request itself.

What differs on a 3.x server, beyond the token's expiry:

- The server rate limits logins per client address: three sign-ins per ten
  seconds, ten other authentication requests per minute. Refusals are reported
  as `ErrRateLimited`, wrapped in a `RateLimitedError` that carries the wait the
  server asked for. Several clients started at once should share a session
  token instead of each signing in.
- `Resync` connects again instead of logging in again, because the server has
  no command that resends every list. Commands in flight finish first.
- There is no replay guard on one-time codes, so two clients may log in with the
  same code.
- `GetMonitors` returns compact entries: ID, name, type, active, parent, path
  name, interval and tags, with every other field at its zero value.
  The list identifies monitors; `GetMonitor` and `GetMonitorAs` describe them,
  on both server lines.
- User IDs are opaque strings instead of numbers. The `UserID` field of
  notifications, proxies and Docker hosts is therefore an `ids.UserID`, which
  decodes both forms and writes each back the way it came; `Int64()` returns the
  number a 2.x server uses.

## Documentation

Full documentation available at [pkg.go.dev](https://pkg.go.dev/github.com/breml/go-uptime-kuma-client)

## API Coverage

This client supports the following Uptime Kuma features:

- ✅ Monitors (all types)
- ✅ Notifications
- ✅ Tags
- ✅ Proxies
- ✅ Maintenance Windows
- ✅ Status Pages

This is a work in progress, more Uptime Kuma features might be added in the future.

## License

MIT License - see [LICENSE](LICENSE) for details
