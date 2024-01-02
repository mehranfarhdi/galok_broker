package messagetopic

type AddTopic struct {
	MessageCode int      `json:"message_code"`
	Topics      []string `json:"topics"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
}

func NewAddTopic(messageCode int, Topics []string, username, password string) AddTopic {
	return AddTopic{
		MessageCode: messageCode,
		Topics:      Topics,
		Password:    password,
		Username:    username,
	}
}

func (m AddTopic) GetTopics() []string {
	return m.Topics
}

func (m AddTopic) GetMessageType() int {
	return m.MessageCode
}

func (m AddTopic) GetUsername() string {
	return m.Username
}

func (m AddTopic) GetPassword() string {
	return m.Password
}
