package messagetopic

type Message struct {
	MessageCode int      `json:"message_code"`
	Topics      []string `json:"topics"`
	Data        string   `json:"data"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
}

func NewMessage(messageCode int, topics []string, data, username, password string) Message {
	return Message{
		MessageCode: messageCode,
		Topics:      topics,
		Data:        data,
		Username:    username,
		Password:    password,
	}
}

func (m Message) GetMessageCode() int {
	return m.MessageCode
}

func (m Message) GetTopics() []string {
	return m.Topics
}

func (m Message) GetData() string {
	return m.Data
}

func (m Message) GetUsername() string {
	return m.Username
}

func (m Message) GetPassword() string {
	return m.Password
}
