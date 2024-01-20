package conf

import (
	"encoding/json"
	"github.com/mehranfarhdi/galok_broker/src/enviroment"
	"io/ioutil"
	"log"
	"strconv"
)

type FiberConfig struct {
	IpBind    string
	PortServe int
}

type MessageBrokerConfig struct {
	ProtocolType string
	IpBind       string
	PortServe    int
}

type ServerConfig struct {
	Connectors MessageBrokerConfig
	FiberConf  FiberConfig
	DebugServe bool
}

type Config struct {
	TopicList    TopicList
	ServerConfig ServerConfig
}

type ConfigLoader struct {
	topicList    *TopicList
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

func (cl *ConfigLoader) AddTopic(topicName string) {
	cl.topicList.Topics = append(cl.topicList.Topics, topicName)
}

// LoadConfig read the server config from env and call loadTopic file
func (cl *ConfigLoader) LoadConfig() error {
	//read ip bind in messagetopic broker
	ipBind, err := enviroment.Load(ip_key)
	if err != nil {
		log.Fatalf("err to read %s from enviroment messagetopic: %s", ip_key, err)
		return err
	}

	portBindStr, err := enviroment.Load(port_key)
	if err != nil {
		log.Fatalf("err to read %s from enviroment messagetopic: %s", port_key, err)
		return err
	}

	portBind, err := strconv.Atoi(portBindStr)
	if err != nil {
		// Handle the error if the conversion fails
		log.Fatalf("Error converting string to int: %s", err)
		return err
	}

	protocoltype, err := enviroment.Load(protocolType_key)
	if err != nil {
		log.Fatalf("err to read %s from enviroment messagetopic: %s", protocolType_key, err)
		return err
	}

	//debug, err := enviroment.Load(debug_key)
	//if err != nil {
	//	log.Fatalf("err to read %s from enviroment messagetopic: %s", debug_key, err)
	//	return err
	//}
	// read rest config from messagetopic broker
	ipRest, err := enviroment.Load(restIp_key)
	if err != nil {
		log.Fatalf("err to read %s from enviroment messagetopic: %s", restIp_key, err)
		return err
	}

	portRestStr, err := enviroment.Load(restPort_key)
	if err != nil {
		log.Fatalf("err to read %s from enviroment messagetopic: %s", restPort_key, err)
		return err
	}

	portRest, err := strconv.Atoi(portRestStr)
	if err != nil {
		// Handle the error if the conversion fails
		log.Fatalf("Error converting string to int: %s", err)
		return err
	}

	log.Print("read config successful")

	cl.serverConfig = &ServerConfig{
		Connectors: MessageBrokerConfig{
			ProtocolType: protocoltype,
			IpBind:       ipBind,
			PortServe:    portBind,
		},
		FiberConf: FiberConfig{
			IpBind:    ipRest,
			PortServe: portRest,
		},
		DebugServe: true,
	}

	//load topic conf
	cl.loadTopics("./topicConf/topics.json")

	return nil
}
