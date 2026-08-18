# NAS fork provenance

This module is the in-tree Free6GC fork of `github.com/acore2026/free6gc-lib/nas`.

- Imported version: `v1.2.4`
- Imported commit: `485e516130fdd1617335958a8f5a3dd20dffc159`
- Original upstream: `github.com/free5gc/nas`
- License: Apache-2.0 (`LICENSE.txt`)

The module path intentionally remains `github.com/acore2026/free6gc-lib/nas`. Free6GC
consumer modules replace that path with the sibling `free6gc-nas` directory so
standalone Go commands and container builds use this source tree.

The repository contains generated NAS message and information-element code.
Review `generate.sh` before regeneration because it downloads a pinned 3GPP
archive and replaces generated source and testdata directories.
