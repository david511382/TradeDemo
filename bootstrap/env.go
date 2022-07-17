package bootstrap

import (
	"os"
)

func GetEnvConfig() string {
	return os.Getenv("CONFIG")
}

func SetEnvConfig(s string) error {
	return os.Setenv("CONFIG", s)
}

func GetEnvPort() string {
	return os.Getenv("PORT")
}

func GetEnvWorkDir() string {
	return os.Getenv("WORK_DIR")
}

func SetEnvWorkDir(s string) error {
	return os.Setenv("WORK_DIR", s)
}
