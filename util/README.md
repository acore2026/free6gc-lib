# Free6GC Utilities

Imported from `github.com/free5gc/util v1.3.1` into the
`github.com/acore2026/free6gc-lib` module. The upstream source packages and tests
are retained, together with the original Apache-2.0 license in `LICENSE.txt`.
The upstream module manifests are not retained as nested modules.

Imports of upstream util, NAS, and NGAP packages use the corresponding packages
inside `free6gc-lib`. Consumers import, for example:

```go
import (
    "github.com/acore2026/free6gc-lib/util/milenage"
    "github.com/acore2026/free6gc-lib/util/ueauth"
)
```

Run tests from the `free6gc-lib` root:

```sh
go test ./util/...
```

The `pfcp/` and `validator/` packages and the exported `milenage.F2345`
entry point are retained from `github.com/acore2026/util v1.3.3`.
Its UE-authentication implementation matches the upstream version.
