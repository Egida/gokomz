package servers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/quietjoy/gocom/pkg/db"
	"github.com/quietjoy/gocom/pkg/models"
	"github.com/quietjoy/gocom/pkg/utils"
)

func RunServer() {
	// Initialize the DB
	dsn := db.GenerateMysqlDSN("mysql", "3306", "root", "password", "gocom")
	mysql, err := db.NewMysqlDB(dsn)
	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}
	log.Println("Successfully connected to DB")

	// Migrate the models
	err = db.Migrate(mysql, &models.Client{}, &models.ControlCommand{})
	if err != nil {
		log.Fatal("Error migrating models: ", err)
	}
	log.Println("Successfully migrated models")

	// Start the server
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		// Return heartbeat
		return c.JSON(http.StatusOK, "OK")
	})

	// Get a list of all registered clients
	e.GET("/admin/clients", func(c echo.Context) error {
		log.Println("Getting list of clients")
		// Get the list of clients from the DB
		clients, err := db.GetClients(mysql)
		if err != nil {
			log.Println(fmt.Sprintf("Error getting clients: %s", err))
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, clients)
	})

	// Get commands for a specific client
	e.GET("/:uuid/commands", func(c echo.Context) error {
		clientUUID, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			log.Println(fmt.Sprintf("Error parsing client UUID: %s", err))
			return c.JSON(http.StatusBadRequest, nil)
		}

		// Get the client by UUID from the database
		client, err := db.GetClientByUUID(mysql, clientUUID.String())
		if err != nil {
			// check if error was because client was not found
			if err == gorm.ErrRecordNotFound {
				log.Println(fmt.Sprintf("Client not found: %s", err))
				return c.JSON(http.StatusNotFound, nil)
			}

			log.Println(fmt.Sprintf("Error getting client by UUID: %s", err))
			return c.JSON(http.StatusInternalServerError, nil)
		}

		if client.SourceIP != utils.GetRemoteIPFromRemoteAddr(c.Request().RemoteAddr) {
			// For now just log if this happens
			log.Println(fmt.Sprintf("Client IP mismatch. Registered: %s, Actual: %s", client.SourceIP, c.Request().RemoteAddr))
			// TODO: Handle more gracefully
			// return c.JSON(http.StatusUnauthorized, nil)
		}

		// Get the commands for the client from the database
		commands, err := db.GetCommands(mysql, client.ID)
		if err != nil {
			log.Println(fmt.Sprintf("Error getting commands for client %s: %s", clientUUID, err))
			return c.JSON(http.StatusInternalServerError, nil)
		}

		// If there are no commands, return 204
		if len(commands) == 0 {
			log.Println(fmt.Sprintf("No commands for client %s", clientUUID))
			return c.JSON(http.StatusNoContent, nil)
		}

		// return the commands for the Client
		return c.JSON(http.StatusOK, commands)
	})

	// Register a command for a specific client
	e.POST("/:uuid/command", func(c echo.Context) error {
		clientUUID, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			log.Println(fmt.Sprintf("Error parsing client ID: %s", err))
			return c.JSON(http.StatusBadRequest, nil)
		}

		// Get the client by UUID from the database
		client, err := db.GetClientByUUID(mysql, clientUUID.String())
		if err != nil {
			// check if error was because client was not found
			if err == gorm.ErrRecordNotFound {
				log.Println(fmt.Sprintf("Client not found: %s", err))
				return c.JSON(http.StatusNotFound, nil)
			}
			log.Println(fmt.Sprintf("Error getting client by UUID: %s", err))
			return c.JSON(http.StatusInternalServerError, nil)
		}

		// Check source IP against what was registered
		if client.SourceIP != utils.GetRemoteIPFromRemoteAddr(c.Request().RemoteAddr) {
			// For now just log if this happens
			log.Println(fmt.Sprintf("Client IP mismatch. Registered: %s, Actual: %s", client.SourceIP, c.Request().RemoteAddr))
		}

		// Get the command from the request body
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			log.Fatal(err)
		}

		// The command is a json encoded string, so we need to decode it
		var controlCommand models.ControlCommand
		json.Unmarshal(bodyBytes, &controlCommand)
		controlCommand.ClientID = client.ID
		log.Println(fmt.Sprintf("Registering command for client %s: %s %s", clientUUID, controlCommand.Command, controlCommand.Arguments))

		// Save command to the database
		err = db.SaveCommand(mysql, controlCommand)
		return c.JSON(http.StatusOK, nil)
	})

	// Get information about a specific client
	e.GET("/:uuid/info", func(c echo.Context) error {
		clientUUID, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			log.Println(fmt.Sprintf("Error parsing client UUID: %s", err))
			return c.JSON(http.StatusBadRequest, nil)
		}

		// Get the client by UUID from the database
		client, err := db.GetClientByUUID(mysql, clientUUID.String())
		if err != nil {
			// check if error was because client was not found
			if err == gorm.ErrRecordNotFound {
				log.Println(fmt.Sprintf("Client not found: %s", err))
				return c.JSON(http.StatusNotFound, nil)
			}
			log.Println(fmt.Sprintf("Error getting client by UUID: %s", err))
			return c.JSON(http.StatusInternalServerError, nil)
		}

		// Check source IP against what was registered
		if client.SourceIP != utils.GetRemoteIPFromRemoteAddr(c.Request().RemoteAddr) {
			// For now just log if this happens
			log.Println(fmt.Sprintf("Client IP mismatch. Registered: %s, Actual: %s", client.SourceIP, c.Request().RemoteAddr))
		}

		return c.JSON(http.StatusOK, client)
	})

	e.POST("/register", func(c echo.Context) error {
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)

		// create a new client
		client := models.Client{
			UUID:        uuid.New(),
			SourceIP:    utils.GetRemoteIPFromRemoteAddr(c.Request().RemoteAddr),
			ForwardedIP: c.Request().Header["X-Forwarded-For"],
		}

		// save client to the database and get the ID
		mysql.Create(&client)

		return c.JSON(http.StatusOK, client.UUID)
	})

	e.POST("/:uuid/:commandID/error", func(c echo.Context) error {
		clientUUID, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			log.Println(fmt.Sprintf("Error parsing client UUID: %s", err))
			return c.JSON(http.StatusBadRequest, nil)
		}

		// Get the client by UUID from the database
		client, err := db.GetClientByUUID(mysql, clientUUID.String())
		if err != nil {
			// check if error was because client was not found
			if err == gorm.ErrRecordNotFound {
				log.Println(fmt.Sprintf("Client not found: %s", err))
				return c.JSON(http.StatusNotFound, nil)
			}
			log.Println(fmt.Sprintf("Error getting client by UUID: %s", err))
			return c.JSON(http.StatusInternalServerError, nil)
		}

		// Check source IP against what was registered
		if client.SourceIP != utils.GetRemoteIPFromRemoteAddr(c.Request().RemoteAddr) {
			// For now just log if this happens
			log.Println(fmt.Sprintf("Client IP mismatch. Registered: %s, Actual: %s", client.SourceIP, c.Request().RemoteAddr))
		}

		// Get the command from the request body
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			log.Fatal(err)
		}
		clientError := string(bodyBytes)
		log.Println(fmt.Sprintf("Error from client %s: %s", clientUUID.String(), clientError))

		// Convert commandID to uint
		commandID, err := strconv.ParseUint(c.Param("commandID"), 10, 64)
		if err != nil {
			log.Println(fmt.Sprintf("Error parsing command ID: %s", err))
			return c.JSON(http.StatusBadRequest, nil)
		}
		// Get command from the database by ID
		command, err := db.GetCommandByID(mysql, uint(commandID))
		if err != nil {
			// check if error was because command was not found
			if err == gorm.ErrRecordNotFound {
				log.Println(fmt.Sprintf("Command not found: %s", err))
				return c.JSON(http.StatusNotFound, nil)
			}
			log.Println(fmt.Sprintf("Error getting command by ID: %s", err))
			return c.JSON(http.StatusInternalServerError, nil)
		}

		// Update command status to error
		command.Status = "ERROR"
		command.Output = clientError
		err = db.Update(mysql, command)

		return c.JSON(http.StatusOK, nil)
	})

	e.POST("/:uuid/:commandID/success", func(c echo.Context) error {
		clientUUID, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			log.Println(fmt.Sprintf("Error parsing client UUID: %s", err))
			return c.JSON(http.StatusBadRequest, nil)
		}

		// Get the client by UUID from the database
		client, err := db.GetClientByUUID(mysql, clientUUID.String())
		if err != nil {
			// check if error was because client was not found
			if err == gorm.ErrRecordNotFound {
				log.Println(fmt.Sprintf("Client not found: %s", err))
				return c.JSON(http.StatusNotFound, nil)
			}
			log.Println(fmt.Sprintf("Error getting client by UUID: %s", err))
			return c.JSON(http.StatusInternalServerError, nil)
		}

		// Check source IP against what was registered
		if client.SourceIP != utils.GetRemoteIPFromRemoteAddr(c.Request().RemoteAddr) {
			// For now just log if this happens
			log.Println(fmt.Sprintf("Client IP mismatch. Registered: %s, Actual: %s", client.SourceIP, c.Request().RemoteAddr))
		}

		// Get the command from the request body
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			log.Fatal(err)
		}
		clientOutput := string(bodyBytes)
		log.Println(fmt.Sprintf("Response from client %s: %s", clientUUID.String(), clientOutput))

		// Convert commandID to uint
		commandID, err := strconv.ParseUint(c.Param("commandID"), 10, 64)
		if err != nil {
			log.Println(fmt.Sprintf("Error parsing command ID: %s", err))
			return c.JSON(http.StatusBadRequest, nil)
		}

		// Get command from the database by ID
		command, err := db.GetCommandByID(mysql, uint(commandID))
		if err != nil {
			// check if error was because command was not found
			if err == gorm.ErrRecordNotFound {
				log.Println(fmt.Sprintf("Command not found: %s", err))
				return c.JSON(http.StatusNotFound, nil)
			}
			log.Println(fmt.Sprintf("Error getting command by ID: %s", err))
			return c.JSON(http.StatusInternalServerError, nil)
		}

		// Update command status to error
		command.Status = "SUCCESS"
		command.Output = clientOutput
		err = db.Update(mysql, command)
		log.Println(fmt.Sprintf("Updated command %d to status %s", command.ID, command.Status))

		return c.JSON(http.StatusOK, nil)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
