package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dimeko/wserver/models"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type Client struct {
	uuid uuid.UUID
	conn *websocket.Conn
}

type Connections map[uuid.UUID]*Client

type ChatMessage struct {
	Sndr_id uuid.UUID
	Rcv_id  uuid.UUID
	// Room      string
	Broadcast bool
	Message   string
}

type Server struct {
	conns        Connections
	message_orch chan ChatMessage
}

func NewServer() *Server {
	return &Server{
		conns:        make(Connections),
		message_orch: make(chan ChatMessage),
	}
}

func (server *Server) WebSocketHandler(ws *websocket.Conn) {
	client_uuid, err := uuid.NewUUID()
	if err != nil {
		log.Printf("Connection %s could not be maintaned", ws.RemoteAddr())
	}
	log.Printf("New client connected: UUID: %s, Address: %s\n", client_uuid, ws.RemoteAddr())
	server.conns[client_uuid] = &Client{
		uuid: client_uuid,
		conn: ws,
	}
	server.keepConnection(ws, client_uuid, server.message_orch)
}

func (server *Server) keepConnection(ws *websocket.Conn, client_uuid uuid.UUID, message_orch chan ChatMessage) {
	buf := make([]byte, 1024)
	for {
		size, err := ws.Read(buf)
		if err != nil {
			log.Printf("Something wrong with client %s. Error %s", ws.RemoteAddr(), err.Error())
			if err == io.EOF {
				delete(server.conns, client_uuid)
				break
			}
		}
		msg := buf[:size]
		var decoded_msg = &models.ClientToServerMsg{}
		log.Printf("Before decoding message: %s\n", msg)

		err = json.Unmarshal(msg, decoded_msg)
		if err != nil {
			log.Printf("Could not decode message. Error: %s", string(err.Error()))
			continue
		}

		var recv_uuid uuid.UUID
		if decoded_msg.Broadcast {
			recv_uuid = uuid.Nil
		} else {
			recv_uuid, err = uuid.Parse(decoded_msg.Rcv_id)
			if err != nil {
				log.Printf("Could not decode rcv_id. Error: %s", string(err.Error()))
				continue
			}
		}

		new_msg := ChatMessage{
			Message:   decoded_msg.Message,
			Rcv_id:    recv_uuid,
			Sndr_id:   client_uuid,
			Broadcast: decoded_msg.Broadcast,
		}
		message_orch <- new_msg
	}
}

func StartServer() {
	// websocket.
	fmt.Println("Starting server")
	hubMapping := NewServer()
	go HubOrchestrator(hubMapping)
	http.Handle("/ws", websocket.Handler(hubMapping.WebSocketHandler))
	if err := http.ListenAndServe(":1337", nil); err != nil {
		panic(err)
	}
}
