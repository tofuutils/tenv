package misc

import "os"

const (
	RootEnv = "TENV_ROOT_PATH"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
