package websocket

import (
	"errors"
	"net/http"
	"sync"

	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/grpc/service/auth"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

// Manager is used to hold references to all Clients Registered, and Broadcasting etc
type Manager struct {
	clients ClientList

	// Using a syncMutex here to be able to lcok state before editing clients
	// Could also use Channels to block
	sync.RWMutex
	// handlers are functions that are used to handle Events
	handlers    map[string]EventHandler
	logger      *logger.Logger
	authService *auth.AuthService
}

// NewManager is used to initalize all the values inside the manager
func NewManager(logger *logger.Logger, authService *auth.AuthService) *Manager {
	m := &Manager{
		clients:     make(ClientList),
		handlers:    make(map[string]EventHandler),
		logger:      logger,
		authService: authService,
	}
	m.setupEventHandlers()
	return m
}

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessageHandler
	m.handlers[EventChangeRoom] = ChatRoomHandler
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (m *Manager) ServeWS(ctx *fiber.Ctx) {
	// Grab the access token and username from the query
	accessToken := ctx.Query("access_token")
	username := ctx.Query("username") // Extract username

	if accessToken == "" || username == "" {
		m.logger.Error("Access token or username not provided")
		ctx.Status(http.StatusUnauthorized).SendString("Unauthorized")
		return
	}

	// Verify the token
	_, err := m.authService.VerifyToken(ctx.Context(), domain.Auth{AccessToken: accessToken})
	if err != nil {
		// Tell the user its not authorized
		m.logger.Error(err.Error())
		ctx.Status(http.StatusUnauthorized).SendString("Unauthorized")
	}

	config := websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Origins:         []string{"https://localhost:3002"},
	}

	websocket.New(func(c *websocket.Conn) {
		// Ensure connection is closed only once
		var once sync.Once
		closeConnection := func() {
			once.Do(func() {
				m.logger.Info("Closing WebSocket connection")
				c.Close()
			})
		}

		// Create a new client
		client := NewClient(c, m, m.logger)

		// Add the client to the manager's list
		m.addClient(client)

		// Wait for both goroutines to finish
		var wg sync.WaitGroup
		wg.Add(2)

		// Start readMessages in a goroutine
		go func() {
			defer wg.Done()
			client.readMessages()
			closeConnection()
		}()

		// Start writeMessages in a goroutine
		go func() {
			defer wg.Done()
			client.writeMessages()
			closeConnection()
		}()

		// Wait for goroutines to finish
		wg.Wait()

		// Remove the client after goroutines finish
		m.removeClient(client)

	}, config)(ctx)
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}
