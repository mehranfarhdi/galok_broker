package dist

import (
	"bufio"
	"encoding/json"
	"net"
	"testing"
	"time"

	testUtils "github.com/mehranfarhdi/galok_broker/src/util"
)

var conn1, conn2, serv1, serv2 *net.Conn
var buf1, buf2 *bufio.Reader
var serviceUnderTest *Notifier

const TEST_STRING string = "test123"

func initConnections(t *testing.T) {
	conn1, _, serv1, buf1 = testUtils.InitPipe()
	conn2, _, serv2, buf2 = testUtils.InitPipe()

	//Create connections
	//	listener := listen(t)
	//	conn1, buf1 = dial(t)
	//	conn2, buf2 = dial(t)

	//	return listener
}

func tearDownConnection() {
	//	l.Close()
	(*conn1).Close()
	(*conn2).Close()
}

func initNotifyService(t *testing.T) {
	serviceUnderTest = new(Notifier)
	serviceUnderTest.Init()
	go serviceUnderTest.StartNotifier()
}

func TestNotifyCorrectly(t *testing.T) {
	initNotifyService(t)
	//	l := initConnections(t)
	//	defer tearDownConnection(l)
	initConnections(t)

	connections := make([]*net.Conn, 2)
	connections[0] = conn1
	connections[1] = conn2

	notification := Notification{
		Connections: &connections,
		Data:        "test123\n",
	}

	serviceUnderTest.Queue <- &notification

	//
	// Test for client 1
	//
	received1, _, err := buf1.ReadLine()
	receivedObject1 := Notification{}
	json.Unmarshal(received1, &receivedObject1)
	if err != nil {
		t.Fail()
	}

	if notification.Data != receivedObject1.Data {
		t.Fail()
	}

	//
	// Test for client 2
	//
	received2, _, err := buf2.ReadLine()
	receivedObject2 := Notification{}
	json.Unmarshal(received2, &receivedObject2)
	if err != nil {
		t.Fail()
	}

	if notification.Data != receivedObject2.Data {
		t.Fail()
	}
}

func TestNotInitializedCreatesError(t *testing.T) {
	serviceUnderTest = new(Notifier)
	// This is missing: serviceUnderTest.Init()
	// There must be an error here:
	err := serviceUnderTest.StartNotifier()

	if err == nil {
		t.Fatal("The service should return an error.")
	}
}

func TestSendToExitChanWillExitCorrectly(t *testing.T) {
	serviceUnderTest = new(Notifier)
	serviceUnderTest.Init()

	go func(service *Notifier, t *testing.T) {
		err := service.StartNotifier()

		if err != nil {
			t.Fatal()
		}
	}(serviceUnderTest, t)

	// Do we need this?
	time.Sleep(time.Millisecond)

	serviceUnderTest.Exit <- true
}

func TestSendStringWorks(t *testing.T) {
	conn1, _, serv1, buf1 = testUtils.InitPipe()

	if buf1.Buffered() != 0 {
		t.Error("Buffered bytes must be 0")
		t.Fail()
	}

	go func(conn *net.Conn) { SendStringTo(conn, TEST_STRING) }(conn1)

	data, _, err := buf1.ReadLine()

	if err != nil {
		t.Fail()
	}

	if string(data) != TEST_STRING {
		t.Error("Buffered bytes must be !=0")
		t.Fail()
	}
}
