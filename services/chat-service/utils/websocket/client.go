package websocket

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/logger"
	"github.com/gofiber/websocket/v2"
)

var (
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontend visitor
type Client struct {
	// the websocket connection
	connection *websocket.Conn

	// manager is the manager used to manage the client
	manager *Manager
	// egress is used to avoid concurrent writes on the WebSocket
	egress chan Event
	// chatroom is used to know what room user is in
	chatroom string

	logger *logger.Logger
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager, logger *logger.Logger) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
		logger:     logger,
	}
}

// readMessages will start the client to read messages and handle them
// appropriatly.
// This is suppose to be ran as a goroutine
func (c *Client) readMessages() {
	defer func() {
		// Graceful Close the Connection once this
		// function is done
		c.manager.removeClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.logger.Error(fmt.Sprintf("error setting read deadline: %v", err))
		return
	}

	// Set the connection to read Pong messages
	c.connection.SetPongHandler(c.pongHandler)

	// Set the read limit to 512 bytes
	c.connection.SetReadLimit(512)

	// Loop Forever
	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error(fmt.Sprintf("error reading message: %v", err))
			}
			break // Break the loop to close conn & Cleanup
		}
		// Marshal incoming data into a Event struct
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			c.logger.Error(fmt.Sprintf("error unmarshalling message: %v", err))
			break // Breaking the connection here might be harsh xD
		}
		// Route the Event
		if err := c.manager.routeEvent(request, c); err != nil {
			c.logger.Error(fmt.Sprintf("error routing message: %v", err))
			break
		}
	}
}

// writeMessages is a process that listens for new messages to output to the Client
func (c *Client) writeMessages() {
	defer func() {
		// Graceful close if this triggers a closing
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			// Ok will be false Incase the egress channel is closed
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// Log that the connection is closed and the reason
					c.logger.Error(fmt.Sprintf("connection closed: %v", err))
				}
				// Return to close the goroutine
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				c.logger.Error(fmt.Sprintf("error marshalling message: %v", err))
				return
			}
			// Write a Regular text message to the connection
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				c.logger.Error(fmt.Sprintf("error writing message: %v", err))
				return
			}
			c.logger.Info("sent message")

		case <-ticker.C:
			c.logger.Info("ping")

			// Send a Ping to the Client
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.logger.Error(fmt.Sprintf("error writing ping: %v", err))
				return
			}
		}
	}
}

// pongHandler is used to handle Pong messages from the client
func (c *Client) pongHandler(pongMsg string) error {
	c.logger.Info("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
