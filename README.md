# envconf

[![Build Status](https://travis-ci.org/devfans/envconf.svg?branch=master)](https://travis-ci.org/devfans/envconf)
[![Go Report Card](https://goreportcard.com/badge/github.com/devfans/envconf)](https://goreportcard.com/report/github.com/devfans/envconf)
[![GoDoc](https://godoc.org/github.com/devfans/envconf?status.svg)](https://godoc.org/github.com/devfans/envconf)

Boostrap for operations with config file or env variables.
Read/Save config files like "~/.app" with sections and Set/Get env variables.

# Get Started

```
import "github.com/devfans/envconf"

func main() {
  config := envconf.NewConfig("~/.app")

  // get name from config
  name := config.Get("name")

  // get sections
  server := config.GetSection("server")
  client := config.GetSection("client")

  // add new keys
  server.Put("ip", "localhost")
  client.Put("ip", "0.0.0.0")

  // save to disk as file "~/.app"
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
