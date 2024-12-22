# How to use tenv as a library

## Get started

### Prerequisites

**tenv** requires [Go](https://go.dev) version [1.23](https://go.dev/doc/devel/release#go1.23.0) or above.

### Getting tenv module

`tenvlib` package is available since tenv v3.2

```console
go get -u github.com/tofuutils/tenv/v4@latest
```

### Basic example

```go
package main

import (
    "context"
    "fmt"

    "github.com/tofuutils/tenv/v4/config/cmdconst"
    "github.com/tofuutils/tenv/v4/versionmanager/tenvlib"
)

func main() {
    tenv, err := tenvlib.Make(tenvlib.AutoInstall, tenvlib.IgnoreEnv, tenvlib.DisableDisplay)
    if err != nil {
        fmt.Println("init failed :", err)
        return
    }

    err = tenv.DetectedCommandProxy(context.Background(), cmdconst.TofuName, "version")
    if err != nil {
        fmt.Println("proxy call failed :", err)
    }
}
```

## Documentation

See the [API documentation on go.dev](https://pkg.go.dev/github.com/tofuutils/tenv/v4/versionmanager/tenvlib) and [examples](https://github.com/tofuutils/tenv/tree/main/versionmanager/tenvlib/examples).

### Overview

Available Tenv struct creation options :

- `AddTool(toolName string, builderFunc builder.BuilderFunc)`, extend `tenvlib` to support other tool use cases.
- `AutoInstall`, shortcut to force auto install feature enabling in `config.Config`.
- `DisableDisplay`, do not display or log anything.
- `IgnoreEnv`, ignore **tenv** environment variables (`TENV_AUTO_INSTALL`, `TOFUENV_TOFU_VERSION`, etc.).
- `WithConfig(conf *config.Config)`, replace default `Config` (one from a `InitConfigFromEnv` or `DefaultConfig` call depending on `IgnoreEnv` usage).
- `WithDisplayer(displayer loghelper.Displayer)`, replace default `Displayer` with a custom to handle `tenvlib` output (standard and log).
- `WithHCLParser(hclParser *hclparse.Parser)`, use passed `Parser` instead of creating a new one.

Tenv methods list :

- `[Detected]Command[Proxy]`
- `Detect`
- `Evaluate`
- `Install[Multiple]`
- `List[Local|Remote]`
- `LocallyInstalled`
- `[Res|S]etDefault[Constraint|Version]`, manage `constraint` and `version` files in `<rootPath>/<tool>/`
- `Uninstall[Multiple]`

Happy hacking !
