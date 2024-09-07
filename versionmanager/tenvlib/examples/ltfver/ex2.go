package main

import (
	"context"
	"fmt"

	"github.com/tofuutils/tenv/v3/config"
	"github.com/tofuutils/tenv/v3/config/cmdconst"
	"github.com/tofuutils/tenv/v3/versionmanager/semantic"
	"github.com/tofuutils/tenv/v3/versionmanager/tenvlib"
)

func main() {
	conf, err := config.DefaultConfig()
	if err != nil {
		fmt.Println("init failed :", err)
	}

	conf.SkipInstall = false

	tenv, err := tenvlib.Make(tenvlib.WithConfig(&conf), tenvlib.DisableDisplay)
	if err != nil {
		fmt.Println("should not occur when calling WithConfig :", err)
	}

	ctx := context.Background()
	version, err := tenv.Evaluate(ctx, cmdconst.TerraformName, semantic.LatestKey)
	if err != nil {
		fmt.Println("eval failed :", err)
	}

	err = tenv.Uninstall(ctx, cmdconst.TerraformName, version)
	if err != nil {
		fmt.Println("uninstall failed :", err)
	}

	conf.ForceRemote = true

	remoteVersion, err := tenv.Evaluate(ctx, cmdconst.TerraformName, semantic.LatestKey)
	if err != nil {
		fmt.Println("eval remote failed :", err)
	}

	fmt.Println("Last Terraform version :", version, "(local),", remoteVersion, "(remote)")
}
