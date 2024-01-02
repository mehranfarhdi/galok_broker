package notification

import (
	"errors"
	msg "github.com/mehranfarhdi/galok_broker/src/messagetopic"
	"net"
	"sync"
)

type Notification struct {
	Connections *[]*net.Conn
	Topic       string
	Data        string
}

type Notifier struct {
	Queue       chan *Notification
	Errors      chan *Notification
	Exit        chan bool
	initialized bool
	mutex       *sync.Mutex
}

// Send String gives us to add \n to the byte data
func SendString(connection *net.Conn, data string) error {
	_, err := (*connection).Write([]byte(data + "\n"))
	return err
}

// Init creates all neccessary channel (queues) to handle notifications.
func (n *Notifier) Init() {
	n.Queue = make(chan *Notification)
	n.Errors = make(chan *Notification)
	n.Exit = make(chan bool)

	n.mutex = &sync.Mutex{}

	n.initialized = true
}

// StartNotifier listens to incoming notification requests.
func (n *Notifier) StartNotifier() error {
	if !n.initialized {
		return errors.New("TopicNotifyService not initialized")
	}

	for {
		select {
		case notification := <-n.Queue:
			go n.sendNotif(notification)
		case <-n.Exit:
			break
		}
	}
}

func (n *Notifier) sendNotif(notification *Notification) {
	message := msg.Message{
		MessageCode: msg.MT_MESSAGE,
		Topics:      []string{notification.Topic},
		Data:        notification.Data,
		Username:    "",
		Password:    "",
	}
}

// SendMessage enqueues the request to send the messagetopic into the queue.
// Therefore the sending itself will happen a bit later because the background
// thread will read from the queue.
func (n *Notifier) SendMessage(connections []*net.Conn, topic, message string) {
	// create notification
	notification := &Notification{
		Connections: &connections,
		Topic:       topic,
		Data:        message,
	}

	n.Queue <- notification
}
