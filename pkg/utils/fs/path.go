package fs

import (
	"fmt"
	"github.com/opentofuutils/tenv/pkg/utils/env"
	log "github.com/sirupsen/logrus"
)

func GetPath(name string) string {
	rootDir := env.GetEnv(env.RootEnv, "")

	switch name {
	case "root_dir":
		return rootDir
	case "bin_dir":
		return fmt.Sprintf("%s/bin", rootDir)
	case "misc_dir":
		return fmt.Sprintf("%s/misc", rootDir)
	case "tfenv_dir":
		return fmt.Sprintf("%s/bin/tfenv", rootDir)
	case "tofuenv_dir":
		return fmt.Sprintf("%s/bin/tofuenv", rootDir)
	case "tfenv_exec":
		return fmt.Sprintf("%s/bin/tfenv/bin/tfenv", rootDir)
	case "tofuenv_exec":
		return fmt.Sprintf("%s/bin/tofuenv/bin/tofuenv", rootDir)

	default:
		log.Warn("Unknown day")
		return ""
	}

}
