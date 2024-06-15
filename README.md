# envconf

[![Build Status](https://travis-ci.org/devfans/envconf.svg?branch=master)](https://travis-ci.org/devfans/envconf)
[![Go Report Card](https://goreportcard.com/badge/github.com/devfans/envconf)](https://goreportcard.com/report/github.com/devfans/envconf)
[![GoDoc](https://godoc.org/github.com/devfans/envconf?status.svg)](https://godoc.org/github.com/devfans/envconf) [![Join the chat at https://gitter.im/devfans/envconf](https://badges.gitter.im/devfans/envconf.svg)](https://gitter.im/devfans/envconf?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Boostrap for operations with config file or env variables.
Read/Save config files like "~/.app" with sections and Set/Get env variables.

# Get Started

Config file sample:
```
[main]
name = test

[server]

[client]
```

Usages:
```
import "github.com/devfans/envconf"

func main() {
  config := envconf.NewConfig("~/.app")

  // get name key from config from default section: [main]
  name := config.Get("name")

  // get sections
  server := config.GetSection("server")
  client := config.GetSection("client")

  // add new keys
  server.Put("ip", "localhost")
  client.Put("ip", "0.0.0.0")

  // save to disk as file "~/.app"
  config.Save() 
  
  // other usages
  config.Section = "server"     // switch current section
  serverIp := config.String("ip")  // localhost

  config.Get("SERVER_IP", "ip") // get env first if env variable is not null

  config.Getenv("SERVER_IP")
  config.Setenv("SERVER_IP", "localhost")
}
```

## dotenv

```
# file: .env

[main]
use_section = case1

[case1]
name = a

[case2]
name = b

```

Example

```

import (
	"os"
	"testing"

	"github.com/devfans/envconf/dotenv"
)

func TestEnv(t *testing.T) {
	t.Log(os.Getenv("a"))
	t.Log(os.Getenv("b"))
	t.Log(dotenv.Int("a"))
	t.Log(dotenv.Uint("b"))
	t.Log(dotenv.Bool("c"))
	t.Log(dotenv.Bool("d"))
	t.Log(os.Getenv("test"))
	t.Log(dotenv.String("test"))
	t.Log(dotenv.EnvConf().Get("a"))
	t.Log(dotenv.EnvConf().Get("b"))
}

```
