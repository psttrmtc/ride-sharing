package messaging

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ErrConnectionNotFound = errors.New("connection not found")
)

// connWrapper is a wrapper around the websocket connection to allow for thread-safe operations
// This is necessary because the websocket connection is not thread-safe

type connWrapper struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}
type ConnectionManager struct {
	connections map[string]*connWrapper // local connections
	mutex       sync.RWMutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
