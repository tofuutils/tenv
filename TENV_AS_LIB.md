# How to use tenv as a library

## Get started

### Prerequisites

**tenv** requires [Go](https://go.dev) version [1.21](https://go.dev/doc/devel/release#go1.21.0) or above.

### Getting tenv module

```console
go get -u github.com/tofuutils/tenv

```

### Basic example

```go
package main

import (
    "context"
    "fmt"

    "github.com/tofuutils/tenv/v3/config/cmdconst"
    "github.com/tofuutils/tenv/v3/versionmanager/tenvlib"
)

func main() {
    tenv, err := tenvlib.Make(tenvlib.AutoInstall, tenvlib.IgnoreEnv, tenvlib.DisableDisplay)
    if err != nil {
        fmt.Println("init failed :", err)
    }

    err = tenv.DetectedCommandProxy(context.Background(), cmdconst.TofuName, "version")
    if err != nil {
        fmt.Println("proxy call failed :", err)
    }
}

```

Happy hacking !
