package connection

import (
	"fmt"
	"log"
	"net"
)

type Listener struct {
	typeListen        string
	listener          net.Listener
	initialized       bool
	host              string
	port              int
	connectionChannel func(conn *net.Conn)
}

func (l *Listener) Init(host string, port int, connectionChannel func(conn *net.Conn)) {
	log.Printf("Init server:%s:%d\n", host, port)
	l.initialized = true
	l.typeListen = "tcp"
	l.host = host
	l.port = port
	l.listen()
	l.connectionChannel = connectionChannel

}

// Run listens to the port of this service and will start the handler.
func (l *Listener) Run() {
	if !l.initialized {
		log.Println("Listening Service not initialized!")
	}

	for {
		conn, err := l.waitForConnection()

		if err == nil {
			l.connectionChannel(conn)
		} else {
			log.Fatalln(err.Error())
		}
	}
}

func (l *Listener) waitForConnection() (*net.Conn, error) {
	log.Println(l.typeListen)
	log.Println(l.host)
	log.Println(l.initialized)
	log.Println(l.port)

	connction, err := l.listener.Accept()
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	} else {
		log.Printf("Get connection on: %s:%d", l.host, l.port)
		return &connction, nil
	}
}

func (l *Listener) listen() {
	listener, err := net.Listen(l.typeListen, fmt.Sprintf("%s:%d", l.host, l.port))

	if err == nil && listener != nil {
		log.Printf("Listen on %s:%d\n", l.host, l.port)
		l.listener = listener
	} else if err != nil {
		log.Fatalln(err.Error())
		log.Fatalln("Maybe the port is not free?")
	} else if listener == nil {
		log.Fatalf("listener empty on host:%s and port %d\n", l.host, l.port)
	}
}
