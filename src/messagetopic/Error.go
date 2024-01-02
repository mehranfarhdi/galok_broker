package messagetopic

type Error struct {
	MessageCode int    `json:"message_code"`
	ErrorCode   int    `json:"error_code"`
	Error       string `json:"error"`
}

func NewError(messageType int, errorCode int, error string) Error {
	return Error{MessageCode: messageType, ErrorCode: errorCode, Error: error}
}

func (m Error) GetMessageType() int {
	return m.MessageCode
}

func (e Error) GetErrorCode() int {
	return e.ErrorCode
}

func (e Error) GetError() string {
	return e.Error
}
