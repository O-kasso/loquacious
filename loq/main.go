package main

import (
	"github.com/o-kasso/loquacious/listen"
	"github.com/o-kasso/loquacious/talk"
	"log"
	"os"
)

// expects path to valid SSML file as argument.
func main() {
	// TODO: validate that second arg is valid SSML
	if len(os.Args) <= 2 {
		log.Fatal("Usage: `loq talk hello-world.ssml`")
	}

	log.Print("Verifying presence of Google Cloud Platform credentials...")
	verifyGoogleCredentials()

	subcommand := os.Args[1]
	arg := os.Args[2]
	switch subcommand {
	case "talk":
		talk.Talk(arg)
	case "listen":
		listen.Record(arg)
	default:
		log.Println("Please use either the `talk` or `listen` subcommand.")
	}
}

func verifyGoogleCredentials() {
	configFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	config, err := os.Stat(configFile)
	if err != nil || config.Size() <= 0 {
		log.Fatal("Could not find a valid Google Cloud Platform service key set to GOOGLE_APPLICATION_CREDENTIALS environment variable.")
	}
}
