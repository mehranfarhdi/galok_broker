package message

type RemoveTopic struct {
	MessageType int      `json:"message_type"`
	Topics      []string `json:"topics"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
}

func NewRemoveTopic(messagetype int, topics []string, username, password string) RemoveTopic {
	return RemoveTopic{
		MessageType: messagetype,
		Topics:      topics,
		Username:    username,
		Password:    password,
	}
}

func (m RemoveTopic) GetMessageType() int {
	return m.MessageType
}

func (t RemoveTopic) GetTopics() []string {
	return t.Topics
}

func (t RemoveTopic) GetUsername() string {
	return t.Username
}

func (t RemoveTopic) GetPassword() string {
	return t.Password
}
