package message

type CloseConnection struct {
	MessageCode int    `json:"message_code"`
	UserName    string `json:"username"`
	Password    string `json:"password"`
}

func NewCloseTopic(messageCode int, username, password string) CloseConnection {
	return CloseConnection{
		MessageCode: messageCode,
		UserName:    username,
		Password:    password,
	}
}

func (m CloseConnection) GetMessageCode() int {
	return m.MessageCode
}

func (m CloseConnection) GetUsername() string {
	return m.UserName
}

func (m CloseConnection) GetPassword() string {
	return m.Password
}
