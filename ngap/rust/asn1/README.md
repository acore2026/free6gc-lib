# ASN.1 Schema

`ngap-15.8.0.asn1` is the complete repository-owned reference schema. It
originated from a local UERANSIM NGAP schema based on 3GPP TS 38.413 release
15.8.0 and was normalized to UTF-8.

`ngap-15.8.0-min.asn1` is the schema compiled by `rust/build.rs`. Its dependency
closure is limited to NG Setup, initial registration, PDU session
establishment, and their error and transfer types.

## Adding a Message

Use the full schema as the source and update the minimal schema in this order:

1. Add the elementary procedure and its procedure-code constant.
2. Add the applicable initiating, successful, and unsuccessful message
   definitions, then register the procedure in the appropriate class set.
3. Copy every referenced protocol-IE object set, IE type, container,
   extension, constraint constant, and common data type until the dependency
   closure is complete.
4. Add Go/Rust APER vectors for each new top-level outcome and embedded
   transfer type.

The compiler reports unresolved dependencies, but successful compilation alone
does not prove that an information-object set is complete. Always decode and
re-encode a Go-produced vector exactly.

The Rust build generates types into Cargo's `OUT_DIR`; generated source is
never checked in or edited directly. After changing the minimal schema, verify:

```text
cargo test --manifest-path rust/Cargo.toml
go test ./...
```
