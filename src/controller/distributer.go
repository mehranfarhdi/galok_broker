package controller

import (
	"github.com/mehranfarhdi/galok_broker/src/conf"
	"net"
)

type Controller struct {
	connection *net.Conn
	isClose    bool
	config     *conf.Config
	regTopics  []string
	ErrorEvent []func(*Controller, []string, string)
	SendEvent  []func(Controller, []string, string)
}

const MAX_PRINTING_LENGTH int = 80

//Controle
