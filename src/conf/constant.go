package conf

// messagetopic broker config
var (
	ip_key           = "IP_BIND"
	protocolType_key = "PROTOCOL_TYPE"
	port_key         = "PORT_BIND"
	debug_key        = "DEBUG"
)

// rest api config
var (
	restPort_key  = "REST_PORT"
	topicFileName = "./topicConf/topics.json"
)

// data base config
var (
	dbDriver   = "DB_DRIVER"
	dbUser     = "DB_USER"
	dbPassword = "DB_PASSWORD"
	dbPort     = "DB_PORT"
	dbHost     = "DB_HOST"
	dbName     = "DB_NAME"
)
