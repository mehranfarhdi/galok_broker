package message

type Abstract struct {
	MessageCode int    `json:"message_code"`
	UserName    string `json:"username"`
	Password    string `json:"password"`
}

func NewAbstract(messageCode int, username, password string) Abstract {
	return Abstract{
		MessageCode: messageCode,
		UserName:    username,
		Password:    password,
	}
}

func (m Abstract) GetMessageCode() int {
	return m.MessageCode
}

func (m Abstract) GetUsername() string {
	return m.UserName
}

func (m Abstract) GetPassword() string {
	return m.Password
}
