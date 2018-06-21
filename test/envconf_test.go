package main

import (
	"envconf"
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

}
