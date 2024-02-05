package main

import (
	"github.com/devfans/envconf"
	"testing"
)

func TestConf(t *testing.T) {
	config := envconf.NewConfig("~/.test")
	server := config.GetSection("server")
	server.Put("address", "localhost")
	address := server.Get("address")
	if address != "localhost" {
		t.Errorf("Address does not match when getting from section")
	}

	address = server.GetConf("server_address", "0.0.0.0")
	if address != "0.0.0.0" {
		t.Errorf("Address does not match when getting from conf with default value")
	}

	config.Section = "server"
	address = config.Get("address")

	if address != "localhost" {
		t.Errorf("Address does not match, when getting from config")
	}

	config.Setenv("name-test", "name-test")
	name := config.Getenv("name-test")
	if name != "name-test" {
		t.Errorf("Name doest not match when use env as store")
	}

	address = config.GetEnv("server_address", "0.0.0.0")
	if address != "0.0.0.0" {
		t.Errorf("Address does not match when getting from env with default value")
	}
	config.Setenv("server_address", "10.0.0.0")
	address = config.GetEnv("server_address", "0.0.0.0")
	if address != "10.0.0.0" {
		t.Errorf("Address does not match when getting from env with default value")
	}

	valueStr := "\"https://xxfef\"//test"
	value := envconf.ParseValue(valueStr)
	if value != "https://xxfef" {
		t.Errorf("Failed to parse value, parsed %s", value)
	}
	t.Log(value, envconf.ParseValue("323//3234"))
}
