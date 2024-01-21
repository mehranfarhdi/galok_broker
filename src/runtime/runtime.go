package runtime

import (
	"fmt"
	"github.com/hauke96/kingpin"
	"github.com/joho/godotenv"
	"github.com/mehranfarhdi/galok_broker/src/api/controllers"
	"github.com/mehranfarhdi/galok_broker/src/api/seed"
	"github.com/mehranfarhdi/galok_broker/src/conf"
	"github.com/mehranfarhdi/galok_broker/src/connection"
	"log"
	"net"
	"os"
)

var server = controllers.Server{}

const VERSION string = "v0.2.4"

var (
	app = kingpin.New("Galok", "A simple messaging service written in go")
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func RunRest() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	seed.Load(server.DB)

	server.Run(":8080")

}

// config for messaging

func configureCLI() {
	app.Author("Mehran Farhadi Bajestani")
	app.Version(VERSION)
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')
}

func startServer(config *conf.Config) {
	log.Println("Initialize services")

	listeningServices := initConnectionService(config)

	log.Println("Start connection listener")
	for _, listeningService := range listeningServices {
		go func(listeningService connection.Listener) {
			//TODO evaluate the need of a routine that restarts the service automatically when a error occurred. Something like: Error occurrec --> wait 5 seconds --> create service --> call Run()
			listeningService.Run()
		}(listeningService)
	}

	log.Println("\nThere we go, I'm ready to server ... eh ... serve\n")

	//TODO remove this and pass channels for closing
	select {}
}

// initConnectionService creates connection services bases on the given configuration.
func initConnectionService(config *conf.Config) []connection.Listener {
	log.Println("Initialize connection services")

	amountConnectors := len(config.ServerConfig.Connectors)

	listeningServices := make([]connection.Listener, amountConnectors)

	for i, connector := range config.ServerConfig.Connectors {
		// connection service
		connectionService := connection.Connector{}
		connectionService.Init(config.TopicList.Topics)

		// listening service
		newConnectionClosure := func(conn *net.Conn) {
			connectionService.HandleConnectionAsync(conn, config)
		}
		listeningService := connection.Listener{}
		listeningService.Init(connector.IpBind, connector.PortServe, newConnectionClosure)

		listeningServices[i] = listeningService
	}

	return listeningServices
}
