package server

import (
	"bytes"
	"github.com/op/go-logging"
	"github.com/google/uuid"
	"net"
	"launchomega.com/MYTH/manager"
)

var log = logging.MustGetLogger("")

func makeID() string {
	return uuid.Must(uuid.Parse(uuid.Must(uuid.NewRandom()).String())).String()
}

func StartServer() {
	listener, err := net.Listen("tcp", ":4208")
	if err != nil {
		log.Info(err)
	}
	clientManager := manager.ClientManager{
		Clients:		make(map[*manager.Client]bool),
		Broadcast:		make(chan *manager.MessageInfo),
		Register:		make(chan *manager.Client),
		Unregister:		make(chan *manager.Client),
	}
	go clientManager.Start()
	for {
		connection, _ := listener.Accept()
		if err != nil {
			log.Info(err)
		}
		uuidStr := makeID()
		log.Infof("New connection started with id: %+v", uuidStr)
		newBuffer := bytes.NewBuffer(make([]byte, 4096))
		client := &manager.Client{Socket: connection, Data: make(chan []byte), Id: uuidStr, Buffer: newBuffer}
		clientManager.Register <- client
		go clientManager.Receive(client)
		go clientManager.Send(client)
	}
}
