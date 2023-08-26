package socket

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Room is a pool of connections
type Room struct {
	// Registered connections.
	RoomID     string `json:"room_id"`
	Private    bool   `json:"private"`
	MaxPlayers int    `json:"max_players"`
	Register   chan *Client
	UnRegister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewRoom(private bool, maxPlayers int) *Room {
	return &Room{
		RoomID:     uuid.New().String()[0:8],
		Private:    private,
		MaxPlayers: maxPlayers,
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (room *Room) Start() {
	for {
		select {
		case client := <-room.Register:
			room.registerClient(client)
		case client := <-room.UnRegister:
			room.unregisterClient(client)
		case message := <-room.Broadcast:
			room.broadcastMessage(message.encode())
		}
	}
}

func (room *Room) registerClient(client *Client) {
	if len(room.Clients) >= room.MaxPlayers {
		log.Info("Room is full")
		return
	} else {
		room.Clients[client] = true
		if len(room.Clients) > 1 {
			log.Info("A new user has joined the chat")
		}
	}

}

func (room *Room) unregisterClient(client *Client) {
	delete(room.Clients, client)
	if len(room.Clients) == 0 {
		delete(client.WsServer.Rooms, room)
		log.Info(fmt.Sprintf("Room %s has been deleted", room.RoomID))
		return
	}
	log.Info("A user left the chat")
	msg, err := json.Marshal(Message{
		Action: ChatAction,
		Target: MessageRoom{
			RoomID: room.RoomID,
		},
		Sender: MessageClient{
			ClientName: client.ClientName,
			ClientID:   client.ClientID,
			AvatarUrl:  client.AvatarUrl,
		},
		Payload: fmt.Sprintf("%s has left the room", client.ClientName),
	})
	if err != nil {
		log.Error(err)
	}
	for client := range room.Clients {
		client.send <- msg
	}
}

func (room *Room) broadcastMessage(message []byte) {
	log.Info("Broadcasting message to all clients in room")
	log.Info(fmt.Sprintf("Number of clients in room %s: %d", room.RoomID, len(room.Clients)))
	for client := range room.Clients {
		client.send <- message
	}
}

func (room *Room) GetAllClientInRoom(){
	for client := range room.Clients {
		log.Info(client.ClientName)
	}
}