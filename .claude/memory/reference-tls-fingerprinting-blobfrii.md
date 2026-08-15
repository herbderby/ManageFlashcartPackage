---
name: reference-tls-fingerprinting-blobfrii
description: Why Go's standard crypto/tls fails against the blobfrii CDN, how uTLS with HelloChrome_Auto works around the JA3 fingerprint check, and the HTTP/2 ALPN gotcha.
metadata:
  type: reference
---

# TLS Fingerprinting on `us.dl.blobfrii.com`

## The problem

- `us.dl.blobfrii.com` performs **JA3 TLS fingerprinting** during
  the handshake.
- Go's standard `crypto/tls` produces a fingerprint that the CDN
  rejects — the connection is killed with **TCP RST during the
  TLS handshake** before any application data flows.
- Default `net/http` clients therefore cannot reach this CDN at
  all. The failure looks like a network problem (connection
  reset) but is actually intentional protocol-level rejection.

## The fix

- Use **`github.com/refraction-networking/utls`** to construct
  TLS connections that mimic a known browser fingerprint.
- Apply the `HelloChrome_Auto` ClientHello spec to match current
  Chrome's TLS signature.
- The CDN accepts the Chrome fingerprint and the handshake
  completes normally.

## The HTTP/2 ALPN gotcha

- The Chrome fingerprint advertises **both `h2` and `http/1.1`**
  in ALPN by default.
- Go's `net/http` `Transport` **cannot do HTTP/2** over a custom
  `DialTLS` connection — its HTTP/2 client requires control of
  the TLS handshake itself.
- If the CDN negotiates `h2` (because Chrome offered it), the
  connection succeeds at the TLS layer but no HTTP/2 client is
  there to drive it — requests hang.
- **Fix:** modify the `HelloChrome_Auto` spec to advertise
  **HTTP/1.1 only** in ALPN. The CDN happily falls back to
  HTTP/1.1, and `net/http`'s standard `Transport` handles it
  fine.

## Cost

- The `utls` dependency adds **~800 KB** to the binary. This is
  a real cost — uTLS pulls in its own copy of the TLS state
  machine.
- Worth it: without this workaround, the tool simply cannot
  fetch from blobfrii at all.
