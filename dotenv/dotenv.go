package dotenv

import (
	"os"

	"github.com/devfans/envconf"
)

// Config is the default envconf.Config instance parse from ".env" or file indicated by env name "ENV"
var config *envconf.Config

func init() {
	envFile := ".env"
	if name := os.Getenv("ENV"); name != "" {
		envFile = name
	}
	if _, err := os.Stat(envFile); err != nil {
		config = envconf.NewEmptyConfig()
		return
	}
	config = envconf.NewConfig(envFile)
	for _, key := range config.List() {
		os.Setenv(key, config.GetConf(key).String())
	}
}

// EnvConf return the global config instance
func EnvConf() *envconf.Config {
	return config
}

// String parse env value
//
// args set: (name)
// args set: (name, defaultValue)
//
func String(args... interface{}) string {
	return config.GetEnv(args...).String()
}

// Int parse env value as int64
//
// args set: (name)
// args set: (name, defaultValue)
func Int(args... interface{}) int64 {
	return config.GetEnv(args...).Int()
}

// Uint parse env value as uint64
//
// args set: (name)
// args set: (name, defaultValue)
func Uint(args... interface{}) uint64 {
	return config.GetEnv(args...).Uint()
}

// Bool parse env value as bool
//
// args set: (name)
// args set: (name, defaultValue)
func Bool(args... interface{}) bool {
	return config.GetEnv(args...).Bool()
}