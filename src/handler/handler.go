package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"strings"

	"github.com/mehranfarhdi/galok_broker/src/conf"
	"github.com/mehranfarhdi/galok_broker/src/messagetopic"
	"github.com/mehranfarhdi/galok_broker/src/util"
)

type Handler struct {
	connection       *net.Conn
	connectionClosed bool
	config           *conf.Config
	registeredTopics []string
	SendEvent        []func(Handler, []string, string)
	ErrorEvent       []func(*Handler, int, string)
}

const MAX_PRINTING_LENGTH int = 80

// Init initializes the handler with the given connection.
func (h *Handler) Init(connection *net.Conn, config *conf.Config) {
	h.connection = connection
	h.config = config
}

// HandleConnection starts a routine to handle registration and sending messages.
// This will run until the client logs out, so run this in a goroutine.
func (h *Handler) HandleConnection() {
	// Not initialized
	if h.connection == nil {
		log.Fatalln("Connection not set!")
	}

	messageTypes := []int{messagetopic.MT_REGISTER,
		messagetopic.MT_LOGOUT,
		messagetopic.MT_CLOSE,
		messagetopic.MT_SEND}

	handler := []func(messagetopic.Message){h.handleRegistration,
		h.handleLogout,
		h.handleClose,
		h.handleSending}

	reader := bufio.NewReader(*h.connection)

	// Now a arbitrary amount of registration, logout, close and send messages is allowed
	for !h.connectionClosed {
		h.waitFor(
			messageTypes,
			handler,
			reader)
	}

	// TODO fire closed-event so that the distributor can react on that
}

// waitFor wats until on of the given message types arrived.
// The i-th argument in the messageTypes array must match to the i-th argument in the handler array.
func (h *Handler) waitFor(messageTypes []int, handler []func(message messagetopic.Message), reader *bufio.Reader) {

	// Check if the arrays match and error/fatal here
	if len(messageTypes) != len(handler) {
		if len(messageTypes) > len(handler) {
			// Fatal here to prevent a "slice bounds out of range" error during runtime
			log.Fatalln("There're more defined message types then functions mapped to them.")
		} else {
			log.Fatalln("There're more defined functions then message types here. Some message types might not be covered. Fix that!")
		}
	}

	rawMessage, err := reader.ReadString('\n')

	if err == nil {
		// the length of the message that should be printed
		maxOutputLength := int(math.Min(float64(len(rawMessage))-1, float64(MAX_PRINTING_LENGTH)))
		output := rawMessage[:maxOutputLength]
		if MAX_PRINTING_LENGTH < len(rawMessage)-1 {
			output += " [...]"
		}
		log.Println(output)

		// JSON to Message-struct
		message := getMessageFromJSON(rawMessage)

		// check type
		for i := 0; i < len(messageTypes); i++ {
			messageType := messageTypes[i]
			log.Printf("cheack if recived type %d is type %d", message.MessageCode, messageTypes)
			if message.MessageCode == messageType {
				log.Printf("Handle %d type", messageType)
				handler[i](message)
				break
			}
		}
	} else {
		log.Println("The connection will be closed. Reason: " + err.Error())
		h.exit()
		h.connectionClosed = true
	}
}

// getMessageFromJSON converts the given json-data into a message object.
func getMessageFromJSON(jsonData string) messagetopic.Message {
	message := messagetopic.Message{}
	json.Unmarshal([]byte(jsonData), &message)
	return message
}

// handleRegistration registeres this connection to the topics specified in the message.
func (h *Handler) handleRegistration(message messagetopic.Message) {
	log.Println("Try to register to topics " + fmt.Sprintf("%#v", message.Topics))

	// A comma separated list of all topics, the client is not allowed to register to
	forbiddenTopics := ""
	alreadyRegisteredTopics := ""

	for _, topic := range message.Topics {
		//TODO create a service for this. This should later take care of different user rights
		if !util.ContainsString(h.config.TopicList.Topics, topic) {
			forbiddenTopics += topic + ","
			log.Println("Clients wants to register on invalid topic (" + topic + ").")

		} else if util.ContainsString(h.registeredTopics, topic) {
			alreadyRegisteredTopics += topic + ","
			log.Println("Client already registered on " + topic)

		} else {
			h.registeredTopics = append(h.registeredTopics, topic)
			log.Println("Register " + topic)
		}
	}

	// Send error message for forbidden topics and cut trailing comma
	if len(forbiddenTopics) != 0 {
		forbiddenTopics = strings.TrimSuffix(forbiddenTopics, ",")

		for _, event := range h.ErrorEvent {
			event(h, messagetopic.ERR_REG_INVALID_TOPIC, forbiddenTopics)
		}
	}

	// Send error message for already registered topics and cut trailing comma
	if len(alreadyRegisteredTopics) != 0 {
		alreadyRegisteredTopics = strings.TrimSuffix(alreadyRegisteredTopics, ",")

		for _, event := range h.ErrorEvent {
			event(h, messagetopic.ERR_REG_ALREDY_REGISTERED, alreadyRegisteredTopics)
		}
	}
}

// handleSending send the given message to all clients interested in the topics specified in the message.
func (h *Handler) handleSending(message messagetopic.Message) {
	for _, event := range h.SendEvent {
		event(*h, message.Topics, message.Data)
	}
}

// handleLogout logs the client out.
func (h *Handler) handleLogout(message messagetopic.Message) {
	log.Println(fmt.Sprintf("Unsubscribe from topics %#v", message.Topics))
	h.logout(message.Topics)
}

// handleClose logs the client out from all topics and closes the connection.
func (h *Handler) handleClose(message messagetopic.Message) {
	h.exit()
}

// exit logs the client out from all topics and closes the connection.
func (h *Handler) exit() {
	log.Println("Unsubscribe from all topics")
	h.logout(h.registeredTopics)

	log.Println("Close connection")
	(*h.connection).Close()
	h.connectionClosed = true
}

// logout will logs the client out from the given topics.
func (h *Handler) logout(topics []string) {
	for _, topic := range topics {
		h.registeredTopics = util.RemoveString(h.registeredTopics, topic)
	}

	h.registeredTopics = util.RemoveStrings(h.registeredTopics, topics)
}

func (h *Handler) IsRegisteredTo(topic string) bool {
	return util.ContainsString(h.registeredTopics, topic)
}
