package dist

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/mehranfarhdi/galok_broker/src/messagetopic"
)

// SendStringTo sends the given string with an \n character to the given connection.
func SendStringTo(connection *net.Conn, data string) error {
	_, err := (*connection).Write([]byte(data + "\n"))

	return err
}

type Notifier struct {
	Queue       chan *Notification
	Errors      chan *Notification
	Exit        chan bool
	initialized bool
	mutex       *sync.Mutex
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
			go n.sendNotification(notification)
		case <-n.Exit:
			break
		}
	}
}

// SendMessage enqueues the request to send the message into the queue.
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

// sendNotification sends the notification or an error if there's one.
func (n *Notifier) sendNotification(notification *Notification) {
	message := messagetopic.Message{
		MessageCode: messagetopic.MT_MESSAGE,
		Topics:      []string{notification.Topic},
		Data:        notification.Data,
	}

	if len(notification.Data) > 10 {
		log.Print("send message with data: " + notification.Data[0:10] + "[...]\n")
	} else {
		log.Print("send message with data: " + notification.Data)
	}

	messageByteArray, err := json.Marshal(message)
	messageString := string(messageByteArray)

	if err != nil {
		log.Fatal("Error parsing message data: " + err.Error())

		for _, connection := range *notification.Connections {
			n.SendError(connection, messagetopic.ERR_SEND_FAILED, err.Error())
		}
		return
	}

	for _, connection := range *notification.Connections {
		err := SendStringTo(connection, messageString)

		if err != nil {
			log.Fatalln(fmt.Sprintf("Could not send message to %s", (*connection).RemoteAddr()))
		}
	}
}

// TODO also create a queue for the errors
// SendError directly sends the error, there's no asynchronous queue here.
func (n *Notifier) SendError(connection *net.Conn, errorCode int, message string) {
	errorMessage := messagetopic.Error{
		MessageCode: messagetopic.MT_ERROR,
		ErrorCode:   errorCode,
		Error:       message,
	}

	data, err := json.Marshal(errorMessage)

	if err == nil {
		log.Fatalln("Sending error")
		SendStringTo(connection, string(data)+"\n")
	} else {
		log.Fatalln("Error while sending error: " + err.Error())
	}

}
