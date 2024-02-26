package dotenv

import (
	"os"

	"github.com/devfans/envconf"
)

// Config is the default envconf.Config instance parse from ".env" or file indicated by env name "ENV"
var Config *envconf.Config

func init() {
	envFile := ".env"
	if name := os.Getenv("ENV"); name != "" {
		envFile = name
	}
	if _, err := os.Stat(envFile); err != nil {
		Config = envconf.NewEmptyConfig()
		return
	}
	Config = envconf.NewConfig(envFile)
	for _, key := range Config.List() {
		os.Setenv(key, Config.GetConf(key))
	}
}

// String parse env value
//
// args set: (name)
// args set: (name, defaultValue)
//
func String(args... interface{}) string {
	return Config.String(args...)
}

// Int parse env value as int64
//
// args set: (name)
// args set: (name, defaultValue)
func Int(args... interface{}) int64 {
	return Config.Int(args...)
}

// Uint parse env value as uint64
//
// args set: (name)
// args set: (name, defaultValue)
func Uint(args... interface{}) uint64 {
	return Config.Uint(args...)
}

// Bool parse env value as bool
//
// args set: (name)
// args set: (name, defaultValue)
func Bool(args... interface{}) bool {
	return Config.Bool(args...)
}