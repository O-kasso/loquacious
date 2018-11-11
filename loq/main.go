package main

import (
	"log"
	"os"
	"strconv"

	"github.com/o-kasso/loquacious/listen"
	"github.com/o-kasso/loquacious/talk"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: `loq [command] [args]`")
	}

	log.Print("Verifying presence of Google Cloud Platform credentials...")
	verifyGoogleCredentials()

	subcommand := os.Args[1]
	switch subcommand {
	case "talk":
		if os.Args[2] == "--demo" {
			talk.Demo()
		} else {
			talk.Talk(os.Args[2])
		}
	case "listen":
		timeLimit := getTimeLimit(os.Args[2])
		listen.Record(timeLimit)
	default:
		log.Println("Please use either the `talk` or `listen` subcommand.")
	}
}

func getTimeLimit(arg string) int {
	if i, err := strconv.Atoi(arg); err == nil {
		return i
	}

	log.Println("Using default recording time limit of 30 seconds")
	return 30
}

func verifyGoogleCredentials() {
	configFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	config, err := os.Stat(configFile)
	if err != nil || config.Size() <= 0 {
		log.Fatal("Could not find a valid Google Cloud Platform service key set to GOOGLE_APPLICATION_CREDENTIALS environment variable.")
	}
}
