package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/seme4eg/sadbot/session"
	"github.com/seme4eg/sadbot/utils"
)

func main() {
	// parse config file
	config, err := utils.GetConfig()
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err)
	}

	// Create new Discord Session
	session, err := session.OpenSession("Bot "+config.Token, config.Prefix)
	if err != nil {
		log.Fatalf("Failed to create a discord session: %s", err)
	}

	// ensure that session will be gracefully closed whenever the function exits
	defer session.Close()

	// run until code is terminated
	fmt.Println("sadbot is now running. Press Ctrl-C to exit.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-c
}
