package message

type Error struct {
	MessageType int    `json:"message_type"`
	ErrorCode   int    `json:"error_code"`
	Error       string `json:"error"`
}

func NewError(messageType int, errorCode int, error string) Error {
	return Error{MessageType: messageType, ErrorCode: errorCode, Error: error}
}

func (m Error) GetMessageType() int {
	return m.MessageType
}

func (e Error) GetErrorCode() int {
	return e.ErrorCode
}

func (e Error) GetError() string {
	return e.Error
}
