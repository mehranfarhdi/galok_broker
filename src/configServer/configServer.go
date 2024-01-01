package configServer

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type TopicList struct {
	Topics []string `json:"topics"`
}

type FiberConfig struct {
	IpBind 	  string
	PortServe int
}

type MessageBrokerConfig struct {
	ProtocolType string
	IpBind       string
	PortServe    int
}

type ServerConfig struct {
	Connectors []MessageBrokerConfig
	FiberConf  FiberConfig
	DebugServe bool
}

type Config struct {
	TopicList	 TopicList
	ServerConfig ServerConfig
}

type ConfigLoader struct {
	topicList 	 *TopicList
	serverConfig *ServerConfig
}

// LoadTopics reads the config file for the topics and fills the TopicConfig field.
// The default location is /topicConf/topics.json
func (cl *ConfigLoader) loadTopics(filename string) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatalf("Error reading %s", filename)
		log.Fatalln(err.Error())
	}

	cl.topicList = &TopicList{}
	json.Unmarshal(data, cl.serverConfig)
}


type LoadCofig