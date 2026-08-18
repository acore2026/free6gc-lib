# NGAP Codec for Rust

`ngap-codec` provides typed NGAP Aligned PER (APER) encoding and decoding
compatible with the Go codec in this repository. ASN.1 types are generated at
build time from the checked-in minimal schema, so using the crate does not
require an external ASN.1 compiler. The complete Release 15.8.0 schema is
retained beside it as the source for future additions.

## Supported Surface

The top-level API is tailored to NG Setup, initial registration, and PDU
session establishment. All generated transfer and embedded container types
implement `AperMessage`. Unknown or out-of-scope procedure and IE identifiers,
malformed messages, and truncated input return decode errors.

## Usage

```rust
use ngap_codec::{decode, encode, NgapPdu};

fn round_trip(input: &[u8]) -> ngap_codec::Result<(NgapPdu, Vec<u8>)> {
    let pdu = decode(input)?;
    let encoded = encode(&pdu)?;
    Ok((pdu, encoded))
}
```

Transfer containers implement `AperMessage` directly:

```no_run
use ngap_codec::{types::PDUSessionResourceSetupRequestTransfer, AperMessage};

# let bytes: &[u8] = &[];
let transfer = PDUSessionResourceSetupRequestTransfer::decode_aper(bytes)?;
let encoded = transfer.encode_aper()?;
# Ok::<(), ngap_codec::CodecError>(())
```

The NGAP SCTP payload protocol identifier is exported as `PPID`.

## Registration and Session Establishment

The codec covers the NGAP messages used by the initial registration and
UE-requested PDU session establishment paths:

| Procedure | NGAP messages and containers |
| --- | --- |
| Initial registration | Initial UE Message, Downlink/Uplink NAS Transport, Initial Context Setup Request/Response/Failure |
| PDU session establishment | PDU Session Resource Setup Request/Response and request, response, or unsuccessful transfer containers |

PDU session setup carried inside Initial Context Setup is supported through
`PDUSessionResourceSetupListCxtReq` and
`PDUSessionResourceSetupListCxtRes`. Optional contingency messages such as NAS
Non-Delivery Indication and UE Radio Capability Info Indication are also part
of the generated top-level schema.

A complete PDU Session Resource Setup Request construction example is available
under `examples/`:

```text
cargo run --manifest-path rust/Cargo.toml \
  --example pdu_session_resource_setup_request
```

## Development

Run from the repository root:

```text
cargo fmt --manifest-path rust/Cargo.toml --check
cargo clippy --manifest-path rust/Cargo.toml --all-targets -- -D warnings
cargo test --manifest-path rust/Cargo.toml
go test ./...
```

Rust and Go interoperability tests consume the same canonical vectors from
`testdata/interop/`. Add a shared vector and assertions in both suites when
changing the schema or codec behavior.
