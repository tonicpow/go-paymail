package main

import (
	"log"
	"time"

	"github.com/tonicpow/go-paymail/server"
)

func main() {

	// initialize the demo database
	if err := InitDemoDatabase(); err != nil {
		log.Fatal(err.Error())
	}

	// Custom server with lots of customizable goodies
	config, err := server.NewConfig(
		new(demoServiceProvider),
		server.WithBasicRoutes(),
		server.WithDomain("localhost"), // todo: make this work locally?
		server.WithDomain("another.com"),
		server.WithDomain("test.com"),
		server.WithGenericCapabilities(),
		server.WithPort(3000),
		server.WithServiceName("BsvAliasCustom"),
		server.WithTimeout(15*time.Second),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create & start the server
	server.Start(server.Create(config))
}
