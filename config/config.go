package config

import (
	"os"
	"strings"
)

var Debug bool

func InitConfig() {
	debugEnv := strings.ToLower(os.Getenv("DEBUG"))
	if debugEnv == "true" || debugEnv == "1" {
		Debug = true
	} else {
		Debug = false
	}
}
