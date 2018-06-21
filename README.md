# envconf

[![Build Status](https://travis-ci.org/devfans/envconf.svg?branch=master)](https://travis-ci.org/devfans/envconf)
[![Go Report Card](https://goreportcard.com/badge/github.com/devfans/envconf)](https://goreportcard.com/report/github.com/devfans/envconf)

Boostrap for operations with config file or env variables.
Read/Save config files like "~/.app" with sections and Set/Get env variables.

# Get Started

```
import "github.com/devfans/envconf"

func main() {
  config := envconf.NewConfig("~/.app")
  config.Put("name", "my-app")
  server := config.GetSection("server")
  client := config.GetSection("client")
  server.Put("ip", "localhost")
  client.Put("ip", "0.0.0.0")

  config.Save() 
  config.Section = "server"     // switch current section
  serverIp := config.Get("ip")  // localhost

  config.Get("SERVER_IP", "ip") // get env first if env variable is not null

  config.Getenv("SERVER_IP")
  config.Setenv("SERVER_IP", "localhost")
}
```
```
~/.app:
[main]
name = my-app

[server]
ip = localhost

[client]
ip = 0.0.0.0
```
