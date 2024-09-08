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
	conf, err := config.DefaultConfig() // does not read environment variables
	if err != nil {
		fmt.Println("init failed :", err)

		return
	}

	conf.SkipInstall = false // tenvlib.AutoInstall option equivalent

	tenv, err := tenvlib.Make(tenvlib.WithConfig(&conf), tenvlib.DisableDisplay)
	if err != nil {
		fmt.Println("should not occur when calling WithConfig :", err)

		return
	}

	ctx := context.Background()
	version, err := tenv.Evaluate(ctx, cmdconst.TerraformName, semantic.LatestKey)
	if err != nil {
		fmt.Println("eval failed :", err)

		return
	}

	conf.ForceRemote = true

	remoteVersion, err := tenv.Evaluate(ctx, cmdconst.TerraformName, semantic.LatestKey)
	if err != nil {
		fmt.Println("eval remote failed :", err)

		return
	}

	if version != remoteVersion {
		err = tenv.Uninstall(ctx, cmdconst.TerraformName, version)
		if err != nil {
			fmt.Println("uninstall failed :", err)
		}
	}

	fmt.Println("Last Terraform version :", version, "(local),", remoteVersion, "(remote)")
}
