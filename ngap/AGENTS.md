# Repository Guidelines

## Project Structure & Module Organization

This repository provides NGAP encoding and decoding for Go and a compatible
Rust codec.

- `ngap.go` contains the public encoder, decoder, and PDU-printing helpers.
- `ngapType/` contains the protocol PDU, information-element, and transfer types. Keep ASN.1/PER tags and presence constants aligned with the relevant 3GPP specification.
- `ngapConvert/` converts NGAP values to and from `github.com/acore2026/free6gc-lib/openapi/models`.
- `rust/` contains the Rust crate, vendored ASN.1 schema, and Rust integration tests.
- `logger/` exposes the library logger and its configuration helpers.
- `testdata/interop/` contains APER vectors shared by the Go and Rust suites.
- `ngap_test.go` contains unit and regression tests. `fuzz_test.go` and `testdata/fuzz/FuzzNGAP/` contain the decoder fuzz target and seed corpus.
- `.github/workflows/` and `.golangci.yml` define CI and lint expectations.

## Build, Test, and Development Commands

Use Go 1.25.5, matching `go.mod` and CI.

- `go build ./...` compiles every package.
- `go test ./...` runs the full unit-test suite.
- `go test -run TestDecoder ./...` runs a focused regression test.
- `go test -fuzz=FuzzNGAP -fuzztime=30s` exercises the decoder against generated inputs and the checked-in corpus.
- `golangci-lint run` runs the configured linters and formatters; CI uses golangci-lint v2.7.2.
- `go mod tidy` normalizes dependencies after import changes. Review both `go.mod` and `go.sum`.
- `cargo test --manifest-path rust/Cargo.toml` builds the generated Rust types and runs interoperability tests.
- `cargo fmt --manifest-path rust/Cargo.toml --check` and `cargo clippy --manifest-path rust/Cargo.toml --all-targets -- -D warnings` enforce Rust style.

## Coding Style & Naming Conventions

Format Go with `gofmt`; lint also enforces `gofumpt` and `gci`. Use Go naming and preserve initialisms such as `NGAP`, `AMF`, and `PLMN`. Format Rust with `rustfmt`; use `snake_case` functions/modules and `PascalCase` types. Do not edit generated Rust; update `rust/asn1/ngap-15.8.0-min.asn1` from the full schema, or change the wrapper API.

## Testing Guidelines

Name Go tests `TestXxx` and Rust integration tests after behavior in `snake_case`. Add regression cases for malformed or newly supported PDUs. Codec changes must preserve shared vectors exactly and should pass a short Go fuzz run. Add durable crash inputs under `testdata/fuzz/FuzzNGAP/`.

## Commit & Pull Request Guidelines

Recent history favors concise, imperative Conventional Commit subjects such as `fix: ...`, `feat: ...`, `test: ...`, and `chore: ...`. Keep each commit focused. Pull requests should explain the protocol behavior changed, identify relevant 3GPP fields or tags, link issues, and include test evidence. Ensure build, tests, and lint pass; screenshots are unnecessary for this library.
