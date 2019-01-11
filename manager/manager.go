package manager

import (
	"bytes"
	"github.com/op/go-logging"
	"net"
)

type MessageInfo struct {
	Message			[]byte
	Id				string
}

type ClientManager struct {
	Clients			map[*Client]bool
	Broadcast		chan *MessageInfo
	Register		chan *Client
	Unregister		chan *Client
}

type Client struct {
	Socket			net.Conn
	Data 			chan []byte
	Id				string
	Buffer			*bytes.Buffer
}

var log = logging.MustGetLogger("")

func (manager *ClientManager) Start() {
	for {
		select {
		case connection := <-manager.Register:
			manager.Clients[connection] = true
			log.Info("A connection was established.")
		case connection := <-manager.Unregister:
			close(connection.Data)
			delete(manager.Clients, connection)
			log.Info("A connection was terminated.")
		case message := <-manager.Broadcast:
			for connection := range manager.Clients {
				if message.Id == connection.Id { // Don't echo messages back
					continue
				}
				log.Infof("Words are being written to %+v", connection.Id)
				select {
				case connection.Data <- message.Message:
				default:
					close(connection.Data)
					delete(manager.Clients, connection)
				}
			}
		}
	}
}

func (manager *ClientManager) Receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.Socket.Read(message)
		if err != nil {
			manager.Unregister <- client
			client.Socket.Close()
			break
		}
		//if length > 0 {
		//	manager.Broadcast <- message
		//}
		client.Buffer.Write(message)
		for x := 0; x < length; x++ {
			mInfo := &MessageInfo{
				Message: client.Buffer.Bytes(),
				Id: client.Id,
			}
			if message[x] == 13 {
				manager.Broadcast <- mInfo
				client.Buffer.Reset()
				//log.Info("A return was hit: %+v", message)
			}
		}
	}
}

func (manager *ClientManager) Send(client *Client) {
	defer client.Socket.Close()
	for {
		select {
		case message, ok := <-client.Data:
			if !ok {
				return
			}
			client.Socket.Write(message)
		}
	}
}