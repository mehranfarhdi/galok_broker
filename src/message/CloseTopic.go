package message

type CloseTopic struct {
	MessageCode int    `json:"message_code"`
	UserName    string `json:"username"`
	Password    string `json:"password"`
}

func NewCloseTopic(messageCode int, username, password string) CloseTopic {
	return CloseTopic{
		MessageCode: messageCode,
		UserName:    username,
		Password:    password,
	}
}

func (m CloseTopic) GetMessageCode() int {
	return m.MessageCode
}

func (m CloseTopic) GetUsername() string {
	return m.UserName
}

func (m CloseTopic) GetPassword() string {
	return m.Password
}
