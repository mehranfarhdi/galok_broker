package message

// send message from client to server
/*

The types of messages are:

register(topics, username, password) 	: Will register the client to a topic
send(topics, data, username, password)	: Sends the data to all subscribers of the topics
message(data, username, password)		: The receiving message a client gets
logout(topics, username, password)		: The clients unsubscribes from the given topics
close(topics, username, password)				: The client closes the connection and unsubscribes from all topics
*/

// type code from client to server (request)
const (
	MT_REGISTER = 101
	MT_SEND     = 102
	MT_LOGOUT   = 103
	MT_CLOSE    = 104
)

// type code from server to client ()
const (
	MT_MESSAGE = 205
	MT_ERROR   = 206
)

// Error codes
const (
	ERR_REG_INVALID_TOPIC     = 501
	ERR_REG_ALREDY_REGISTERED = 502
	ERR_SEND_FAILED           = 503
	ERR_USER_NOT_ALLOWED      = 504
	ERR_USER_NOT_FOUND        = 505
)