# ACN NAS Message Format — Revised Prototype Specification

## Background

3GPP SA2 contribution S2-2606208 studies 6G core-network support for UE AI-agent identification, registration, discovery, and communication. It includes NAS-based variants for Agent ID assignment, registration/deregistration, and discovery, while leaving the detailed NAS transport and parsing mechanism for later coordination.

This document defines a private prototype format. ACN messages are **not new 5GMM message types**. They are private ACN control PDUs carried inside the existing:

- `UL NAS TRANSPORT` (`0x67`)
- `DL NAS TRANSPORT` (`0x68`)

The outer NAS message uses the existing NAS security procedures. The ACN codec starts only after NAS unprotection and ends before NAS protection.

`Payload Container Type = 0x0E` is a prototype-local assignment. It is not a 3GPP allocation and must be enabled only between ACN-aware UE and network implementations.

## Scope

This specification covers:

| ACN operation | HTTP interface | Transport |
|---|---|---|
| `ACN_AGENT_REGISTER_REQUEST/ACCEPT/REJECT` | `/idm/v1/identity-applications` | NAS |
| `ACN_AGENT_DEREGISTER_REQUEST/ACCEPT/REJECT` | `/acn-agent/v1/agent-deletions` | NAS |
| `ACN_AGENT_PROFILE_UPDATE_REQUEST/ACCEPT/REJECT` | `/arf/v1/agent-cards` | NAS |
| `ACN_AGENT_SEARCH_REQUEST/RESPONSE/REJECT` | `/arf/v1/agent-discoveries`, `/arf/v1/agent-info`, `/acn-agent/v1/owner-agents` | NAS |

Task execution and task termination remain user-plane operations and are outside this document.

## Recommended architecture

| Layer | Implementation |
|---|---|
| 5GMM transport | Reuse existing UL/DL NAS transport codecs |
| ACN identification | Prototype-local Payload Container Type `0x0E` |
| ACN header | Version, message type, transaction ID |
| ACN body | Handwritten `acn` package |
| NAS security | Existing NAS protection/unprotection |
| Transaction tracking | Stateful UE adapter and network service, not codec |

```text
Application / peripheral
        │ HTTP
        ▼
UE ACN NAS Adapter
        │ ACN encode
        ▼
UL NAS TRANSPORT
        │ existing NAS protection
        ▼
AMF / NAS security service
        │ NAS unprotection + transport decode
        ▼
ACN router / IDM / discovery service
```

## Plain and protected NAS

All hexadecimal examples in this document represent the **plain inner 5GMM message**:

- before NAS protection on transmission; or
- after NAS unprotection on reception.

Over the air, standard NAS security may add the security header, MAC, sequence number, and encryption.

## NAS envelope

```text
┌──────────────────────────────────────────────────────────────┐
│ Extended Protocol Discriminator          1 Byte   0x7E       │
│ Security Header Type / Spare             1 Byte   0x00       │
│ 5GMM Message Type                        1 Byte              │
│ Payload Container Type                   1 Byte   0x0E       │
│ Payload Container Length                 2 Byte              │
├──────────────────────────────────────────────────────────────┤
│ ACN Version                              1 Byte   0x01       │
│ ACN Message Type                         1 Byte              │
│ Transaction ID                           1 Byte              │
│ Message-specific fields                  Variable            │
└──────────────────────────────────────────────────────────────┘
```

| Field | Rule |
|---|---|
| 5GMM Message Type | Uplink `0x67`; downlink `0x68` |
| Payload Container Type | Prototype-local `0x0E` |
| Payload Container Length | Length of the complete ACN PDU |
| ACN Version | `0x01` |
| Transaction ID | `0x01`–`0xFF`; `0x00` is reserved |

No fragmentation is defined. The complete ACN PDU must be at most `65,535` bytes. Implementations must expose a configurable `MaxACNPayload`; the controlled demo may set it to `65,535`.

## Transaction rules

- The requester allocates a non-zero Transaction ID.
- The response echoes the request Transaction ID.
- The transaction namespace is per UE ACN context.
- The codec validates only that the value is non-zero.
- Duplicate detection, timeouts, retransmission, completion, and ID reuse are handled above the codec.

## Message types

| Value | Message | Direction |
|---:|---|---|
| `0x01` | `ACN_AGENT_REGISTER_REQUEST` | UL |
| `0x02` | `ACN_AGENT_REGISTER_ACCEPT` | DL |
| `0x03` | `ACN_AGENT_REGISTER_REJECT` | DL |
| `0x04` | `ACN_AGENT_DEREGISTER_REQUEST` | UL |
| `0x05` | `ACN_AGENT_DEREGISTER_ACCEPT` | DL |
| `0x06` | `ACN_AGENT_PROFILE_UPDATE_REQUEST` | UL |
| `0x07` | `ACN_AGENT_PROFILE_UPDATE_ACCEPT` | DL |
| `0x08` | `ACN_AGENT_PROFILE_UPDATE_REJECT` | DL |
| `0x09` | `ACN_AGENT_SEARCH_REQUEST` | UL |
| `0x0A` | `ACN_AGENT_SEARCH_RESPONSE` | DL |
| `0x0B` | `ACN_AGENT_SEARCH_REJECT` | DL |
| `0x0C` | `ACN_AGENT_DEREGISTER_REJECT` | DL |

`0x0C` is appended to preserve the existing prototype values.

## Field encoding

| Encoding | Structure |
|---|---|
| Fixed | Field value only |
| LV | `Length: 1 Byte` + `Value` |
| LV-E | `Length: 2 Byte` + `Value` |
| JSON Container | `JSON Length: 2 Byte` + exact UTF-8 JSON bytes |
| Agent Record | `Record Length: 2 Byte` + record value |

General rules:

- Multi-byte integers use network byte order.
- Text is valid UTF-8.
- Timestamp is `uint64` Unix Epoch milliseconds.
- Public keys and signatures are raw binary. The HTTP adapter Base64-decodes them before ACN encoding.
- Mandatory variable-length fields are present and non-empty.
- Unknown ACN versions, message types, search types, or wrong-direction messages are rejected.
- Trailing bytes after a complete message are rejected.

Default version 1 field limits:

| Field | Limit |
|---|---:|
| Agent ID / Owner ID | 1–1,024 bytes |
| Task ID | 1–1,024 bytes |
| Agent Name | 1–255 bytes |
| Description | 1–2,048 bytes |
| Capability | 1–255 bytes |
| Public Key / Signature | 1–8,192 bytes |
| Region | 1–64 bytes |
| OS | 1–128 bytes |
| Software Version | 1–64 bytes |
| JSON Container | 1–65,535 bytes, additionally bounded by `MaxACNPayload` |

The codec enforces these limits by default. An implementation may configure
smaller limits for a deployment, but increasing them does not change the wire
format and must not exceed the capacity of the corresponding LV or LV-E
encoding.

For version 1, reserved or unknown values of Deregistration Reason, Priority,
Agent Status, Reject Cause, and Failed Field are rejected. A conditional
variable-length field is non-empty when its presence bit is set.

## JSON container rules

Only `vc0` and `vc_list` use JSON.

- The codec validates the bytes with `json.Valid`.
- The codec preserves the exact bytes and does not unmarshal and re-marshal them.
- `vc0` must be a non-empty JSON object.
- `vc_list` must be a non-empty JSON array.
- Compact JSON is recommended but is not treated as canonical JSON.
- If a signature must bind a JSON container, the adapter/IDM contract must bind the signature to the hash of the exact JSON bytes.
- JSON semantic validation, duplicate-key policy, VC proof validation, and credential authorization are outside the codec.

A Go implementation should use `json.RawMessage`.

## Security boundary

NAS security authenticates and protects the UE-to-AMF hop. The ACN backend must also receive trusted UE context from the AMF or adapter, for example SUPI, serving PLMN, and access type.

Payload-carried `owner` and `agent_id` values must be checked against this authenticated context.

Timestamp freshness, acceptable clock skew, future-time handling, and replay
detection are application policy above the codec. The codec validates only the
unsigned 64-bit wire representation.

The ACN codec treats public keys and signatures as opaque bytes. The adapter/IDM contract defines:

- public-key format;
- key algorithm;
- signature algorithm; and
- signed byte sequence.

## `ACN_AGENT_REGISTER_REQUEST`

Carries the complete request from `POST /idm/v1/identity-applications`.

### Message structure

| Field | Encoding | Presence | HTTP field |
|---|---|:---:|---|
| ACN Version | 1 Byte | M | — |
| ACN Message Type | 1 Byte | M | — |
| Transaction ID | 1 Byte | M | — |
| Owner | LV-E UTF-8 | M | `owner` |
| Agent Name | LV-E UTF-8 | M | `name` |
| Agent Public Key | LV-E binary | M | `public_key` after Base64 decode |
| Description | LV-E UTF-8 | M | `description` |
| Timestamp | 8 Byte | M | `timestamp` |
| Signature | LV-E binary | M | `signature` after Base64 decode |
| Region | LV-E UTF-8 | M | `metadata.region` |
| OS | LV-E UTF-8 | M | `metadata.os` |
| Software Version | LV-E UTF-8 | M | `metadata.version` |

`signature_encoding` is not encoded in NAS. The demo HTTP adapter requires `base64`.

### HTTP interface

```http
POST /idm/v1/identity-applications
Content-Type: application/json
```

```json
{
  "owner": "u1",
  "name": "Alice",
  "public_key": "AQIDBA==",
  "description": "Model-X",
  "timestamp": "2026-03-27T13:25:40Z",
  "signature": "c2ln",
  "signature_encoding": "base64",
  "metadata": {
    "region": "CN",
    "os": "Linux",
    "version": "1.0.0"
  }
}
```

### Hexadecimal example

```text
7E 00 67 0E 00 3C
│  │  │  │  └──── Payload Container Length = 0x003C = 60 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── UL NAS TRANSPORT = 0x67
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 01 21
│  │  └────────── Transaction ID = 0x21
│  └───────────── ACN_AGENT_REGISTER_REQUEST = 0x01
└──────────────── ACN Version = 0x01
```

```text
00 02 75 31
└──────────────── Owner = "u1"

00 05 41 6C 69 63 65
└──────────────── Agent Name = "Alice"

00 04 01 02 03 04
└──────────────── Agent Public Key = 4 binary bytes

00 07 4D 6F 64 65 6C 2D 58
└──────────────── Description = "Model-X"

00 00 01 9D 2F 78 D0 20
└──────────────── Timestamp = 2026-03-27T13:25:40Z

00 03 73 69 67
└──────────────── Signature = 3 binary bytes

00 02 43 4E
└──────────────── Region = "CN"

00 05 4C 69 6E 75 78
└──────────────── OS = "Linux"

00 05 31 2E 30 2E 30
└──────────────── Software Version = "1.0.0"
```

Complete message:

```text
7E 00 67 0E 00 3C 01 01 21 00 02 75 31 00 05 41
6C 69 63 65 00 04 01 02 03 04 00 07 4D 6F 64 65
6C 2D 58 00 00 01 9D 2F 78 D0 20 00 03 73 69 67
00 02 43 4E 00 05 4C 69 6E 75 78 00 05 31 2E 30
2E 30
```

## `ACN_AGENT_REGISTER_ACCEPT`

Returns the assigned Agent ID and complete `vc0`. The message type represents `result = success`; no separate result field is encoded.

### Message structure

| Field | Encoding | Presence | HTTP field |
|---|---|:---:|---|
| ACN Version | 1 Byte | M | — |
| ACN Message Type | 1 Byte | M | `result = success` |
| Transaction ID | 1 Byte | M | — |
| Agent ID | LV-E UTF-8 | M | `agent_id` |
| VC0 JSON | JSON Container | M | complete `vc0` object |

### HTTP interface

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "result": "success",
  "agent_id": "a1",
  "vc0": {
    "context": ["c"],
    "id": "cred1",
    "type": ["VC", "SIM"],
    "issuer": "idm",
    "valid_from": "2026",
    "valid_until": "2027",
    "claims": {
      "agent_name": "Alice",
      "agent_id": "a1",
      "agent_attribute": "binding",
      "master_id": "m1",
      "self_id": "s1"
    },
    "proof": {
      "creator": "idm#1",
      "signature_value": "sig"
    }
  }
}
```

### Hexadecimal example

```text
7E 00 68 0E 01 13
│  │  │  │  └──── Payload Container Length = 0x0113 = 275 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── DL NAS TRANSPORT = 0x68
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 02 21
│  │  └────────── Transaction ID = 0x21
│  └───────────── ACN_AGENT_REGISTER_ACCEPT = 0x02
└──────────────── ACN Version = 0x01
```

```text
00 02 61 31
└──────────────── Agent ID = "a1"

01 0A
└──────────────── VC0 JSON Length = 0x010A = 266 Byte
```

Complete message:

```text
7E 00 68 0E 01 13 01 02 21 00 02 61 31 01 0A 7B
22 63 6F 6E 74 65 78 74 22 3A 5B 22 63 22 5D 2C
22 69 64 22 3A 22 63 72 65 64 31 22 2C 22 74 79
70 65 22 3A 5B 22 56 43 22 2C 22 53 49 4D 22 5D
2C 22 69 73 73 75 65 72 22 3A 22 69 64 6D 22 2C
22 76 61 6C 69 64 5F 66 72 6F 6D 22 3A 22 32 30
32 36 22 2C 22 76 61 6C 69 64 5F 75 6E 74 69 6C
22 3A 22 32 30 32 37 22 2C 22 63 6C 61 69 6D 73
22 3A 7B 22 61 67 65 6E 74 5F 6E 61 6D 65 22 3A
22 41 6C 69 63 65 22 2C 22 61 67 65 6E 74 5F 69
64 22 3A 22 61 31 22 2C 22 61 67 65 6E 74 5F 61
74 74 72 69 62 75 74 65 22 3A 22 62 69 6E 64 69
6E 67 22 2C 22 6D 61 73 74 65 72 5F 69 64 22 3A
22 6D 31 22 2C 22 73 65 6C 66 5F 69 64 22 3A 22
73 31 22 7D 2C 22 70 72 6F 6F 66 22 3A 7B 22 63
72 65 61 74 6F 72 22 3A 22 69 64 6D 23 31 22 2C
22 73 69 67 6E 61 74 75 72 65 5F 76 61 6C 75 65
22 3A 22 73 69 67 22 7D 7D
```

## `ACN_AGENT_REGISTER_REJECT`

### Message structure

| Field | Encoding | Presence |
|---|---|:---:|
| ACN Version | 1 Byte | M |
| ACN Message Type | 1 Byte | M |
| Transaction ID | 1 Byte | M |
| Reject Cause | 1 Byte | M |
| Failed Field | 1 Byte | M |

Reject Cause:

| Value | Meaning |
|---:|---|
| `0x01` | Invalid mandatory field |
| `0x02` | Owner not allowed for current UE |
| `0x03` | Invalid public key |
| `0x04` | Invalid signature |
| `0x05` | Invalid timestamp |
| `0x06` | Agent ID allocation failed |
| `0x07` | Internal error |

Failed Field:

| Value | Field |
|---:|---|
| `0x00` | Not specified |
| `0x01` | Owner |
| `0x02` | Agent Name |
| `0x03` | Public Key |
| `0x04` | Description |
| `0x05` | Timestamp |
| `0x06` | Signature |
| `0x07` | Region |
| `0x08` | OS |
| `0x09` | Software Version |

### HTTP interface

```http
HTTP/1.1 403 Forbidden
Content-Type: application/json
```

```json
{
  "result": "failure",
  "code": "SIGNATURE_INVALID"
}
```

### Hexadecimal example

```text
7E 00 68 0E 00 05
│  │  │  │  └──── Payload Container Length = 0x0005 = 5 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── DL NAS TRANSPORT = 0x68
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 03 21
│  │  └────────── Transaction ID = 0x21
│  └───────────── ACN_AGENT_REGISTER_REJECT = 0x03
└──────────────── ACN Version = 0x01
```

```text
04
└──────────────── Reject Cause = Invalid signature

06
└──────────────── Failed Field = Signature
```

Complete message:

```text
7E 00 68 0E 00 05 01 03 21 04 06
```

## `ACN_AGENT_DEREGISTER_REQUEST`

### Message structure

| Field | Encoding | Presence | HTTP field |
|---|---|:---:|---|
| ACN Version | 1 Byte | M | — |
| ACN Message Type | 1 Byte | M | — |
| Transaction ID | 1 Byte | M | — |
| Agent ID | LV-E UTF-8 | M | `agent_id` |
| Deregistration Reason | 1 Byte | M | `reason` |
| Timestamp | 8 Byte | M | `timestamp` |
| Signature | LV-E binary | M | `signature` after Base64 decode |

Deregistration Reason:

| Value | HTTP value |
|---:|---|
| `0x00` | `normal` |
| `0x01` | `uninstalled` |
| `0x02` | `replaced` |
| `0x03` | `user_request` |
| `0x04` | `security_event` |
| `0x05` | `retired` |
| `0xFF` | `other` |

### HTTP interface

```http
POST /acn-agent/v1/agent-deletions
Content-Type: application/json
```

```json
{
  "agent_id": "a1",
  "reason": "retired",
  "timestamp": "2024-03-23T12:00:00Z",
  "signature": "c2ln",
  "signature_encoding": "base64"
}
```

### Hexadecimal example

```text
7E 00 67 0E 00 15
│  │  │  │  └──── Payload Container Length = 0x0015 = 21 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── UL NAS TRANSPORT = 0x67
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 04 22
│  │  └────────── Transaction ID = 0x22
│  └───────────── ACN_AGENT_DEREGISTER_REQUEST = 0x04
└──────────────── ACN Version = 0x01
```

```text
00 02 61 31
└──────────────── Agent ID = "a1"

05
└──────────────── Deregistration Reason = retired

00 00 01 8E 6B 2E 9A 00
└──────────────── Timestamp = 2024-03-23T12:00:00Z

00 03 73 69 67
└──────────────── Signature = 3 binary bytes
```

Complete message:

```text
7E 00 67 0E 00 15 01 04 22 00 02 61 31 05 00 00
01 8E 6B 2E 9A 00 00 03 73 69 67
```

## `ACN_AGENT_DEREGISTER_ACCEPT`

The operation is idempotent. An already deregistered Agent ID may receive `DEREGISTER_ACCEPT`.

### Message structure

| Field | Encoding | Presence |
|---|---|:---:|
| ACN Version | 1 Byte | M |
| ACN Message Type | 1 Byte | M |
| Transaction ID | 1 Byte | M |

### HTTP interface

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{}
```

### Hexadecimal example

```text
7E 00 68 0E 00 03
│  │  │  │  └──── Payload Container Length = 0x0003 = 3 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── DL NAS TRANSPORT = 0x68
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 05 22
│  │  └────────── Transaction ID = 0x22
│  └───────────── ACN_AGENT_DEREGISTER_ACCEPT = 0x05
└──────────────── ACN Version = 0x01
```

Complete message:

```text
7E 00 68 0E 00 03 01 05 22
```

## `ACN_AGENT_DEREGISTER_REJECT`

### Message structure

| Field | Encoding | Presence |
|---|---|:---:|
| ACN Version | 1 Byte | M |
| ACN Message Type | 1 Byte | M |
| Transaction ID | 1 Byte | M |
| Reject Cause | 1 Byte | M |
| Failed Field | 1 Byte | M |

Reject Cause:

| Value | Meaning |
|---:|---|
| `0x01` | Unknown Agent ID |
| `0x02` | Agent not bound to current UE |
| `0x03` | Operation not authorized |
| `0x04` | Invalid signature |
| `0x05` | Invalid timestamp |
| `0x06` | Backend failure |

Failed Field:

| Value | Field |
|---:|---|
| `0x00` | Not specified |
| `0x01` | Agent ID |
| `0x02` | Reason |
| `0x03` | Timestamp |
| `0x04` | Signature |
| `0x05` | UE/Agent binding |

### HTTP interface

```http
HTTP/1.1 403 Forbidden
Content-Type: application/json
```

```json
{
  "result": "failure",
  "code": "SIGNATURE_INVALID"
}
```

### Hexadecimal example

```text
7E 00 68 0E 00 05
│  │  │  │  └──── Payload Container Length = 0x0005 = 5 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── DL NAS TRANSPORT = 0x68
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 0C 22
│  │  └────────── Transaction ID = 0x22
│  └───────────── ACN_AGENT_DEREGISTER_REJECT = 0x0C
└──────────────── ACN Version = 0x01
```

```text
04
└──────────────── Reject Cause = Invalid signature

04
└──────────────── Failed Field = Signature
```

Complete message:

```text
7E 00 68 0E 00 05 01 0C 22 04 04
```

## `ACN_AGENT_PROFILE_UPDATE_REQUEST`

Only the nested `vc_list` array is JSON-encoded. All other fields use binary NAS encoding.

### Message structure

| Field | Encoding | Presence | HTTP field |
|---|---|:---:|---|
| ACN Version | 1 Byte | M | — |
| ACN Message Type | 1 Byte | M | — |
| Transaction ID | 1 Byte | M | — |
| Agent ID | LV-E UTF-8 | M | `agent_id` |
| Priority | 1 Byte | M | `priority` |
| Timestamp | 8 Byte | M | `timestamp` |
| Signature | LV-E binary | M | `signature` after Base64 decode |
| VC List JSON | JSON Container | M | complete `vc_list` array |

Priority:

| Value | Meaning |
|---:|---|
| `0x00` | Unspecified |
| `0x01` | High |
| `0x02` | Normal |
| `0x03` | Low |
| `0x04`–`0xFF` | Reserved |

### HTTP interface

```http
POST /arf/v1/agent-cards
Content-Type: application/json
```

```json
{
  "agent_id": "a1",
  "priority": 2,
  "timestamp": "2024-03-23T12:00:00Z",
  "signature": "c2ln",
  "signature_encoding": "base64",
  "vc_list": [
    {
      "id": "cred1",
      "type": ["VC"],
      "claims": {
        "agent_attribute": "service"
      }
    },
    {
      "id": "cred2",
      "type": ["VC"],
      "claims": {
        "agent_attribute": "fall"
      }
    }
  ]
}
```

### Hexadecimal example

```text
7E 00 67 0E 00 9D
│  │  │  │  └──── Payload Container Length = 0x009D = 157 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── UL NAS TRANSPORT = 0x67
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 06 23
│  │  └────────── Transaction ID = 0x23
│  └───────────── ACN_AGENT_PROFILE_UPDATE_REQUEST = 0x06
└──────────────── ACN Version = 0x01
```

```text
00 02 61 31
└──────────────── Agent ID = "a1"

02
└──────────────── Priority = Normal

00 00 01 8E 6B 2E 9A 00
└──────────────── Timestamp = 2024-03-23T12:00:00Z

00 03 73 69 67
└──────────────── Signature = 3 binary bytes

00 86
└──────────────── VC List JSON Length = 0x0086 = 134 Byte
```

Complete message:

```text
7E 00 67 0E 00 9D 01 06 23 00 02 61 31 02 00 00
01 8E 6B 2E 9A 00 00 03 73 69 67 00 86 5B 7B 22
69 64 22 3A 22 63 72 65 64 31 22 2C 22 74 79 70
65 22 3A 5B 22 56 43 22 5D 2C 22 63 6C 61 69 6D
73 22 3A 7B 22 61 67 65 6E 74 5F 61 74 74 72 69
62 75 74 65 22 3A 22 73 65 72 76 69 63 65 22 7D
7D 2C 7B 22 69 64 22 3A 22 63 72 65 64 32 22 2C
22 74 79 70 65 22 3A 5B 22 56 43 22 5D 2C 22 63
6C 61 69 6D 73 22 3A 7B 22 61 67 65 6E 74 5F 61
74 74 72 69 62 75 74 65 22 3A 22 66 61 6C 6C 22
7D 7D 5D
```

## `ACN_AGENT_PROFILE_UPDATE_ACCEPT`

### Message structure

| Field | Encoding | Presence |
|---|---|:---:|
| ACN Version | 1 Byte | M |
| ACN Message Type | 1 Byte | M |
| Transaction ID | 1 Byte | M |

### HTTP interface

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{}
```

### Hexadecimal example

```text
7E 00 68 0E 00 03
│  │  │  │  └──── Payload Container Length = 0x0003 = 3 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── DL NAS TRANSPORT = 0x68
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 07 23
│  │  └────────── Transaction ID = 0x23
│  └───────────── ACN_AGENT_PROFILE_UPDATE_ACCEPT = 0x07
└──────────────── ACN Version = 0x01
```

Complete message:

```text
7E 00 68 0E 00 03 01 07 23
```

## `ACN_AGENT_PROFILE_UPDATE_REJECT`

### Message structure

| Field | Encoding | Presence |
|---|---|:---:|
| ACN Version | 1 Byte | M |
| ACN Message Type | 1 Byte | M |
| Transaction ID | 1 Byte | M |
| Reject Cause | 1 Byte | M |
| Failed Field | 1 Byte | M |

Reject Cause:

| Value | Meaning |
|---:|---|
| `0x01` | Invalid Agent ID |
| `0x02` | Update not authorized |
| `0x03` | Invalid priority |
| `0x04` | Invalid VC List |
| `0x05` | Invalid signature |
| `0x06` | Invalid timestamp |
| `0x07` | Internal error |

Failed Field:

| Value | Field |
|---:|---|
| `0x00` | Not specified |
| `0x01` | Agent ID |
| `0x02` | Priority |
| `0x03` | Timestamp |
| `0x04` | Signature |
| `0x05` | VC List |

### HTTP interface

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
  "result": "failure",
  "code": "VC_LIST_INVALID"
}
```

### Hexadecimal example

```text
7E 00 68 0E 00 05
│  │  │  │  └──── Payload Container Length = 0x0005 = 5 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── DL NAS TRANSPORT = 0x68
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 08 23
│  │  └────────── Transaction ID = 0x23
│  └───────────── ACN_AGENT_PROFILE_UPDATE_REJECT = 0x08
└──────────────── ACN Version = 0x01
```

```text
04
└──────────────── Reject Cause = Invalid VC List

05
└──────────────── Failed Field = VC List
```

Complete message:

```text
7E 00 68 0E 00 05 01 08 23 04 05
```

## `ACN_AGENT_SEARCH_REQUEST`

One message type covers capability discovery, exact Agent information, and Owner-agent enumeration.

### Common structure

| Field | Encoding | Presence |
|---|---|:---:|
| ACN Version | 1 Byte | M |
| ACN Message Type | 1 Byte | M |
| Transaction ID | 1 Byte | M |
| Search Type | 1 Byte | M |
| Search-specific fields | Variable | M |

Search Type:

| Value | Operation |
|---:|---|
| `0x01` | Capability Discovery |
| `0x02` | Agent Info |
| `0x03` | Owner Agents |

Unknown search types are rejected.

### Capability Discovery

Fields:

| Field | Encoding | Presence |
|---|---|:---:|
| Source Agent ID | LV-E UTF-8 | M |
| Task ID | LV-E UTF-8 | M |
| Timestamp | 8 Byte | M |
| Capability Count | 1 Byte | M |
| Required Capability | Repeated LV UTF-8 | M |

`Capability Count` must be at least one.

HTTP interface:

```http
POST /arf/v1/agent-discoveries
Content-Type: application/json
```

```json
{
  "task_id": "t1",
  "agent_id": "a1",
  "required_capabilities": ["camera", "radar", "four_legs"],
  "timestamp": "2024-03-23T12:00:00Z"
}
```

Hexadecimal example:

```text
7E 00 67 0E 00 2C
│  │  │  │  └──── Payload Container Length = 0x002C = 44 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── UL NAS TRANSPORT = 0x67
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 09 24
│  │  └────────── Transaction ID = 0x24
│  └───────────── ACN_AGENT_SEARCH_REQUEST = 0x09
└──────────────── ACN Version = 0x01
```

```text
01
└──────────────── Search Type = Capability Discovery (0x01)
```

Complete message:

```text
7E 00 67 0E 00 2C 01 09 24 01 00 02 61 31 00 02
74 31 00 00 01 8E 6B 2E 9A 00 03 06 63 61 6D 65
72 61 05 72 61 64 61 72 09 66 6F 75 72 5F 6C 65
67 73
```

### Agent Info

Fields:

| Field | Encoding | Presence |
|---|---|:---:|
| Target Agent ID | LV-E UTF-8 | M |

HTTP interface:

```http
POST /arf/v1/agent-info
Content-Type: application/json
```

```json
{
  "agent_id": "agent-111"
}
```

Hexadecimal example:

```text
7E 00 67 0E 00 0F
│  │  │  │  └──── Payload Container Length = 0x000F = 15 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── UL NAS TRANSPORT = 0x67
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 09 25
│  │  └────────── Transaction ID = 0x25
│  └───────────── ACN_AGENT_SEARCH_REQUEST = 0x09
└──────────────── ACN Version = 0x01
```

```text
02
└──────────────── Search Type = Agent Info (0x02)
```

Complete message:

```text
7E 00 67 0E 00 0F 01 09 25 02 00 09 61 67 65 6E
74 2D 31 31 31
```

### Owner Agents

Fields:

| Field | Encoding | Presence |
|---|---|:---:|
| Owner ID | LV-E UTF-8 | M |

HTTP interface:

```http
POST /acn-agent/v1/owner-agents
Content-Type: application/json
```

```json
{
  "owner": "owner-123"
}
```

Hexadecimal example:

```text
7E 00 67 0E 00 0F
│  │  │  │  └──── Payload Container Length = 0x000F = 15 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── UL NAS TRANSPORT = 0x67
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 09 26
│  │  └────────── Transaction ID = 0x26
│  └───────────── ACN_AGENT_SEARCH_REQUEST = 0x09
└──────────────── ACN Version = 0x01
```

```text
03
└──────────────── Search Type = Owner Agents (0x03)
```

Complete message:

```text
7E 00 67 0E 00 0F 01 09 26 03 00 09 6F 77 6E 65
72 2D 31 32 33
```

## `ACN_AGENT_SEARCH_RESPONSE`

### Message structure

| Field | Encoding | Presence |
|---|---|:---:|
| ACN Version | 1 Byte | M |
| ACN Message Type | 1 Byte | M |
| Transaction ID | 1 Byte | M |
| Search Type | 1 Byte | M |
| Result Count | 2 Byte | M |
| Agent Record | Repeated | C |

The response Search Type must equal the request Search Type.

Cardinality:

| Search Type | Allowed records |
|---|---:|
| Capability Discovery | `0..N` |
| Agent Info | `0..1` |
| Owner Agents | `0..N` |

`Result Count` must equal the encoded Agent Record count.

### Agent Record

```text
Record Length                  2 Byte
Agent ID                       LV-E UTF-8
Presence Bitmap                1 Byte
Agent Name                     LV-E UTF-8      conditional
Description                    LV-E UTF-8      conditional
Agent Status                   1 Byte          conditional
Priority                       1 Byte          conditional
Capability Count               1 Byte          conditional
Capability                     LV UTF-8        repeated
```

`Record Length` counts only the bytes after the two-byte Record Length field. It does not include its own prefix.

Presence Bitmap:

| Bit | Field |
|---:|---|
| 0 | Agent Name |
| 1 | Description |
| 2 | Agent Status |
| 3 | Priority |
| 4 | Capabilities |
| 5–7 | Reserved; must be zero |

Agent Status:

| Value | Meaning |
|---:|---|
| `0x00` | Unknown |
| `0x01` | Offline |
| `0x02` | Online |
| `0x03` | Busy |

The decoder must consume exactly `Record Length` bytes for each record and reject reserved presence bits or trailing bytes.

When the Capabilities presence bit is set, Capability Count is in the range
`1..255` and exactly that many non-empty Capability values follow. When the bit
is clear, neither Capability Count nor Capability values are encoded.

### Capability Discovery response

HTTP interface:

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "total": 1,
  "agents": [
    {
      "agent_id": "agent-222",
      "agent_name": "Inspection Robot B",
      "description": "Model-X",
      "agent_status": "online",
      "priority": 2,
      "agent_capabilities": ["camera", "radar", "four_legs"]
    }
  ]
}
```

Hexadecimal example:

```text
7E 00 68 0E 00 4B
│  │  │  │  └──── Payload Container Length = 0x004B = 75 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── DL NAS TRANSPORT = 0x68
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 0A 24
│  │  └────────── Transaction ID = 0x24
│  └───────────── ACN_AGENT_SEARCH_RESPONSE = 0x0A
└──────────────── ACN Version = 0x01
```

```text
01
└──────────────── Search Type = Capability Discovery (0x01)
```

```text
0001
└──────────────── Result Count = 1
```

Complete message:

```text
7E 00 68 0E 00 4B 01 0A 24 01 00 01 00 43 00 09
61 67 65 6E 74 2D 32 32 32 1F 00 12 49 6E 73 70
65 63 74 69 6F 6E 20 52 6F 62 6F 74 20 42 00 07
4D 6F 64 65 6C 2D 58 02 02 03 06 63 61 6D 65 72
61 05 72 61 64 61 72 09 66 6F 75 72 5F 6C 65 67
73
```

### Agent Info response

HTTP interface:

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "agent_id": "agent-111",
  "agent_name": "Inspection Robot A",
  "agent_status": "online",
  "agent_capabilities": ["camera", "monitoring", "night_vision"],
  "priority": 1
}
```

Hexadecimal example:

```text
7E 00 68 0E 00 4A
│  │  │  │  └──── Payload Container Length = 0x004A = 74 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── DL NAS TRANSPORT = 0x68
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 0A 25
│  │  └────────── Transaction ID = 0x25
│  └───────────── ACN_AGENT_SEARCH_RESPONSE = 0x0A
└──────────────── ACN Version = 0x01
```

```text
02
└──────────────── Search Type = Agent Info (0x02)
```

```text
0001
└──────────────── Result Count = 1
```

Complete message:

```text
7E 00 68 0E 00 4A 01 0A 25 02 00 01 00 42 00 09
61 67 65 6E 74 2D 31 31 31 1D 00 12 49 6E 73 70
65 63 74 69 6F 6E 20 52 6F 62 6F 74 20 41 02 01
03 06 63 61 6D 65 72 61 0A 6D 6F 6E 69 74 6F 72
69 6E 67 0C 6E 69 67 68 74 5F 76 69 73 69 6F 6E
```

### Owner Agents response

The Owner ID is not repeated in NAS. The HTTP adapter reconstructs it from the pending request.

HTTP interface:

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "owner": "owner-123",
  "total": 2,
  "agents": [
    {
      "agent_id": "agent-111",
      "agent_name": "Inspection Robot A",
      "description": "Model-X"
    },
    {
      "agent_id": "agent-222",
      "agent_name": "Inspection Robot B",
      "description": "Model-X"
    }
  ]
}
```

Hexadecimal example:

```text
7E 00 68 0E 00 5C
│  │  │  │  └──── Payload Container Length = 0x005C = 92 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── DL NAS TRANSPORT = 0x68
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 0A 26
│  │  └────────── Transaction ID = 0x26
│  └───────────── ACN_AGENT_SEARCH_RESPONSE = 0x0A
└──────────────── ACN Version = 0x01
```

```text
03
└──────────────── Search Type = Owner Agents (0x03)
```

```text
0002
└──────────────── Result Count = 2
```

Complete message:

```text
7E 00 68 0E 00 5C 01 0A 26 03 00 02 00 29 00 09
61 67 65 6E 74 2D 31 31 31 03 00 12 49 6E 73 70
65 63 74 69 6F 6E 20 52 6F 62 6F 74 20 41 00 07
4D 6F 64 65 6C 2D 58 00 29 00 09 61 67 65 6E 74
2D 32 32 32 03 00 12 49 6E 73 70 65 63 74 69 6F
6E 20 52 6F 62 6F 74 20 42 00 07 4D 6F 64 65 6C
2D 58
```

### No-match response

HTTP interface:

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{}
```

Hexadecimal example:

```text
7E 00 68 0E 00 06
│  │  │  │  └──── Payload Container Length = 0x0006 = 6 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── DL NAS TRANSPORT = 0x68
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 0A 25
│  │  └────────── Transaction ID = 0x25
│  └───────────── ACN_AGENT_SEARCH_RESPONSE = 0x0A
└──────────────── ACN Version = 0x01
```

```text
02
└──────────────── Search Type = Agent Info (0x02)
```

```text
0000
└──────────────── Result Count = 0
```

Complete message:

```text
7E 00 68 0E 00 06 01 0A 25 02 00 00
```

## `ACN_AGENT_SEARCH_REJECT`

### Message structure

| Field | Encoding | Presence |
|---|---|:---:|
| ACN Version | 1 Byte | M |
| ACN Message Type | 1 Byte | M |
| Transaction ID | 1 Byte | M |
| Reject Cause | 1 Byte | M |
| Failed Field | 1 Byte | M |

Reject Cause:

| Value | Meaning |
|---:|---|
| `0x01` | Invalid Search Type |
| `0x02` | Missing or invalid mandatory field |
| `0x03` | Request not authorized |
| `0x04` | Invalid source Agent |
| `0x05` | Invalid target Agent or Owner |
| `0x06` | Internal error |

Failed Field:

| Value | Field |
|---:|---|
| `0x00` | Not specified |
| `0x01` | Source Agent ID |
| `0x02` | Task ID |
| `0x03` | Target Agent ID |
| `0x04` | Owner ID |
| `0x05` | Required Capability |
| `0x06` | Timestamp |

### HTTP interface

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
  "result": "failure",
  "code": "INVALID_CAPABILITY"
}
```

### Hexadecimal example

```text
7E 00 68 0E 00 05
│  │  │  │  └──── Payload Container Length = 0x0005 = 5 Byte
│  │  │  └─────── Payload Container Type = 0x0E
│  │  └────────── DL NAS TRANSPORT = 0x68
│  └───────────── Security Header Type / Spare = 0x00
└──────────────── EPD = 0x7E
```

```text
01 0B 24
│  │  └────────── Transaction ID = 0x24
│  └───────────── ACN_AGENT_SEARCH_REJECT = 0x0B
└──────────────── ACN Version = 0x01
```

```text
02
└──────────────── Reject Cause = Missing or invalid mandatory field

05
└──────────────── Failed Field = Required Capability
```

Complete message:

```text
7E 00 68 0E 00 05 01 0B 24 02 05
```

## Implementation package

Recommended layout:

```text
acn/
  types.go
  message.go
  codec.go
  errors.go
  register.go
  deregister.go
  profile.go
  search.go
  transport.go
  codec_test.go
  fuzz_test.go
```

The 3GPP table generator should not be used for ACN. ACN contains discriminated bodies, repeated records, presence maps, and embedded JSON that are better implemented by handwritten codecs.

Recommended API:

```go
payload, err := acn.Marshal(request)
nasMessage, err := acn.WrapUplink(payload)

payload, err := acn.ExtractUplink(plainNASMessage)
request, err := acn.Unmarshal(acn.Uplink, payload)
```

For downlink:

```go
payload, err := acn.Marshal(response)
nasMessage, err := acn.WrapDownlink(payload)

payload, err := acn.ExtractDownlink(plainNASMessage)
response, err := acn.Unmarshal(acn.Downlink, payload)
```

The ACN package does not perform NAS security processing and does not accept protected over-the-air wire bytes directly.

## Repository integration

The existing generated NAS transport files do not need ACN-specific message definitions if they preserve the Payload Container Type value and expose the complete Payload Container:

```text
nasMessage/NAS_ULNASTransport.go
nasMessage/NAS_DLNASTransport.go
nasType/NAS_PayloadContainer.go
```

The existing NAS security service should protect and unprotect ACN transport messages without ACN-specific cryptographic changes.

If the current codec service rejects `ULNASTransport` after decoding, end-to-end integration requires one of:

- extending the service and its API to expose ACN payloads; or
- letting the ACN adapter extract and decode the transport message directly after NAS unprotection.

## Decoder validation

The decoder must reject:

- Payload Container Length mismatch;
- payload larger than `MaxACNPayload`;
- Transaction ID `0x00`;
- unknown ACN version or message type;
- wrong UL/DL direction;
- zero-length mandatory fields;
- invalid UTF-8;
- invalid JSON;
- empty `vc0` object;
- empty `vc_list` array;
- unknown Search Type;
- Capability Discovery with zero capabilities;
- Agent Info response with more than one record;
- Result Count mismatch;
- malformed Agent Record length;
- reserved presence-bitmap bits;
- trailing bytes.

## Test requirements

Add golden encode/decode tests for every hexadecimal example, plus:

- malformed outer and inner lengths;
- maximum and over-maximum payloads;
- wrong direction;
- unknown version/type;
- zero Transaction ID;
- zero-length mandatory fields;
- invalid UTF-8;
- invalid and empty JSON containers;
- search cardinality violations;
- record truncation and overrun;
- reserved presence bits;
- trailing bytes;
- fuzz decode/encode round trips.

## References

- 3GPP SA2 contribution S2-2606208, KI#19, solution update for 6G CN support of UE AI-agent communication.
- 3GPP TS 24.501 archive: `https://www.3gpp.org/ftp/Specs/archive/24_series/24.501/`
