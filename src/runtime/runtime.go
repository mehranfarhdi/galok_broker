package runtime

import (
	"fmt"
	"github.com/hauke96/kingpin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/mehranfarhdi/galok_broker/src/api/controllers"
	"github.com/mehranfarhdi/galok_broker/src/api/seed"
	"github.com/mehranfarhdi/galok_broker/src/conf"
	"github.com/mehranfarhdi/galok_broker/src/connection"
	"log"
	"net"
	"sync"
)

const VERSION string = "v0.2.4"

var server = controllers.Server{}

var (
	app = kingpin.New("Galok", "A simple messaging service written in go")
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func RunRest(db *gorm.DB, port string) {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(db)

	seed.Load(server.DB)

	server.Run(fmt.Sprintf(":%s", port))

}

// config for messaging

func configureCLI() {
	app.Author("Mehran Farhadi Bajestani")
	app.Version(VERSION)
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')
}

func startServer(config *conf.Config, db *gorm.DB) {
	log.Println("Initialize services")

	listeningServices := initConnectionService(config, db)

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
func initConnectionService(config *conf.Config, db *gorm.DB) []connection.Listener {
	log.Println("Initialize connection services")

	amountConnectors := len(config.ServerConfig.Connectors)

	log.Println("how match:", amountConnectors)

	listeningServices := make([]connection.Listener, amountConnectors)

	for i, connector := range config.ServerConfig.Connectors {
		// connection service
		connectionService := connection.Connector{}
		connectionService.Init(config.TopicList.Topics)

		// listening service
		newConnectionClosure := func(conn *net.Conn) {
			connectionService.HandleConnectionAsync(conn, config, db)
		}
		listeningService := connection.Listener{}
		listeningService.Init(connector.IpBind, connector.PortServe, newConnectionClosure)

		listeningServices[i] = listeningService
	}

	return listeningServices
}

func RunDataBaseAndConf() {
	// get config for data base
	configLoader := conf.ConfigLoader{}
	configLoader.LoadConfig()
	config := configLoader.GetConfig()
	dataBase := conf.InitializeDB(config.DBConf.Dbdriver, config.DBConf.DbUser, config.DBConf.DbPassword, config.DBConf.DbPort, config.DBConf.DbHost, config.DBConf.DbName)

	// make channels to signal when each function has completed

	var wg sync.WaitGroup

	wg.Add(2)
	go RunRest(dataBase, "8080")

	go startServer(&config, dataBase)

	fmt.Println("Waiting for goroutines to finish...")
	wg.Wait()
	fmt.Println("Done!")

}
