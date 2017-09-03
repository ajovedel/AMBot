package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const Version = "v0.0.3-alpha"

var discordSession, _ = discordgo.New()

// Read in all configuration options from both environment variables and
// command line arguments.
func init() {

	// Discord Authentication Token
	discordSession.Token = os.Getenv("DG_TOKEN")
	if discordSession.Token == "" {
		flag.StringVar(&discordSession.Token, "t", "", "Discord Authentication Token")
	}
}

func main() {

	// Declare any variables needed later.
	var err error

	// Print out a fancy logo!
	fmt.Printf(` 
			________  .__                               .___
			\______ \ |__| ______ ____   ___________  __| _/
			||    |  \|  |/  ___// ___\ /  _ \_  __ \/ __ | 
			||    '   \  |\___ \/ /_/  >  <_> )  | \/ /_/ | 
			||______  /__/____  >___  / \____/|__|  \____ | 
			\_______\/        \/_____/   %-16s\/`+"\n\n", Version)

	// Parse command line arguments
	flag.Parse()

	// Verify a Token was provided
	if discordSession.Token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	// Verify the Token is valid and grab user information
	discordSession.State.User, err = discordSession.User("@me")
	if err != nil {
		log.Printf("error fetching user information, %s\n", err)
	}

	// Open a websocket connection to Discord
	err = discordSession.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
	}

	_, err = discordSession.ChannelMessageSend("348922901531066371", "Hello World")
	if err != nil {
		log.Printf("Derp, %s", err)
	}

	// load sounds
	err = loadSound()
	if err != nil {
		fmt.Println("Error loading sound: ", err)
		fmt.Println("Please copy $GOPATH/src/github.com/bwmarrin/examples/airhorn/airhorn.dca to this directory.")
		return
	}

	// set our fancy handlers
	setHandlers(discordSession)

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	discordSession.Close()
	ambotDB.Close()

	// Exit Normally.
}
