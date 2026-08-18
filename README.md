# Free6GC Library (`free6gc-lib`)

Consolidated core 5G protocols, OpenAPI models, and support libraries for Free6GC.

## Included Modules

- [`aper/`](./aper): ASN.1 ALIGNED PER (Packed Encoding Rules) encoder and decoder.
- [`openapi/`](./openapi): 3GPP 5G OpenAPI models and client SDKs.
- [`nas/`](./nas): Non-Access Stratum (5GS NAS) protocol encoder/decoder, security context, and Agent Coordination Network (`acn`) extensions.
- [`ngap/`](./ngap): Next Generation Application Protocol (NGAP) codec.

## Installation

```sh
go get github.com/acore2026/free6gc-lib
```

## Usage Example

```go
import (
    "github.com/acore2026/free6gc-lib/aper"
    "github.com/acore2026/free6gc-lib/nas"
    "github.com/acore2026/free6gc-lib/nas/nasMessage"
    "github.com/acore2026/free6gc-lib/ngap"
    "github.com/acore2026/free6gc-lib/openapi/models"
)
```

## Running Tests

```sh
go test ./...
```
