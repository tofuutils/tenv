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

		return
	}

	err = tenv.DetectedCommandProxy(context.Background(), cmdconst.TofuName, "version")
	if err != nil {
		fmt.Println("proxy call failed :", err)
	}
}
