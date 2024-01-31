package conf

import (
	"encoding/json"
	"github.com/mehranfarhdi/galok_broker/src/enviroment"
	"io/ioutil"
	"log"
	"strconv"
)

type TopicList struct {
	Topics []string `json:"topics"`
}

type DataBaseConf struct {
	Dbdriver   string
	DbUser     string
	DbPassword string
	DbPort     string
	DbHost     string
	DbName     string
}

type FiberConfig struct {
	PortServe string
}

type MessageBrokerConfig struct {
	ProtocolType string
	IpBind       string
	PortServe    int
}

type ServerConfig struct {
	DBConf     DataBaseConf
	Connectors []MessageBrokerConfig
	FiberConf  FiberConfig
	DebugServe bool
}

type Config struct {
	DBConf       DataBaseConf
	TopicList    TopicList
	ServerConfig ServerConfig
}

type ConfigLoader struct {
	dataBaseConf *DataBaseConf
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

	portRestStr, err := enviroment.Load(restPort_key)
	if err != nil {
		log.Fatalf("err to read %s from enviroment messagetopic: %s", restPort_key, err)
		return err
	}

	dbname, err := enviroment.Load(dbName)
	if err != nil {
		log.Fatalf("err to read %s from enviroment messagetopic: %s", dbName, err)
	}

	dbdriver, err := enviroment.Load(dbDriver)
	if err != nil {
		log.Fatalf("err to read %s from enviroment messagetopic: %s", dbDriver, err)
	}

	dbport, err := enviroment.Load(dbPort)
	if err != nil {
		log.Fatalf("err to read %s from enviroment messagetopic: %s", dbDriver, err)
	}

	dbhost, err := enviroment.Load(dbHost)
	if err != nil {
		log.Fatalf("err to read %s from enviroment messagetopic: %s", dbDriver, err)
	}

	dbpassword, err := enviroment.Load(dbPassword)
	if err != nil {
		log.Fatalf("err to read %s from enviroment messagetopic: %s", dbDriver, err)
	}

	dbuser, err := enviroment.Load(dbUser)
	if err != nil {
		log.Fatalf("err to read %s from enviroment messagetopic: %s", dbDriver, err)
	}

	log.Print("read config successful")

	cl.serverConfig = &ServerConfig{
		DBConf: DataBaseConf{
			Dbdriver:   dbdriver,
			DbName:     dbname,
			DbHost:     dbhost,
			DbPort:     dbport,
			DbUser:     dbuser,
			DbPassword: dbpassword,
		},
		Connectors: []MessageBrokerConfig{
			{
				ProtocolType: protocoltype,
				IpBind:       ipBind,
				PortServe:    portBind,
			},
		},
		FiberConf: FiberConfig{
			PortServe: portRestStr,
		},
		DebugServe: true,
	}

	//load topic conf
	cl.loadTopics("./topicConf/topics.json")

	return nil
}

func (cl *ConfigLoader) GetConfig() Config {
	return Config{
		DBConf:       *cl.dataBaseConf,
		TopicList:    *cl.topicList,
		ServerConfig: *cl.serverConfig,
	}
}
