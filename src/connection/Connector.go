package connection

import (
	"github.com/jinzhu/gorm"
	"github.com/mehranfarhdi/galok_broker/src/conf"
	"github.com/mehranfarhdi/galok_broker/src/handler"
	"log"
	"net"
	"sync"
)

type Connector struct {
	topics            []string
	connectionHandler []*handler.Handler
	initialized       bool
	mutex             *sync.Mutex
	distributor       *handler.Distributor
}

// Init will initialize the connection service by creating all topic notifier and initializing fields.
func (c *Connector) Init(topics []string) {
	log.Println("Init connection service")

	c.distributor = &handler.Distributor{}
	c.distributor.Init(topics)

	c.topics = topics
	c.mutex = &sync.Mutex{}

	c.initialized = true
}

// HandleConnectionAsync creates a handler for the given connection and runs it in the background.
func (c *Connector) HandleConnectionAsync(conn *net.Conn, config *conf.Config, db *gorm.DB) {
	go c.createAndRunHandler(conn, config, db)
}

// createAndRunHandler sets up a new connection handler by registering to its events and starts it then.
// This should run on a new goroutine.
func (c *Connector) createAndRunHandler(conn *net.Conn, config *conf.Config, db *gorm.DB) {
	log.Println("Create connection handler")

	connHandler := handler.Handler{}
	connHandler.Init(conn, config, db)

	c.lock()
	c.distributor.Add(&connHandler)
	c.connectionHandler = append(c.connectionHandler, &connHandler)
	c.unlock()
	connHandler.HandleConnection()

	c.lock()

	// TODO move whole list of handlers and this removal into distributor
	// find connection handler index
	i := -1
	for j, a := range c.connectionHandler {
		if a == &connHandler {
			i = j
			break
		}
	}

	// remove connection handler
	if i != -1 {
		c.connectionHandler = append(c.connectionHandler[:i], c.connectionHandler[i+1:]...)
	}

	c.unlock()

	(*conn).Close()
}

// lock will prevent race conditions by ensuring that only one goroutine will have access to its fields.
func (c *Connector) lock() {
	c.mutex.Lock()
}

// unlock will free the fields so that other goroutines will have access to them.
func (c *Connector) unlock() {
	c.mutex.Unlock()
}
